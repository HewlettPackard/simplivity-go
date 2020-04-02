package ovc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// VirtualMachineResource handles communications with the the VM resource methods
//
// SimpliVity API docs: https://developer.hpe.com/api/simplivity/endpoint?&path=%2Fvirtual_machines
type VirtualMachineResource resourceClient

// Virtual machine GetAll response
type VirtualMachineList struct {
	Offset  int               `json:"offset,omitempty"`
	Count   int               `json:"count,omitempty"`
	Limit   int               `json:"limit,omitempty"`
	Members []*VirtualMachine `json:"virtual_machines,omitempty"`
}

type ReplicaSetList struct {
	Role string `json:"role,omitempty"`
	Id   string `json:"id,omitempty"`
}

// VirtualMachine represents a SimpliVity VM
type VirtualMachine struct {
	Name                                   string           `json:"name,omitempty"`
	Id                                     string           `json:"id,omitempty"`
	State                                  string           `json:"state,omitempty"`
	CreatedAt                              string           `json:"created_at,omitempty"`
	DeletedAt                              string           `json:"deleted_at,omitempty"`
	PolicyName                             string           `json:"policy_name,omitempty"`
	PolicyId                               string           `json:"policy_id,omitempty"`
	DatastoreName                          string           `json:"datastore_name,omitempty"`
	DatastoreId                            string           `json:"datastore_id,omitempty"`
	OmniStackClusterName                   string           `json:"omnistack_cluster_name,omitempty"`
	OmniStackClusterId                     string           `json:"omnistack_cluster_id,omitempty"`
	AppAwareVMStatus                       string           `json:"app_aware_vm_status,omitempty"`
	HypervisorObjectId                     string           `json:"hypervisor_object_id,omitempty"`
	HypervisorType                         string           `json:"hypervisor_type,omitempty"`
	HypervisorManagementSystem             string           `json:"hypervisor_management_system,omitempty"`
	HypervisorManagementSystemName         string           `json:"hypervisor_management_system_name,omitempty"`
	HostId                                 string           `json:"host_id,omitempty"`
	ComputeClusterParentHypervisorObjectId string           `json:"compute_cluster_parent_hypervisor_object_id,omitempty"`
	ComputeClusterName                     string           `json:"cumpute_cluster_name,omitempty"`
	ClusterGroupIds                        []string         `json:"cluster_group_ids,omitempty"`
	ReplicaSet                             []ReplicaSetList `json:"replica_set,omitempty"`
}

// GetAll returns all the virtual machines filtered by the query parameters.
// Filters:
//   id: The unique identifier (UID) of the virtual_machines to return
//     Accepts: Single value, comma-separated list
//   name: The name of the virtual_machines to return
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   omnistack_cluster_id: The unique identifier (UID) of the omnistack_cluste
//     that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   omnistack_cluster_name: The name of the omnistack_cluster that
//     is associated with the instances to return.
//     Accepts: Single value, comma-separated list.
//   compute_cluster_parent_hypervisor_object_id: The unique identifier (UID)
//     of the hypervisor that contains the omnistack_cluster that is associated
//     with the instances to return
//     Accepts: Single value, comma-separated list.
//   compute_cluster_parent_name: The name of the hypervisor that contains the
//     omnistack_cluster that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   hypervisor_management_system: The IP address of the hypervisor associated
//     with the virtual machine.
//     Accepts: Single value, comma-separated list, pattern using one
//     or more asterisk characters as a wildcard
//   hypervisor_management_system_name: The name of the hypervisor associated
//     with the virtual machine
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   datastore_id: The unique identifier (UID) of the datastore that is associated
//     with the instances to return
//     Accepts: Single value, comma-separated list
//   datastore_name: The name of the datastore that is associated with the
//     instances to return
//     Accepts: Single value, comma-separated list
//   policy_id: The unique identifier (UID) of the policy that is associated
//     with the instances to return
//     Accepts: Single value, comma-separated list
//   policy_name: The name of the policy that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   hypervisor_object_id: The unique identifier (UID) of the hypervisor-based instance
//     that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   created_after: The earliest creation time after the virtual machines to return were
//     created, expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   created_before: The latest creation time before the virtual machines to return were
//     created, expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   state: The state of the virtual_machine that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   app_aware_vm_status: The status of the ability of the virtual machine to take
//     an application-consistent backup that uses Microsoft VSS
//     Accepts: Single value, comma-separated list
//   hypervisor_is_template: An indicator that shows if the virtual machine is a template.
//   host_id: The unique identifier (UID) of the virtual_machine host.
func (v *VirtualMachineResource) GetAll(params GetAllParams) (*VirtualMachineList, error) {
	var (
		path   = "/virtual_machines"
		vmList VirtualMachineList
	)

	qrStr := params.QueryString()

	resp, err := v.client.DoRequest("GET", path, qrStr, nil, nil)
	if err != nil {
		return &vmList, err
	}

	err = json.Unmarshal(resp, &vmList)
	if err != nil {
		return &vmList, err
	}

	return &vmList, nil
}

