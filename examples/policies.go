package main

import (
	"fmt"

	"github.com/HewlettPackard/simplivity-go/ovc"
)

func main() {
	var (
		policyName   = ""
		policyByName *ovc.Policy
	)

	//Create an ovc client
	client, err := ovc.NewClient("username", "password", "ovc_ip", "certificate_path_if_needed")
	if err != nil {
		fmt.Println(err)
	}

	//Get all Policy resources without Filter
	fmt.Println("\nGet all Policies without params")
	policyList, err := client.Policies.GetAll(ovc.GetAllParams{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(policyList.Limit, policyList.Count, policyList.Offset)
	for _, policy := range policyList.Members {
		fmt.Println(policy.Name + "\n")
	}

	//Get All Policy resources with Filters
	fmt.Println("\nGet all policies with params")
	policyList, err = client.Policies.GetAll(ovc.GetAllParams{Limit: 1, Filters: map[string]string{"name": policyName}})
	if err != nil {
		fmt.Println(err)
	}
	for _, policy := range policyList.Members {
		fmt.Println(policy.Name + "\n")
	}

	//Get a Policy resource by its name
	fmt.Println("\nGet a Policy resource by it's name.")
	policyByName, err = client.Policies.GetByName(policyName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(policyByName)

	//Get a Policy resource by its id
	fmt.Println("\nGet a Policy resource by it's id.")
	policyById, err := client.Policies.GetById(policyByName.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(policyById)
}
