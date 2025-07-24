package utils

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

var albClient *elasticloadbalancingv2.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	albClient = elasticloadbalancingv2.NewFromConfig(cfg)
}

func ListAlbArns() ([]string, error) {
	describeLoadBalancersInput := &elasticloadbalancingv2.DescribeLoadBalancersInput{}
	describeLoadBalancersOutput, err := albClient.DescribeLoadBalancers(context.TODO(), describeLoadBalancersInput)
	if err != nil {
		return nil, err
	}
	albArns := make([]string, len(describeLoadBalancersOutput.LoadBalancers))
	for i, lb := range describeLoadBalancersOutput.LoadBalancers {
		albArns[i] = *lb.LoadBalancerArn
	}
	return albArns, nil
}

func ListAlbListenerArns(albName string) ([]string, error) {
	describeListenersInput := &elasticloadbalancingv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(albName),
	}
	describeListenersOutput, err := albClient.DescribeListeners(context.TODO(), describeListenersInput)
	if err != nil {
		return nil, err
	}
	listenerArns := make([]string, len(describeListenersOutput.Listeners))
	for i, listener := range describeListenersOutput.Listeners {
		listenerArns[i] = *listener.ListenerArn
	}
	return listenerArns, nil
}

func HighestAlbListenerRulePriority(listenerArn string) (int, error) {
	describeRulesInput := &elasticloadbalancingv2.DescribeRulesInput{
		ListenerArn: aws.String(listenerArn),
	}
	describeRulesOutput, err := albClient.DescribeRules(context.TODO(), describeRulesInput)
	if err != nil {
		return 0, err
	}
	highestPriority := 0
	for _, rule := range describeRulesOutput.Rules {
		curPriority, err := strconv.Atoi(*rule.Priority)
		if err != nil {
			continue
		}
		if curPriority > highestPriority {
			highestPriority = curPriority
		}
	}
	return highestPriority, nil
}
