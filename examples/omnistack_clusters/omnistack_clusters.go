package main

import (
	"fmt"

	"github.com/HewlettPackard/simplivity-go/ovc"
)

func main() {
	var (
		clusterName   = "RemoteCluster"
		clusterByName *ovc.OmniStackCluster
	)

	//Create an ovc client
	client, err := ovc.NewClient("username", "password", "ovc_ip", "certificate_path_if_needed")
	if err != nil {
		fmt.Println(err)
	}

	//Get all Cluster resources without Filter
	fmt.Println("\nGet all OmniStackClusters without params")
	clusterList, err := client.OmniStackClusters.GetAll(ovc.GetAllParams{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(clusterList.Limit, clusterList.Count, clusterList.Offset)
	for _, cluster := range clusterList.Members {
		fmt.Println(cluster.Name + "\n")
	}

	//Get All Cluster resources with Filters
	fmt.Println("\nGet all policies with params")
	clusterList, err = client.OmniStackClusters.GetAll(ovc.GetAllParams{Limit: 1, Filters: map[string]string{"name": clusterName}})
	if err != nil {
		fmt.Println(err)
	}
	for _, cluster := range clusterList.Members {
		fmt.Println(cluster.Name + "\n")
	}

	//Get a Cluster resource by its name
	fmt.Println("\nGet a Cluster resource by it's name.")
	clusterByName, err = client.OmniStackClusters.GetByName(clusterName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(clusterByName.Name)

	//Get a Cluster resource by its id
	fmt.Println("\nGet a Cluster resource by it's id.")
	clusterById, err := client.OmniStackClusters.GetById(clusterByName.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(clusterById.Name)
}
