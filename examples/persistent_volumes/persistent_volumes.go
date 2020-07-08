package main

import (
	"fmt"
	"time"

	"github.com/HewlettPackard/simplivity-go/ovc"
)

func main() {
	var (
		pvName      = "pvc-test_fcd"
		pvByName    *ovc.PersistentVolume
		currentTime = time.Now().String()
	)

	//Create an ovc client
	client, err := ovc.NewClient("username", "password", "ovc_ip", "certificate_path_if_needed")
	if err != nil {
		fmt.Println(err)
	}

	//Get all persistent volume resources without Filter
	fmt.Println("\nGet all persistent volumes without params")
	pvList, err := client.PersistentVolumes.GetAll(ovc.GetAllParams{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pvList.Limit, pvList.Count, pvList.Offset, pvList.Members[0])
	for _, pv := range pvList.Members {
		fmt.Println(pv.Name + "\n")
	}

	//Get All PV resources with Filters
	fmt.Println("\nGet all PVs with params")
	pvList, err = client.PersistentVolumes.GetAll(ovc.GetAllParams{Limit: 1, Filters: map[string]string{"name": pvName}})
	if err != nil {
		fmt.Println(err)
	}
	for _, pv := range pvList.Members {
		fmt.Println(pv.Name + "\n")
	}

	//Get a PV resource by its name
	fmt.Println("\nGet a PV resource by it's name.")
	pvByName, err = client.PersistentVolumes.GetByName(pvName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pvByName)

	//Get a PV resource by its id
	fmt.Println("\nGet a VM resource by it's id.")
	pvById, err := client.PersistentVolumes.GetById(pvByName.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pvById)

	//Set policy for a single and multiple PVs
	policyList, err := client.Policies.GetAll(ovc.GetAllParams{Limit: 1})
	if err != nil {
		fmt.Println(err)
	}
	if policyList != nil {
		//Get one of the backup policies
		policy := policyList.Members[0]
		fmt.Println("\nPolicy to set is ", policy)

		//Set policy for multiple PVs
		fmt.Println("\nSet policy for multiple PVs")
		pvs := []*ovc.PersistentVolume{pvByName}
		err := client.PersistentVolumes.SetPolicyForMultiplePVs(policy, pvs)
		if err != nil {
			fmt.Println(err)
		}
	}

	//Take a backup of a Persistent Volume
	fmt.Println("Take backup of a Peristent Volume")
	backupName := "backup_" + currentTime
	backReq := &ovc.CreateBackupRequest{Name: backupName}
	_, err = pvByName.CreateBackup(backReq, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Backup operation was Successful")
	}

	// Get backups of a PV
	fmt.Println("\nGet all backups of a PV")
	backupList, err := pvByName.GetBackups()
	if err != nil {
		fmt.Println(err)
	}
	for _, backup := range backupList.Members {
		fmt.Println(backup.Name)
	}
}
