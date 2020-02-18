package main

import (
	"fmt"

	"github.com/HewlettPackard/simplivity-go/ovc"
)

func main() {
	var (
		hostName   = "omnicube308240.cloud.local"
		hostByName *ovc.Host
	)

	//Create an ovc client
	client, err := ovc.NewClient("username", "password", "ovc_ip", "certificate_path_if_needed")
	if err != nil {
		fmt.Println(err)
	}

	//Get all hosts resources without Filter
	fmt.Println("\nGet all Hosts without params")
	hostList, err := client.Hosts.GetAll(ovc.GetAllParams{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hostList.Limit, hostList.Count, hostList.Offset)
	for _, hosts := range hostList.Members {
		fmt.Println(hosts.Name + "\n")
	}

	//Get All hosts resources with Filters
	fmt.Println("\nGet all hosts with params")
	hostList, err = client.Hosts.GetAll(ovc.GetAllParams{Limit: 1, Filters: map[string]string{"name": hostName}})
	if err != nil {
		fmt.Println(err)
	}
	for _, hosts := range hostList.Members {
		fmt.Println(hosts.Name + "\n")
	}

	//Get a hosts resource by its name
	fmt.Println("\nGet a hosts resource by it's name.")
	hostByName, err = client.Hosts.GetByName(hostName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hostByName)

	//Get a hosts resource by its id
	fmt.Println("\nGet a hosts resource by it's id.")
	hostsById, err := client.Hosts.GetById(hostByName.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hostsById)
}
