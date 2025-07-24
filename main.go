package main

import (
	"bufio"
	"devops-tools/utils"
	"fmt"
	"os"
	"strconv"
)

type MenuItem struct {
	Label  string
	Action func() (*Result, error)
}

type Menu struct {
	Label  string
	Items  []MenuItem
	IsMain bool
}

func (menu Menu) Run() {
	for {
		utils.ClearScreen()

		var labels []string

		for _, item := range menu.Items {
			labels = append(labels, item.Label)
		}
		if menu.IsMain {
			labels = append(labels, "Exit")
		} else {
			labels = append(labels, "Back")
		}

		choice := chooseItemWithRetry(menu.Label, labels)

		if choice == len(labels) {
			break
		}

		result, err := menu.Items[choice-1].Action()
		if err != nil {
			utils.PrintError(err.Error())
			utils.PrintPressEnterToContinue()
		} else if result != nil {
			result.Print()
			utils.PrintPressEnterToContinue()
		}
	}
}

type Result struct {
	Title  string
	Values []string
}

func (result Result) Print() {
	utils.PrintResult(result.Title, result.Values)
}

var NotImplementedResult = &Result{Title: "Not implemented", Values: []string{}}

var (
	mainMenu = Menu{
		Label:  "Main",
		IsMain: true,
		Items: []MenuItem{
			{Label: "ALB", Action: func() (*Result, error) {
				albMenu.Run()
				return nil, nil
			}},
			{Label: "ECS", Action: func() (*Result, error) {
				ecsMenu.Run()
				return nil, nil
			}},
			{Label: "EKS", Action: func() (*Result, error) {
				eksMenu.Run()
				return nil, nil
			}},
		},
	}
	albMenu = Menu{
		Label:  "ALB",
		IsMain: false,
		Items: []MenuItem{
			{
				Label: "List ALBs",
				Action: func() (*Result, error) {
					albArns, err := utils.ListAlbArns()
					if err != nil {
						return nil, err
					}
					return &Result{Title: "List of ALBs", Values: albArns}, nil
				},
			},
			{
				Label: "Highest listener rule priority",
				Action: func() (*Result, error) {
					albArns, err := utils.ListAlbArns()
					if err != nil {
						return nil, err
					}
					albChoice := chooseItemWithRetry("Select an ALB", albArns)
					listenerArns, err := utils.ListAlbListenerArns(albArns[albChoice-1])
					if err != nil {
						return nil, err
					}
					listenerChoice := chooseItemWithRetry("Select a listener", listenerArns)
					highestPriority, err := utils.HighestAlbListenerRulePriority(listenerArns[listenerChoice-1])
					if err != nil {
						return nil, err
					}
					return &Result{Title: "Highest listener rule priority", Values: []string{strconv.Itoa(highestPriority)}}, nil
				},
			},
		},
	}
	ecsMenu = Menu{
		Label:  "ECS",
		IsMain: false,
		Items: []MenuItem{
			{Label: "List ECS clusters", Action: func() (*Result, error) {
				clusterArns, err := utils.ListEcsClusters()
				if err != nil {
					return nil, err
				}
				return &Result{Title: "List of ECS clusters", Values: clusterArns}, nil
			}},
			{Label: "List ECS services", Action: func() (*Result, error) {
				clusterArns, err := utils.ListEcsClusters()
				if err != nil {
					return nil, err
				}
				clusterChoice := chooseItemWithRetry("Select a cluster", clusterArns)
				serviceArns, err := utils.ListEcsServices(clusterArns[clusterChoice-1])
				if err != nil {
					return nil, err
				}
				return &Result{Title: "List of ECS services", Values: serviceArns}, nil
			}},
			{Label: "Restart ECS service", Action: func() (*Result, error) {
				clusterArns, err := utils.ListEcsClusters()
				if err != nil {
					return nil, err
				}
				clusterChoice := chooseItemWithRetry("Select a cluster", clusterArns)
				serviceArns, err := utils.ListEcsServices(clusterArns[clusterChoice-1])
				if err != nil {
					return nil, err
				}
				serviceChoice := chooseItemWithRetry("Select a service", serviceArns)
				if confirm(fmt.Sprintf("Are you sure you want to restart the service: %s?", serviceArns[serviceChoice-1])) {
					return nil, utils.RestartEcsService(clusterArns[clusterChoice-1], serviceArns[serviceChoice-1])
				}
				return nil, nil
			}},
			{Label: "Rollback ECS service", Action: func() (*Result, error) {
				clusterArns, err := utils.ListEcsClusters()
				if err != nil {
					return nil, err
				}
				clusterChoice := chooseItemWithRetry("Select a cluster", clusterArns)
				serviceArns, err := utils.ListEcsServices(clusterArns[clusterChoice-1])
				if err != nil {
					return nil, err
				}
				serviceChoice := chooseItemWithRetry("Select a service", serviceArns)
				if confirm(fmt.Sprintf("Are you sure you want to rollback the service: %s?", serviceArns[serviceChoice-1])) {
					return nil, utils.RollbackEcsService(clusterArns[clusterChoice-1], serviceArns[serviceChoice-1])
				}
				return nil, nil
			}},
			{Label: "Get latest ECS service deployment status", Action: func() (*Result, error) {
				clusterArns, err := utils.ListEcsClusters()
				if err != nil {
					return nil, err
				}
				clusterChoice := chooseItemWithRetry("Select a cluster", clusterArns)
				serviceArns, err := utils.ListEcsServices(clusterArns[clusterChoice-1])
				if err != nil {
					return nil, err
				}
				serviceChoice := chooseItemWithRetry("Select a service", serviceArns)
				deploymentStatus, err := utils.GetLatestEcsServiceDeploymentStatus(clusterArns[clusterChoice-1], serviceArns[serviceChoice-1])
				if err != nil {
					return nil, err
				}
				return &Result{Title: "Latest ECS service deployment status", Values: []string{deploymentStatus}}, nil
			}},
		},
	}
	eksMenu = Menu{
		Label:  "EKS",
		IsMain: false,
		Items: []MenuItem{
			{Label: "EKS clusters", Action: func() (*Result, error) {
				return NotImplementedResult, nil
			}},
		},
	}
)

func main() {
	mainMenu.Run()
}

func chooseItem(label string, items []string) (int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	utils.PrintMenuTitle(label)
	for idx, item := range items {
		utils.PrintMenuItem(idx+1, item)
	}

	utils.PrintMenuChoice()

	var line string
	if scanner.Scan() {
		line = scanner.Text()
	} else if scanner.Err() != nil {
		return 0, fmt.Errorf("failed to scan: %s", scanner.Err())
	} else {
		return 0, fmt.Errorf("failed to read input")
	}

	choice, err := strconv.Atoi(line)
	if err != nil {
		return 0, fmt.Errorf("bad choice: %s", err)
	} else if choice < 1 || choice > len(items) {
		return 0, fmt.Errorf("invalid choice: %d", choice)
	} else {
		return choice, nil
	}
}

func chooseItemWithRetry(label string, items []string) int {
	for {
		choice, err := chooseItem(label, items)
		if err == nil {
			return choice
		} else {
			utils.PrintError(err.Error())
			utils.PrintPressEnterToContinue()
			utils.ClearScreen()
		}
	}
}

func confirm(label string) bool {
	choice, err := chooseItem(label, []string{"Yes", "No"})
	if err != nil {
		return false
	}
	return choice == 1
}