// GetBy searches for VM resources with single filter.
func (v *VirtualMachineResource) GetBy(field_name string, value string) ([]*VirtualMachine, error) {
	filters := map[string]string{field_name: value}
	vmList, err := v.GetAll(GetAllParams{Filters: filters})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	vms := vmList.Members

	return vms, nil
}

// GetByName searches for a VM by its name
func (v *VirtualMachineResource) GetByName(name string) (*VirtualMachine, error) {
	vms, err := v.GetBy("name", name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(vms) > 0 {
		vm := vms[0]
		return vm, nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// GetById searches for a VM by its id
func (v *VirtualMachineResource) GetById(id string) (*VirtualMachine, error) {
	vms, err := v.GetBy("id", id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(vms) > 0 {
		vm := vms[0]
		return vm, nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// SetPolicyForMultipleVMs sets a policy for list of VM resources.
func (v *VirtualMachineResource) SetPolicyForMultipleVMs(policy *Policy, vms []*VirtualMachine) error {
	var (
		path = fmt.Sprintf("/virtual_machines/set_policy")
	)

	if len(vms) < 1 {
		return errors.New("Pass a list of VM resoures")
	}

	vm_ids := []string{}
	for _, vm := range vms {
		vm_ids = append(vm_ids, vm.Id)
	}

	body := map[string]interface{}{"policy_id": policy.Id, "virtual_machine_id": vm_ids}

	resp, err := v.client.DoRequest("POST", path, "", body, nil)
	if err != nil {
		return err
	}

	_, err = commonClient.Tasks.WaitForTask(resp)
	if err != nil {
		return err
	}

	return nil
}

// SetPolicy sets policy for single VM resource.
func (v *VirtualMachine) SetPolicy(policy *Policy) error {
	var (
		path = fmt.Sprintf("/virtual_machines/%s/set_policy", v.Id)
	)

	body := map[string]string{"policy_id": policy.Id}
	resp, err := commonClient.DoRequest("POST", path, "", body, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	task, err := commonClient.Tasks.WaitForTask(resp)
	if err != nil {
		log.Println(err)
		return err
	}

	if len(task.AffectedResources) < 1 {
		err_message := "Set policy was not successful. Error code:" + string(task.ErrorCode)
		return errors.New(err_message)
	}

	return nil
}

// Clone creates a clone of the VM.
func (v *VirtualMachine) Clone(new_vm_name string, app_consistent bool) (*VirtualMachine, error) {
	var (
		path = fmt.Sprintf("/virtual_machines/%s/clone", v.Id)
	)

	body := map[string]interface{}{"virtual_machine_name": new_vm_name,
		"app_consistent": app_consistent}
	resp, err := commonClient.DoRequest("POST", path, "", body, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	task, err := commonClient.Tasks.WaitForTask(resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resources := task.AffectedResources
	if len(resources) < 1 {
		err_message := "Clone was not successful. Error code:" + string(task.ErrorCode)
		return nil, errors.New(err_message)
	}

	resource_id := resources[0].ObjectId
	clonedVM, err := commonClient.VirtualMachines.GetById(resource_id)

	return clonedVM, nil
}

// Move moves a VM from one datastore to another.
func (v *VirtualMachine) Move(vm_name string, datastore *Datastore) (*VirtualMachine, error) {
	var (
		path = fmt.Sprintf("/virtual_machines/%s/move", v.Id)
	)

	body := map[string]interface{}{"virtual_machine_name": vm_name,
		"destination_datastore_id": datastore.Id}

	resp, err := commonClient.DoRequest("POST", path, "", body, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	task, err := commonClient.Tasks.WaitForTask(resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resources := task.AffectedResources
	if len(resources) < 1 {
		err_message := "Move was not successful. Error code:" + string(task.ErrorCode)
		return nil, errors.New(err_message)
	}

	resource_id := resources[0].ObjectId
	clonedVM, err := commonClient.VirtualMachines.GetById(resource_id)

	return clonedVM, nil
}

// CreateBackup request body
type CreateBackupRequest struct {
	// The name of the new backup created from this action
	Name string `json:"backup_name,omitempty"`

	// The consistency type of the backup
	// Default: "NONE"
	ConsistencyType string `json:"consistency_type,omitempty"`

	// An indicator to show if the backup represents a snapshot
	// of a virtual machine with data that was first flushed to disk
	// Default: false
	AppConsistent bool `json:"app_consistent,omitempty"`

	// The number of minutes to keep backups
	// Default: 0(indicates that the backup never expires)
	Retention int `json:"retention,omitempty"`

	// omnistack_cluster that stores the new backup
	// Default: local omnistack_cluster
	Destination string `json:"destination_id,omitempty"`
}

// CreateBackup creates a back of the VM.
func (v *VirtualMachine) CreateBackup(req *CreateBackupRequest, dest *OmniStackCluster) (*Backup, error) {
	var (
		path = fmt.Sprintf("/virtual_machines/%s/backup", v.Id)
	)

	if dest != nil {
		req.Destination = dest.Id
	}

	resp, err := commonClient.DoRequest("POST", path, "", req, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	task, err := commonClient.Tasks.WaitForTask(resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resources := task.AffectedResources
	if len(resources) < 1 {
		err_message := "Backup was not successful. Error code:" + string(task.ErrorCode)
		return nil, errors.New(err_message)
	}

	resource_id := resources[0].ObjectId
	backup, err := commonClient.Backups.GetById(resource_id)

	return backup, nil
}

// GetBackups gets all the backups of a VM.
func (v *VirtualMachine) GetBackups() (*BackupList, error) {
	var (
		path       = fmt.Sprintf("/virtual_machines/%s/backups", v.Id)
		backupList BackupList
	)

	resp, err := commonClient.DoRequest("GET", path, "", "", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(resp, &backupList)
	if err != nil {
		log.Println(err)
		return &backupList, err
	}

	return &backupList, nil
}

// SetBackupParameters request body
type SetBackupParametersRequest struct {
	// The user name of the virtual machine
	Username string `json:"guest_username,omitempty"`

	// The password of the virtual machine
	Password string `json:"guest_password,omitempty"`

	// Set to true to disable virtual machine validation logic
	OverrideValidation bool `json:"override_validation,omitempty"`

	// app_aware_type: Set the application aware backup type:
	//   VSS - Application-consistent backup using Microsoft VSS
	//   DEFAULT - Crash-consistent backup
	//   NONE - Application-consistent backup using a VMware snapshot
	AppAwareType string `json:"app_aware_type,omitempty"`
}

// SetBackupParameters sets the virtual machine backup parameters used for application consistent backups.
func (v *VirtualMachine) SetBackupParameters(req *SetBackupParametersRequest) error {
	var (
		path = fmt.Sprintf("/virtual_machines/%s/backup_parameters", v.Id)
	)

	header := map[string]string{"Content-Type": "application/vnd.simplivity.v1.11+json"}
	resp, err := commonClient.DoRequest("POST", path, "", req, header)
	if err != nil {
		log.Println(err)
		return err
	}

	task, err := commonClient.Tasks.WaitForTask(resp)
	if err != nil {
		log.Println(err)
		return err
	}

	resources := task.AffectedResources
	if len(resources) < 1 {
		err_message := "Set backup parameters operation was not successful. Error code:" + string(task.ErrorCode)
		return errors.New(err_message)
	}

	return nil
}

// UpdatePowerState sets power state of the VM.
func (v *VirtualMachine) UpdatePowerState(state string) error {

	var path string

	if state == "off" {
		path = fmt.Sprintf("/virtual_machines/%s/power_off", v.Id)
	} else if state == "on" {
		path = fmt.Sprintf("/virtual_machines/%s/power_on", v.Id)
	} else {
		error_message := "Pass a valid power state"
		return errors.New(error_message)
	}

	req_header := map[string]string{"Content-Type": "application/vnd.simplivity.v1.11+json"}
	resp, err := commonClient.DoRequest("POST", path, "", "", req_header)
	if err != nil {
		log.Println(err)
		return err
	}

	task, err := commonClient.Tasks.WaitForTask(resp)
	if err != nil {
		log.Println(err)
		return err
	}

	if len(task.AffectedResources) < 1 {
		err_message := "Setting power state operation was not successful. Error code:" + string(task.ErrorCode)
		return errors.New(err_message)
	}

	return nil
}
