package ovc

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

// DatastoreResource handles communications with the the Datastore resource methods
//
// SimpliVity API docs: https://developer.hpe.com/api/simplivity/endpoint?&path=%2Fdatastores
type DatastoreResource resourceClient

type DatastoreList struct {
	Offset  int          `json:"offset,omitempty"`
	Count   int          `json:"count,omitempty"`
	Limit   int          `json:"limit,omitempty"`
	Members []*Datastore `json:"datastores,omitempty"`
}

type Shares struct {
	Address string `json:"address,omitempty"`
	Host    string `json:"host,omitempty"`
	Rw      bool   `json:"rw,omitempty"`
}

// Datastore represents a SimpliVity datastore
type Datastore struct {
	Name                                   string    `json:"name,omitempty"`
	Id                                     string    `json:"id,omitempty"`
	ClusterGroupIds                        []string  `json:"cluster_group_ids,omitempty"`
	PolicyID                               string    `json:"policy_id,omitempty"`
	MountDirectory                         string    `json:"mount_directory,omitempty"`
	CreatedAt                              time.Time `json:"created_at,omitempty"`
	PolicyName                             string    `json:"policy_name,omitempty"`
	OmnistackClusterName                   string    `json:"omnistack_cluster_name,omitempty"`
	Deleted                                bool      `json:"deleted,omitempty"`
	HypervisorObjectID                     string    `json:"hypervisor_object_id,omitempty"`
	Size                                   int64     `json:"size,omitempty"`
	ComputeClusterParentHypervisorObjectID string    `json:"compute_cluster_parent_hypervisor_object_id,omitempty"`
	ComputeClusterParentName               string    `json:"compute_cluster_parent_name,omitempty"`
	HypervisorType                         string    `json:"hypervisor_type,omitempty"`
	OmnistackClusterID                     string    `json:"omnistack_cluster_id,omitempty"`
	HypervisorManagementSystem             string    `json:"hypervisor_management_system,omitempty"`
	HypervisorManagementSystemName         string    `json:"hypervisor_management_system_name,omitempty"`
	Shares                                 []Shares  `json:"shares,omitempty"`
}

// GetAll returns all the datastores filtered by the query parameters.
// Filters:
//   id: The unique identifier (UID) of the datastores to return
//     Accepts: Single value, comma-separated list
//   name: The name of the datastores to return
//     Accepts: Single value, comma-separated list, pattern using one
//     or more asterisk characters as a wildcard
//   min_size: The minimum size (in bytes) of datastores to return
//   max_size: The maximum size (in bytes) of datastores to return
//   created_before: The latest creation time before the datastores to return were created,
//     expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   created_after: The earliest creation time after the datastores to return were created,
//     expressed in ISO-8601 form, based on Coordinated Universal Time (UTC)
//   omnistack_cluster_id: The unique identifier (UID) of the omnistack_cluster that is
//     associated with the instances to return
//     Accepts: Single value, comma-separated list
//   omnistack_cluster_name: The name of the omnistack_cluster that is associated with
//     the instances to return
//     Accepts: Single value, comma-separated list
//   compute_cluster_parent_hypervisor_object_id: The unique identifier (UID) of the hypervisor
//     that contains the omnistack_cluster that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   compute_cluster_parent_name: The name of the hypervisor that contains the omnistack
//     cluster that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   hypervisor_management_system_name: The name of the Hypervisor Management System (HMS)
//     associated with the datastore
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   policy_id: The unique identifier (UID) of the policy that is associated with the
//     instances to return
//     Accepts: Single value, comma-separated list
//   policy_name: The name of the policy that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   hypervisor_object_id: The unique identifier (UID) of the hypervisor-based instance
//     that is associated with the instances to return
//     Accepts: Single value, comma-separated list
//   mount_directory: A comma-separated list of fields to include in the returned objects
//     Default: Returns all fields
func (d *DatastoreResource) GetAll(params GetAllParams) (*DatastoreList, error) {
	var (
		path          = "/datastores"
		datastoreList DatastoreList
	)

	qrStr := params.QueryString()

	resp, err := d.client.DoRequest("GET", path, qrStr, nil, nil)
	if err != nil {
		return &datastoreList, err
	}

	err = json.Unmarshal(resp, &datastoreList)
	if err != nil {
		return &datastoreList, err
	}

	return &datastoreList, nil
}

// GetBy gets datastores with single filter
func (d *DatastoreResource) GetBy(field string, value string) ([]*Datastore, error) {
	filters := map[string]string{field: value}
	datastoreList, err := d.GetAll(GetAllParams{Filters: filters})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	datastores := datastoreList.Members

	return datastores, nil
}

// GetByName searches for a datastore by its name
func (d *DatastoreResource) GetByName(name string) (*Datastore, error) {
	datastores, err := d.GetBy("name", name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(datastores) > 0 {
		return datastores[0], nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// GetById searches for a datastore by its id.
func (d *DatastoreResource) GetById(id string) (*Datastore, error) {
	datastores, err := d.GetBy("id", id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(datastores) > 0 {
		return datastores[0], nil
	}

	return nil, errors.New("Resource doesn't exist")
}
