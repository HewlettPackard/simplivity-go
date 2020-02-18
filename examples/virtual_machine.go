package main

import (
	"fmt"
	"time"

	"github.com/HewlettPackard/simplivity-go/ovc"
)

func main() {
	var (
		vmName        = "clone_test"
		datastoreName = "SVT_Montreal01"
		vmByName      *ovc.VirtualMachine
		currentTime   = time.Now().String()
	)

	//Create an ovc client
	client, err := ovc.NewClient("username", "password", "ovc_ip", "certificate_path_if_needed")
	if err != nil {
		fmt.Println(err)
	}

	//Get all VM resources without Filter
	fmt.Println("\nGet all VMs without params")
	vmList, err := client.VirtualMachines.GetAll(ovc.GetAllParams{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(vmList.Count, vmList.Limit)
	for _, vm := range vmList.Members {
		fmt.Println(vm.Name + "\n")
		fmt.Println(vm.Id)
	}

	//Get All VM resources with Filters
	fmt.Println("\nGet all VMs with params")
	vmList, err = client.VirtualMachines.GetAll(ovc.GetAllParams{Limit: 1, Filters: map[string]string{"name": vmName}})
	if err != nil {
		fmt.Println(err)
	}
	for _, vm := range vmList.Members {
		fmt.Println(vm.Name + "\n")
	}

	//Get a VM resource by its name
	fmt.Println("\nGet a VM resource by it's name.")
	vmByName, err = client.VirtualMachines.GetByName(vmName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(vmByName)

	//Get a VM resource by its id
	fmt.Println("\nGet a VM resource by it's id.")
	vmById, err := client.VirtualMachines.GetById(vmByName.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(vmById)

	//Set policy for a single and multiple VMs
	policyList, err := client.Policies.GetAll(ovc.GetAllParams{Limit: 1})
	if err != nil {
		fmt.Println(err)
	}
	if policyList != nil {
		//Get one of the backup policies
		policy := policyList.Members[0]

		//Set policy for multiple VMs
		fmt.Println("\nSet policy for multiple VMs")
		vms := []*ovc.VirtualMachine{vmByName}
		err := client.VirtualMachines.SetPolicyForMultipleVMs(policy, vms)
		if err != nil {
			fmt.Println(err)
		}

		//Set policy for a single VM
		fmt.Println("\nSet policy for a single VM")
		err = vmByName.SetPolicy(policy)
		if err != nil {
			fmt.Println(err)
		}
	}

	//Clone a VM resource
	fmt.Println("\nClone a VM resource")
	clonedVM, err := vmByName.Clone("clone_from_go", false)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Clone operation was Successful")
	}

	//Move a VM to another datastore
	fmt.Println("\nMove a VM to another datastore", clonedVM.Name)
	datastore, err := client.Datastores.GetByName(datastoreName)
	if err != nil {
		fmt.Println(err)
	} else {
		vm, err := clonedVM.Move("move_from_go", datastore)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Move operation was Successful", vm.Name)
		}
	}

	//Take a backup of a Virtual Machine
	fmt.Println("Take backup of a Virtual Machine")
	backupName := "backup_" + currentTime
	backReq := &ovc.CreateBackupRequest{Name: backupName}
	_, err = vmByName.CreateBackup(backReq, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Backup operation was Successful")
	}

	// Get backups of a VM
	fmt.Println("\nGet all backups of a VM")
	backupList, err := vmByName.GetBackups()
	if err != nil {
		fmt.Println(err)
	}
	for _, backup := range backupList.Members {
		fmt.Println(backup.Name)
	}

	//Set backup parameters of a Virtual Machine
	fmt.Println("Set backup parameters of a Virtual Machine")
	req := &ovc.SetBackupParametersRequest{Username: "username", Password: "password", OverrideValidation: false}
	err = vmByName.SetBackupParameters(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Set backup parameters operation was successful")
	}

}
