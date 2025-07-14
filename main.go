package main

import (
	"fmt"
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
		fmt.Printf("\nSelect action:\n")
		if menu.IsMain {
			fmt.Printf("0. Exit\n")
		} else {
			fmt.Printf("0. Back\n")
		}
		for idx, item := range menu.Items {
			fmt.Printf("%d. %s\n", idx+1, item.Label)
		}

		var choice int
		fmt.Printf("Choice: ")
		_, err := fmt.Scanf("%d", &choice)

		if err != nil {
			fmt.Printf("Failed scan: %s\n", err)

		} else if choice < 0 || choice > len(menu.Items) {
			fmt.Printf("Invalid choice: %d\n", choice)

		} else if choice == 0 {
			return nil

		} else {
			if err := menu.Items[choice-1].Action(); err != nil {
				fmt.Printf("Error: %s\n", err)
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

func albUtils() error {
	var menu = Menu{
		Label: "Main",
		Items: []MenuItem{
			{Label: "Highest listener rule priority", Action: notImplemented},
		},
	}
	return menu.Run()
}

func ecsUtils() error {
	fmt.Println("No actions defined for ECS utils")
	return nil
}

func eksUtils() error {
	fmt.Println("No actions defined for EKS utils")
	return nil
}

func notImplemented() error {
	fmt.Println("Not implemented")
	return nil
}
