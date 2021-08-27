package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
)

func main() {

	var name string
	flag.StringVar(&name, "name", "*", "Specify full or partial name of instance.")
	flag.StringVar(&name, "n", "*", "Specify full or partial name of instance.")

	var region string
	flag.StringVar(&region, "region", "us-east-1", "Specify the region you want to search in")
	flag.StringVar(&region, "r", "us-east-1", "Specify the region you want to search in")

	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	//nameQuery := "*" + name + "*"

	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		o.Region = region
	})

	resp, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
				},
			},
			{
				Name: aws.String("tag:Name"),
				Values: []string{
					//nameQuery,
					name,
				},
			},
		},
	})

	if err != nil {
		log.Fatalf("failed to list instances, %v", err)
	}

	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					fmt.Println("Name: " + *t.Value + " - IP Address: " + *i.PrivateIpAddress)

				}
			}
		}

	}

}
