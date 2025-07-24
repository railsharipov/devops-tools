package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

var ecsClient *ecs.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ecsClient = ecs.NewFromConfig(cfg)
}

func ListEcsClusters() ([]string, error) {
	listClustersInput := &ecs.ListClustersInput{}
	clusterArns := []string{}
	for {
		listClustersOutput, err := ecsClient.ListClusters(context.TODO(), listClustersInput)
		if err != nil {
			return nil, err
		}
		clusterArns = append(clusterArns, listClustersOutput.ClusterArns...)
		if listClustersOutput.NextToken == nil {
			break
		}
		listClustersInput.NextToken = listClustersOutput.NextToken
	}

	return clusterArns, nil
}

func ListEcsServices(clusterArn string) ([]string, error) {
	listServicesInput := &ecs.ListServicesInput{
		Cluster: aws.String(clusterArn),
	}
	serviceArns := []string{}
	for {
		listServicesOutput, err := ecsClient.ListServices(context.TODO(), listServicesInput)
		if err != nil {
			return nil, err
		}
		serviceArns = append(serviceArns, listServicesOutput.ServiceArns...)
		if listServicesOutput.NextToken == nil {
			break
		}
		listServicesInput.NextToken = listServicesOutput.NextToken
	}

	return serviceArns, nil
}

func ListEcsTasks(clusterArn, serviceArn string) ([]string, error) {
	listTasksInput := &ecs.ListTasksInput{
		Cluster:     aws.String(clusterArn),
		ServiceName: aws.String(serviceArn),
	}
	tasks := []string{}
	for {
		listTasksOutput, err := ecsClient.ListTasks(context.TODO(), listTasksInput)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, listTasksOutput.TaskArns...)
		if listTasksOutput.NextToken == nil {
			break
		}
		listTasksInput.NextToken = listTasksOutput.NextToken
	}
	return tasks, nil
}

func GetEcServiceTaskDefinition(serviceArn string) (string, error) {
	describeServiceInput := &ecs.DescribeServicesInput{
		Services: []string{serviceArn},
	}
	describeServiceOutput, err := ecsClient.DescribeServices(context.TODO(), describeServiceInput)
	if err != nil {
		return "", err
	}
	return *describeServiceOutput.Services[0].TaskDefinition, nil
}

func GetEcsTaskDefinitionFamily(taskDefinitionArn string) (string, error) {
	describeTaskDefinitionInput := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}
	describeTaskDefinitionOutput, err := ecsClient.DescribeTaskDefinition(context.TODO(), describeTaskDefinitionInput)
	if err != nil {
		return "", err
	}
	return *describeTaskDefinitionOutput.TaskDefinition.Family, nil
}

func ListLatestEcsTaskDefinitions(taskFamily string) ([]string, error) {
	listTaskDefinitionsInput := &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(taskFamily),
	}
	taskDefinitions := []string{}
	listTaskDefinitionsOutput, err := ecsClient.ListTaskDefinitions(context.TODO(), listTaskDefinitionsInput)
	if err != nil {
		return nil, err
	}
	taskDefinitions = append(taskDefinitions, listTaskDefinitionsOutput.TaskDefinitionArns...)

	return taskDefinitions, nil
}

func RestartEcsService(clusterArn, serviceArn string) error {
	updateServiceInput := &ecs.UpdateServiceInput{
		Cluster:            aws.String(clusterArn),
		Service:            aws.String(serviceArn),
		ForceNewDeployment: true,
	}
	_, err := ecsClient.UpdateService(context.TODO(), updateServiceInput)
	log.Printf("Requested ECS service restart for %s in cluster %s", serviceArn, clusterArn)
	return err
}

func RollbackEcsService(clusterArn, serviceArn string) error {
	describeServiceInput := &ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterArn),
		Services: []string{serviceArn},
	}
	describeServiceOutput, err := ecsClient.DescribeServices(context.TODO(), describeServiceInput)
	if err != nil {
		return err
	}
	taskDefinitionArn := *describeServiceOutput.Services[0].TaskDefinition
	describeTaskDefinitionInput := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}
	describeTaskDefinitionOutput, err := ecsClient.DescribeTaskDefinition(context.TODO(), describeTaskDefinitionInput)
	if err != nil {
		return err
	}
	taskDef := describeTaskDefinitionOutput.TaskDefinition
	taskDefFamily := *taskDef.Family
	taskDefRevision := taskDef.Revision
	previousTaskDefinition := fmt.Sprintf("%s:%d", taskDefFamily, taskDefRevision-1)

	updateServiceInput := &ecs.UpdateServiceInput{
		Cluster:            aws.String(clusterArn),
		Service:            aws.String(serviceArn),
		TaskDefinition:     aws.String(previousTaskDefinition),
		ForceNewDeployment: true,
	}
	_, err = ecsClient.UpdateService(context.TODO(), updateServiceInput)
	log.Printf("Requested ECS service rollback for %s in cluster %s", serviceArn, clusterArn)

	return err
}

func GetLatestEcsServiceDeploymentStatus(clusterArn, serviceArn string) (string, error) {
	listServiceDeploymentsInput := &ecs.ListServiceDeploymentsInput{
		Cluster:    aws.String(clusterArn),
		Service:    aws.String(serviceArn),
		MaxResults: aws.Int32(1),
	}
	listServiceDeploymentsOutput, err := ecsClient.ListServiceDeployments(context.TODO(), listServiceDeploymentsInput)
	if err != nil {
		return "", err
	}
	serviceDeployments := listServiceDeploymentsOutput.ServiceDeployments
	if len(serviceDeployments) == 0 {
		return "", fmt.Errorf("no service deployments found for service %s in cluster %s", serviceArn, clusterArn)
	}
	return string(serviceDeployments[0].Status), nil
}
