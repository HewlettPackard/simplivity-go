package ovc

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

// BackupResource handles communications with the the Backup resource methods
//
// SimpliVity API docs: https://developer.hpe.com/api/simplivity/endpoint?&path=%2Fbackups
type BackupResource resourceClient

type BackupList struct {
	Offset  int       `json:"offset,omitempty"`
	Count   int       `json:"count,omitempty"`
	Limit   int       `json:"limit,omitempty"`
	Members []*Backup `json:"backups,omitempty"`
}

// Backup represents a backup of the VM resource in SimpliVity.
type Backup struct {
	Name                                   string    `json:"name,omitempty"`
	Id                                     string    `json:"id,omitempty"`
	VirtualMachineName                     string    `json:"virtual_machine_name,omitempty"`
	CreatedAt                              time.Time `json:"created_at,omitempty"`
	ConsistencyType                        string    `json:"consistency_type,omitempty"`
	Type                                   string    `json:"type,omitempty"`
	DatastoreName                          string    `json:"datastore_name,omitempty"`
	VirtualMachineDeletionTime             int       `json:"virtual_machine_deletion_time,omitempty"`
	VirtualMachineID                       string    `json:"virtual_machine_id,omitempty"`
	ApplicationConsistent                  bool      `json:"application_consistent,omitempty"`
	ComputeClusterParentHypervisorObjectID string    `json:"compute_cluster_parent_hypervisor_object_id,omitempty"`
	State                                  string    `json:"state,omitempty"`
	OmnistackClusterID                     string    `json:"omnistack_cluster_id,omitempty"`
	VirtualMachineType                     string    `json:"virtual_machine_type,omitempty"`
	SentCompletionTime                     string    `json:"sent_completion_time,omitempty"`
	UniqueSizeBytes                        int       `json:"unique_size_bytes,omitempty"`
	ClusterGroupIds                        []string  `json:"cluster_group_ids,omitempty"`
	UniqueSizeTimestamp                    string    `json:"unique_size_timestamp,omitempty"`
	ExpirationTime                         string    `json:"expiration_time,omitempty"`
	OmnistackClusterName                   string    `json:"omnistack_cluster_name,omitempty"`
	Sent                                   int       `json:"sent,omitempty"`
	Size                                   int64     `json:"size,omitempty"`
	VirtualMachineState                    string    `json:"virtual_machine_state,omitempty"`
	DatastoreID                            string    `json:"datastore_id,omitempty"`
	ComputeClusterParentName               string    `json:"compute_cluster_parent_name,omitempty"`
	HypervisorType                         string    `json:"hypervisor_type,omitempty"`
	SentDuration                           int       `json:"sent_duration,omitempty"`
}

// GetAll returns all the backups filtered by the query parameters.
// Filters:
//   id: the unique identifier (UID) of the backups to return
//     Accepts: Single value, comma-separated list
//   name: The name of the backups to return
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard.
//   sent_min: The minimum sent data size (in bytes) of the remote backups to return.
//   sent_max: The maximum sent data size (in bytes) of the remote backups to return.
//   state: The current state of the backups to return
//     Accepts: Single value, comma-separated list
//   type: The type of backups to return.
//     Accepts: Single value, comma-separated list.
//   omnistack_cluster_id: The unique identifier (UID) of the omnistack_cluster
//     that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   omnistack_cluster_name: The name of the omnistack_cluster that is associated
//     with the instances to return
//     Accepts: Single value, comma-separated list
//   compute_cluster_parent_hypervisor_object_id: The unique identifier (UID) of the
//     hypervisor that contains the omnistack_cluster that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   compute_cluster_parent_name: The name of the hypervisor that contains the
//     omnistack_cluster that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   datastore_id: The unique identifier (UID) of the datastore that is associated
//     with the instances to return
//     Accepts: Single value, comma-separated list
//   datastore_name: The name of the datastore that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   expires_before: The latest expiration time before the backups to return expire,
//     expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   expires_after: The earliest expiration time after the backups to return expire,
//     expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   virtual_machine_id: The unique identifier (UID) of the virtual_machine that is
//     associated with the instances to return
//     Accepts: Single value, comma-separated list
//   virtual_machine_name: The name of the virtual_machine that is associated with
//     the instances to return
//     Accepts: Single value, comma-separated list
//   virtual_machine_type: The type of the virtual_machine that is associated with the
//     instances to return
//     Accepts: Single value, comma-separated list
//   size_min: The minimum size (in bytes) of the backups to return
//   size_max: The maximum size (in bytes) of the backups to return
//   application_consistent: The application-consistent setting of the backups to return
//   consistency_type: The consistency type of the backups to return
//     Accepts: Single value, comma-separated list
//   created_before: The latest creation time before the backups to return were created,
//     expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   created_after: The earliest creation time after the backups to return were created,
//     expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   sent_duration_min: The minimum number of seconds that elapsed while replicating
//     the backups to return
//   sent_duration_max: The maximum number of seconds that elapsed while replicating the
//     backups to return
//   sent_completion_before: The latest time before the replication of backups to return was
//     completed, expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   sent_completion_after: The earliest time after the replication of backups to return was
//     completed, expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
func (b *BackupResource) GetAll(params GetAllParams) (*BackupList, error) {
	var (
		path       = "/backups"
		backupList BackupList
	)

	qrStr := params.QueryString()

	resp, err := b.client.DoRequest("GET", path, qrStr, nil, nil)
	if err != nil {
		return &backupList, err
	}

	err = json.Unmarshal(resp, &backupList)
	if err != nil {
		log.Println(err)
		return &backupList, err
	}

	return &backupList, nil
}

// GetBy searches for backups with single filter.
func (b *BackupResource) GetBy(field string, value string) ([]*Backup, error) {
	filters := map[string]string{field: value}
	BackupList, err := b.GetAll(GetAllParams{Filters: filters})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	backups := BackupList.Members

	return backups, nil
}

// GetByName searches for a backup by name.
func (b *BackupResource) GetByName(name string) (*Backup, error) {
	backups, err := b.GetBy("name", name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(backups) > 0 {
		backup := backups[0]
		return backup, nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// GetById searches for a backup by id.
func (b *BackupResource) GetById(id string) (*Backup, error) {
	backups, err := b.GetBy("id", id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(backups) > 0 {
		backup := backups[0]
		return backup, nil
	}

	return nil, errors.New("Resource doesn't exist")
}
