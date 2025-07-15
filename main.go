package main

import (
	"fmt"

	"node40-dev-tools/utils"
)

type MenuItem struct {
	Label  string
	Action func() error
}

type Menu struct {
	Label  string
	Items  []MenuItem
	IsMain bool
}

func (menu Menu) Run() error {
	for {
		var labels []string

		for _, item := range menu.Items {
			labels = append(labels, item.Label)
		}
		if menu.IsMain {
			labels = append(labels, "Exit")
		} else {
			labels = append(labels, "Back")
		}

		choice, err := chooseItem(menu.Label, labels)
		if err != nil {
			return err
		}

		if choice == len(labels) {
			return nil
		} else {
			if err := menu.Items[choice-1].Action(); err != nil {
				utils.PrintError(err.Error())
			}
		}
	}
}

func main() {
	var menu = Menu{
		Label:  "Main",
		IsMain: true,
		Items: []MenuItem{
			{Label: "ALB", Action: albUtils},
			{Label: "ECS", Action: ecsUtils},
			{Label: "EKS", Action: eksUtils},
		},
	}
	menu.Run()
}

func chooseItem(label string, items []string) (int, error) {
	for {
		utils.PrintMenuTitle(label)
		for idx, item := range items {
			utils.PrintMenuItem(idx+1, item)
		}

		var choice int
		utils.PrintMenuChoice()
		_, err := fmt.Scanf("%d", &choice)

		if err != nil {
			utils.PrintError(fmt.Sprintf("Failed scan: %s", err))
		} else if choice < 1 || choice > len(items) {
			utils.PrintError(fmt.Sprintf("Invalid choice: %d", choice))
		} else {
			return choice, nil
		}
	}
}

func albUtils() error {
	var menu = Menu{
		Label:  "ALB",
		IsMain: false,
		Items: []MenuItem{
			{
				Label: "List ALBs",
				Action: func() error {
					albArns, err := utils.ListAlbArns()
					if err != nil {
						return err
					}
					utils.PrintResultTitle("List of ALBs")
					for _, albArn := range albArns {
						utils.PrintResultItem(albArn)
					}
					return nil
				},
			},
			{
				Label: "Highest listener rule priority",
				Action: func() error {
					albArns, err := utils.ListAlbArns()
					if err != nil {
						return err
					}
					albChoice, err := chooseItem("Select an ALB", albArns)
					if err != nil {
						return err
					}
					listenerArns, err := utils.ListAlbListenerArns(albArns[albChoice-1])
					if err != nil {
						return err
					}
					listenerChoice, err := chooseItem("Select a listener", listenerArns)
					if err != nil {
						return err
					}
					highestPriority, err := utils.HighestAlbListenerRulePriority(listenerArns[listenerChoice-1])
					if err != nil {
						return err
					}
					utils.PrintResultValueInt("Highest listener rule priority", highestPriority)
					return nil
				},
			},
		},
	}
	return menu.Run()
}

func ecsUtils() error {
	var menu = Menu{
		Label:  "ECS",
		IsMain: false,
		Items: []MenuItem{
			{Label: "Not implemented", Action: notImplemented},
		},
	}
	return menu.Run()
}

func eksUtils() error {
	var menu = Menu{
		Label:  "EKS",
		IsMain: false,
		Items: []MenuItem{
			{Label: "Not implemented", Action: notImplemented},
		},
	}
	return menu.Run()
}

func notImplemented() error {
	utils.PrintWarning("This feature is not implemented yet")
	return nil
}
