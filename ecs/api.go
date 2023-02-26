package ecs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// AWSClient encapsulates serveral helper method to retrive data from AWS cloud
type AWSClient struct {
	region string
	sn     *session.Session
}

// NewAWSClient returns a new AWSClient
func NewAWSClient(region string) (*AWSClient, error) {
	sn, err := session.NewSession(&aws.Config{
		Region: &region,
	})
	if err != nil {
		return nil, err
	}

	return &AWSClient{
		region: region,
		sn:     sn,
	}, nil
}

// GetECSClusters returns a list of ecs clusters in a region
func (a *AWSClient) GetECSClusters() ([]string, error) {
	svc := ecs.New(a.sn)
	maxResults := int64(100)

	input := &ecs.ListClustersInput{
		MaxResults: &maxResults,
	}
	var clusters []string
	var nextToken string
	for {
		if nextToken != "" {
			input.NextToken = &nextToken
		}

		result, err := svc.ListClusters(input)
		if err != nil {
			return []string{}, err
		}

		for _, arn := range result.ClusterArns {
			cluster := strings.Split((*arn), "cluster/")[1]
			clusters = append(clusters, cluster)
		}

		if result.NextToken == nil {
			break
		} else {
			nextToken = *result.NextToken
		}

	}

	return clusters, nil
}

// GetECSServices returns a list of ecs services for a cluster in a region
func (a *AWSClient) GetECSServices(cluster string) ([]string, error) {
	svc := ecs.New(a.sn)
	maxResults := int64(100)

	input := &ecs.ListServicesInput{
		Cluster:    &cluster,
		MaxResults: &maxResults,
	}
	var services []string
	var nextToken string
	for {
		if nextToken != "" {
			input.NextToken = &nextToken
		}

		result, err := svc.ListServices(input)
		if err != nil {
			return []string{}, err
		}

		for _, arn := range result.ServiceArns {
			service := strings.Split((*arn), fmt.Sprintf("/%v/", cluster))[1]
			services = append(services, service)
		}

		if result.NextToken == nil {
			break
		} else {
			nextToken = *result.NextToken
		}

	}

	return services, nil
}

// GetECSTaskDef returns the current task definition for a given cluster and service
func (a *AWSClient) GetECSTaskDef(cluster, service string) (string, error) {
	svc := ecs.New(a.sn)
	input := &ecs.DescribeServicesInput{
		Cluster:  &cluster,
		Services: []*string{&service},
	}
	result, err := svc.DescribeServices(input)
	if err != nil {
		return "", err
	}

	return *result.Services[0].TaskDefinition, nil
}

// GetENVsFromECSTaskDef returns enviroment variable in a task definition
func (a *AWSClient) GetENVsFromECSTaskDef(taskDef string) (map[string]string, error) {
	svc := ecs.New(a.sn)

	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDef,
	}
	result, err := svc.DescribeTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	taskDefObj := result.TaskDefinition
	if taskDefObj == nil {
		return nil, errors.New("missing task definitions")
	}

	containerDefs := taskDefObj.ContainerDefinitions
	if containerDefs == nil {
		return nil, errors.New("missing container definitions")
	}

	containerDef := containerDefs[0]
	envs := make(map[string]string)

	for _, kvPair := range containerDef.Environment {
		envs[*kvPair.Name] = *kvPair.Value
	}

	return envs, nil
}
