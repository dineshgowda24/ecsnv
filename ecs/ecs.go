package ecs

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

func Run(cluster, service, file, profile string) {
	awsClient, err := NewAWSClient(profile)
	if err != nil {
		log.Fatalln(err)
	}

	if cluster == "" {
		clusters, err := awsClient.GetECSClusters()
		if err != nil {
			log.Fatalln(err)
		}

		if len(clusters) == 0 {
			log.Println("No clusters found")
		}

		cluster = prompt("Select cluster", clusters)
	}

	if service == "" {
		services, err := awsClient.GetECSServices(cluster)
		if err != nil {
			log.Fatalln(err)
		}

		if len(services) == 0 {
			log.Println("No services found")
		}

		service = prompt("Select service", services)
	}

	taskDef, err := awsClient.GetECSTaskDef(cluster, service)
	if err != nil {
		log.Fatalln(err)
	}

	envs, err := awsClient.GetENVsFromECSTaskDef(taskDef)
	if err != nil {
		log.Fatalln(err)
	}

	if file != "" {
		write(envs, file)
	} else {
		print(envs)
	}
}

// print prints to STDOUT
func print(envs map[string]string) {
	for key, value := range envs {
		fmt.Printf("%v=%v\n", key, value)
	}
}

// writes the envs to file
func write(envs map[string]string, file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for key, value := range envs {
		_, err := f.WriteString(fmt.Sprintf("%v=%v\n", key, value))

		if err != nil {
			log.Fatal(err)
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v envs written to %v/%v\n", len(envs), pwd, file)
}

// prompt displays a terminal prompt
func prompt(label string, items []string) string {
	pmt := promptui.Select{
		Label: label,
		Items: items,
		Size:  25,
	}

	_, value, _ := pmt.Run()
	return value
}
