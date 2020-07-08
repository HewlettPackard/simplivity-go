package ovc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

// Support for persistent_volume is added in v1.16 in SimpliVity
var header = map[string]string{
	"Accept":       "application/vnd.simplivity.v1.16+json",
	"Content-Type": "application/vnd.simplivity.v1.16+json",
}

// PersistentVolumeResource handles communications with the the Persistent Volume resource methods
//
// SimpliVity API docs:<link>
type PersistentVolumeResource resourceClient

// Persistent Volumes GetAll response
type PersistentVolumeList struct {
	Offset  int                 `json:"offset,omitempty"`
	Count   int                 `json:"count,omitempty"`
	Limit   int                 `json:"limit,omitempty"`
	Members []*PersistentVolume `json:"persistent_volumes,omitempty"`
}

// PersistentVolume represents a SimpliVity persistent volume
type PersistentVolume struct {
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
	ComputeClusterParentName               string           `json:"compute_cluster_parent_name,omitempty"`
	ClusterGroupIds                        []string         `json:"cluster_group_ids,omitempty"`
	ReplicaSet                             []ReplicaSetList `json:"replica_set,omitempty"`
}

// GetAll returns all the persistent volumes filtered by the query parameters.
// Filters:
//   id: The unique identifier (UID) of the persistent volume to return
//     Accepts: Single value, comma-separated list
//   name: The name of the persistent_volumes to return
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   omnistack_cluster_id: The unique identifier (UID) of the omnistack_cluster
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
//     with the persistent volume.
//     Accepts: Single value, comma-separated list, pattern using one
//     or more asterisk characters as a wildcard
//   hypervisor_management_system_name: The name of the hypervisor associated
//     with the persistent volume
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
//   created_after: The earliest creation time after the persistent volumes to return were
//     created, expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   created_before: The latest creation time before the persistent volumes to return were
//     created, expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   state: The state of the persistent volume that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   app_aware_vm_status: The status of the ability of the persistent volume to take
//     an application-consistent backup that uses Microsoft VSS
//     Accepts: Single value, comma-separated list
//   host_id: The unique identifier (UID) of the persistent_volume host.
func (p *PersistentVolumeResource) GetAll(params GetAllParams) (*PersistentVolumeList, error) {
	var (
		path   = "/persistent_volumes"
		pvList PersistentVolumeList
	)

	qrStr := params.QueryString()
	resp, err := p.client.DoRequest("GET", path, qrStr, nil, header)
	if err != nil {
		return &pvList, err
	}

	err = json.Unmarshal(resp, &pvList)
	if err != nil {
		return &pvList, err
	}

	return &pvList, nil
}

// GetBy searches for PV resources with single filter.
func (p *PersistentVolumeResource) GetBy(fieldName string, value string) ([]*PersistentVolume, error) {
	filters := map[string]string{fieldName: value}
	pvList, err := p.GetAll(GetAllParams{Filters: filters})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pvs := pvList.Members
	return pvs, nil
}

// GetByName searches for a PV by its name
func (p *PersistentVolumeResource) GetByName(name string) (*PersistentVolume, error) {
	pvs, err := p.GetBy("name", name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(pvs) > 0 {
		pv := pvs[0]
		return pv, nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// GetById searches for a PV by its id
func (p *PersistentVolumeResource) GetById(id string) (*PersistentVolume, error) {
	pvs, err := p.GetBy("id", id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(pvs) > 0 {
		pv := pvs[0]
		return pv, nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// SetPolicyForMultiplePVs sets a policy for list of PV resources.
func (p *PersistentVolumeResource) SetPolicyForMultiplePVs(policy *Policy, pvs []*PersistentVolume) error {
	path := fmt.Sprintf("/persistent_volumes/set_policy")
	if len(pvs) < 1 {
		return errors.New("Pass a list of PV resoures")
	}

	pv_ids := []string{}
	for _, pv := range pvs {
		pv_ids = append(pv_ids, pv.Id)
	}

	body := map[string]interface{}{"policy_id": policy.Id, "persistent_volume_id": pv_ids}
	resp, err := p.client.DoRequest("POST", path, "", body, header)
	if err != nil {
		return err
	}

	_, err = commonClient.Tasks.WaitForTask(resp)
	if err != nil {
		return err
	}

	return nil
}

// CreateBackup creates a backup of the PV.
func (p *PersistentVolume) CreateBackup(req *CreateBackupRequest, dest *OmniStackCluster) (*Backup, error) {
	path := fmt.Sprintf("/persistent_volumes/%s/backup", p.Id)
	if dest != nil {
		req.Destination = dest.Id
	}

	resp, err := commonClient.DoRequest("POST", path, "", req, header)
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

// GetBackups gets all the backups of a PV.
func (p *PersistentVolume) GetBackups() (*BackupList, error) {
	backupList, err := commonClient.Backups.GetAll(GetAllParams{Filters: map[string]string{"pv": p.Name}})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return backupList, nil
}
