package ovc

import (
	"encoding/json"
	"errors"
	"log"
)

// HostResource handles communications with the the Host resource methods
//
// SimpliVity API docs: https://developer.hpe.com/api/simplivity/endpoint?&path=%2Fhosts
type HostResource resourceClient

// Host GetAll response
type HostList struct {
	Offset  int     `json:"offset,omitempty"`
	Count   int     `json:"count,omitempty"`
	Limit   int     `json:"limit,omitempty"`
	Members []*Host `json:"hosts,omitempty"`
}

// Host represents a SimpliVity host.
type Host struct {
	Name                                   string   `json:"name,omitempty"`
	Id                                     string   `json:"id,omitempty"`
	ClusterFeatureLevel                    int      `json:"cluster_feature_level,omitempty"`
	PolicyEnabled                          bool     `json:"policy_enabled,omitempty"`
	ComputeClusterHypervisorObjectID       string   `json:"compute_cluster_hypervisor_object_id,omitempty"`
	StorageMask                            string   `json:"storage_mask,omitempty"`
	PotentialFeatureLevel                  int      `json:"potential_feature_level,omitempty"`
	Type                                   string   `json:"type,omitempty"`
	CurrentFeatureLevel                    int      `json:"current_feature_level,omitempty"`
	HypervisorObjectID                     string   `json:"hypervisor_object_id,omitempty"`
	ComputeClusterName                     string   `json:"compute_cluster_name,omitempty"`
	ManagementIP                           string   `json:"management_ip,omitempty"`
	FederationIP                           string   `json:"federation_ip,omitempty"`
	VirtualControllerName                  string   `json:"virtual_controller_name,omitempty"`
	FederationMask                         string   `json:"federation_mask,omitempty"`
	Model                                  string   `json:"model,omitempty"`
	ComputeClusterParentHypervisorObjectID string   `json:"compute_cluster_parent_hypervisor_object_id,omitempty"`
	StorageMtu                             string   `json:"storage_mtu,omitempty"`
	OmnistackClusterID                     string   `json:"omnistack_cluster_id,omitempty"`
	State                                  string   `json:"state,omitempty"`
	UpgradeState                           string   `json:"upgrade_state,omitempty"`
	FederationMtu                          string   `json:"federation_mtu,omitempty"`
	CanRollback                            bool     `json:"can_rollback,omitempty"`
	StorageIP                              string   `json:"storage_ip,omitempty"`
	ClusterGroupIds                        []string `json:"cluster_group_ids,omitempty"`
	ManagementMtu                          string   `json:"management_mtu,omitempty"`
	Version                                string   `json:"version,omitempty"`
	ComputeClusterParentName               string   `json:"compute_cluster_parent_name,omitempty"`
	HypervisorManagementSystem             string   `json:"hypervisor_management_system,omitempty"`
	ManagementMask                         string   `json:"management_mask,omitempty"`
	HypervisorManagementSystemName         string   `json:"hypervisor_management_system_name,omitempty"`
}

// GetAll returns all the hosts filtered by the query parameters.
// Filters:
//   id: The unique identifier (UID) of the host
//     Accepts: Single value, comma-separated list
//   name: The name of the host
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   type: The type of host
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   model: The model of the host
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   version: The version of the host
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   hypervisor_management_system: The IP address of the Hypervisor Management System (HMS)
//     associated with the host
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   hypervisor_management_system_name: The name of the Hypervisor Management System (HMS)
//     associated with the host
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   hypervisor_object_id: The unique identifier (UID) of the hypervisor associated
//     with the host
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   compute_cluster_name: The name of the compute cluster associated with the host
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   compute_cluster_hypervisor_object_id: The unique identifier (UID)
//     of the Hypervisor Management System (HMS) for the associated compute cluster
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   management_ip: The IP address of the HPE OmniStack management module that
//     runs on the host
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   storage_ip: The IP address of the HPE OmniStack storage module that runs on the host
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   federation_ip: The IP address of the federation
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   virtual_controller_name: The name of the Virtual Controller that runs on the host
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   compute_cluster_parent_name: The name of the hypervisor that contains the omnistack
//     cluster that is associated with the instance
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   compute_cluster_parent_hypervisor_object_id: The unique identifier (UID) of the
//     hypervisor that contains the omnistack_cluster that is associated with the instance
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   policy_enabled: An indicator to show the status of the backup policy for the host
//     Valid values:
//     True: The backup policy for the host is enabled.
//     False: The backup policy for the host is disabled.
//   current_feature_level_min: The minimum current feature level of the HPE OmniStack
//     software running on the host
//   current_feature_level_max: The maximum current feature level of the HPE OmniStack
//     software running on the host
//   potential_feature_level_min: The minimum potential feature level of the HPE OmniStack
//     software running on the host
//   potential_feature_level_max: The maximum potential feature level of the HPE OmniStack
//     software running on the host
//   upgrade_state: The state of the most recent HPE OmniStack software upgrade for this
//     host (SUCCESS, FAIL, IN_PROGRESS, NOOP, UNKNOWN)
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   can_rollback: An indicator to show if the current HPE OmniStack software running on
//     the host can roll back to the previous version
//     Valid values:
//     True: The current HPE OmniStack software for the host can roll back to the previous version.
//     False: The current HPE OmniStack software for the host cannot roll back to the previous version.
func (h *HostResource) GetAll(params GetAllParams) (*HostList, error) {
	var (
		path     = "/hosts"
		hostList HostList
	)

	qrStr := params.QueryString()

	resp, err := h.client.DoRequest("GET", path, qrStr, nil, nil)
	if err != nil {
		return &hostList, err
	}

	err = json.Unmarshal(resp, &hostList)
	if err != nil {
		return &hostList, err
	}

	return &hostList, nil
}

// GetBy searches for hosts with single filter.
func (h *HostResource) GetBy(field string, value string) ([]*Host, error) {
	filters := map[string]string{field: value}
	hostList, err := h.GetAll(GetAllParams{Filters: filters})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	hosts := hostList.Members

	return hosts, nil
}

// GetByName searches for a host by its name
func (h *HostResource) GetByName(name string) (*Host, error) {
	hosts, err := h.GetBy("name", name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(hosts) > 0 {
		host := hosts[0]
		return host, nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// GetById searches for a host by its id
func (h *HostResource) GetById(id string) (*Host, error) {
	hosts, err := h.GetBy("id", id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(hosts) > 0 {
		host := hosts[0]
		return host, nil
	}

	return nil, errors.New("Resource doesn't exist")
}
