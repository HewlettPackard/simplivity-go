package main

import (
	"fmt"

	"github.com/HewlettPackard/simplivity-go/ovc"
)

func main() {
	var (
		backupName   = "vm_new_back"
		backupByName *ovc.Backup
	)

	//Create an ovc client
	client, err := ovc.NewClient("sijeesh.kattumunda@demo.local", "Sijenov@2019", "10.30.8.245", "")
	//client, err := ovc.NewClient("username", "password", "ovc_ip", "certificate_path_if_needed")
	if err != nil {
		fmt.Println(err)
	}

	//Get all backup resources without Filter
	fmt.Println("\nGet all backups without params")
	backupList, err := client.Backups.GetAll(ovc.GetAllParams{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(backupList.Limit, backupList.Count, backupList.Offset, backupList.Members[0])
	for _, backup := range backupList.Members {
		fmt.Println(backup.Name + "\n")
	}

	//Get All backup resources with Filters
	fmt.Println("\nGet all policies with params")
	backupList, err = client.Backups.GetAll(ovc.GetAllParams{Limit: 1, Filters: map[string]string{"name": backupName}})
	if err != nil {
		fmt.Println(err)
	}
	for _, backup := range backupList.Members {
		fmt.Println(backup.Name + "\n")
	}

	//Get a backup resource by its name
	fmt.Println("\nGet a backup resource by it's name.")
	backupByName, err = client.Backups.GetByName(backupName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(backupByName)

	//Get a backup resource by its id
	fmt.Println("\nGet a backup resource by it's id.")
	backupById, err := client.Backups.GetById(backupByName.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(backupById)

	err = backupById.Delete()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Deleted")
}
