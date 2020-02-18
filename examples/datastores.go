package main

import (
	"fmt"

	"github.com/HewlettPackard/simplivity-go/ovc"
)

func main() {
	var (
		dstName   = "SVT_Montreal01"
		dstByName *ovc.Datastore
	)

	//Create an ovc client
	client, err := ovc.NewClient("username", "password", "ovc_ip", "certificate_path_if_needed")
	if err != nil {
		fmt.Println(err)
	}

	//Get all Datastore resources without Filter
	fmt.Println("\nGet all datastores without params")
	dstList, err := client.Datastores.GetAll(ovc.GetAllParams{})
	if err != nil {
		fmt.Println(err)
	}
	for _, dst := range dstList.Members {
		fmt.Println(dst.Name + "\n")
	}

	//Get All Datastore resources with Filters
	fmt.Println("\nGet all datastores with params")
	dstList, err = client.Datastores.GetAll(ovc.GetAllParams{Limit: 1, Filters: map[string]string{"name": dstName}})
	if err != nil {
		fmt.Println(err)
	}
	for _, dst := range dstList.Members {
		fmt.Println(dst.Name + "\n")
	}

	//Get a Datastore resource by its name
	fmt.Println("\nGet a datastore resource by it's name.")
	dstByName, err = client.Datastores.GetByName(dstName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dstByName)

	//Get a Datastore resource by its id
	fmt.Println("\nGet a datastore resource by it's id.")
	dstById, err := client.Datastores.GetById(dstByName.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dstById)
}
