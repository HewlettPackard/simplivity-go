package ovc

import (
	"encoding/json"
	"errors"
	"log"
)

// OmniStackClusterResource handles communications with the the Cluster resource methods
//
// SimpliVity API docs: https://developer.hpe.com/api/simplivity/endpoint?&path=%2Fomnistack_clusters
type OmniStackClusterResource resourceClient

// GetAll response fields
type OmniStackClusterList struct {
	Offset  int                 `json:"offset,omitempty"`
	Count   int                 `json:"count,omitempty"`
	Limit   int                 `json:"limit,omitempty"`
	Members []*OmniStackCluster `json:"omnistack_clusters,omitempty"`
}

type InfosightConfiguration struct {
	InfosightRegistered bool   `json:"infosight_registered,omitempty"`
	InfosightEnabled    bool   `json:"infosight_enabled,omitempty"`
	InfosightProxyURL   string `json:"infosight_proxy_url,omitempty"`
}

// OmniStackCluster represents a SimpliVity OmniStack cluster.
type OmniStackCluster struct {
	Name                           string                 `json:"name,omitempty"`
	Id                             string                 `json:"id,omitempty"`
	HypervisorObjectParentName     string                 `json:"hypervisor_object_parent_name,omitempty"`
	ClusterFeatureLevel            int                    `json:"cluster_feature_level,omitempty"`
	ClusterGroupIds                []string               `json:"cluster_group_ids,omitempty"`
	ArbiterConnected               bool                   `json:"arbiter_connected,omitempty"`
	HypervisorObjectParentID       string                 `json:"hypervisor_object_parent_id,omitempty"`
	Type                           string                 `json:"type,omitempty"`
	Version                        string                 `json:"version,omitempty"`
	HypervisorObjectID             string                 `json:"hypervisor_object_id,omitempty"`
	Members                        []string               `json:"members,omitempty"`
	ArbiterAddress                 string                 `json:"arbiter_address,omitempty"`
	HypervisorType                 string                 `json:"hypervisor_type,omitempty"`
	HypervisorManagementSystem     string                 `json:"hypervisor_management_system,omitempty"`
	HypervisorManagementSystemName string                 `json:"hypervisor_management_system_name. omitempty"`
	InfosightConfiguration         InfosightConfiguration `json:"infosight_configuration,omitempty"`
	IwoEnabled                     bool                   `json:"iwo_enabled,omitempty"`
}

// GetAll returns all the OmniStack Clusters filtered by the query parameters.
// Filters:
//   id: The unique identifier (UID) of the omnistack_clusters to return
//     Accepts: Single value, comma-separated list
//   name: The name of the omnistack_clusters to return
//     Accepts: Single value, comma-separated list, pattern using one or more
//     asterisk characters as a wildcard
//   hypervisor_object_id: The unique identifier (UID) of the hypervisor associated
//     with the objects to return
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   hypervisor_object_parent_id: The unique identifier (UID) of the hypervisor that
//     contains the objects to return
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   hypervisor_object_parent_name: The name of the hypervisor that contains the objects
//     to return
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   hypervisor_management_system_name: The name of the hypervisor associated with the
//     omnistack_cluster
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   type: The type of omnistack_clusters to return
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   arbiter_address: The address of the Arbiter connected to the objects to return
//     Accepts: Single value, comma-separated list, pattern using one or more asterisk
//     characters as a wildcard
//   arbiter_connected: An indicator to show if the omnistack_cluster is connected to Arbiter
//     Valid values:
//     True: Only returns omnistack_clusters connected to Arbiters that you identified
//       in arbiter_address
//     False: Only returns omnistack_clusters not connected to Arbiters that you identified
//       in arbiter_address
func (o *OmniStackClusterResource) GetAll(params GetAllParams) (*OmniStackClusterList, error) {
	var (
		path        = "/omnistack_clusters"
		clusterList OmniStackClusterList
	)

	qrStr := params.QueryString()

	resp, err := o.client.DoRequest("GET", path, qrStr, nil, nil)
	if err != nil {
		return &clusterList, err
	}

	err = json.Unmarshal(resp, &clusterList)
	if err != nil {
		return &clusterList, err
	}

	return &clusterList, nil
}

// GetBy searches for OmniStack Clusters with single filter.
func (o *OmniStackClusterResource) GetBy(field string, value string) ([]*OmniStackCluster, error) {
	filters := map[string]string{field: value}
	clusterList, err := o.GetAll(GetAllParams{Filters: filters})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	clusters := clusterList.Members

	return clusters, nil
}

// GetByName searches for an OmniStack Cluster by its name
func (o *OmniStackClusterResource) GetByName(name string) (*OmniStackCluster, error) {
	clusters, err := o.GetBy("name", name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(clusters) > 0 {
		cluster := clusters[0]
		return cluster, nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// GetById searches for an OmniStack Cluster by its id.
func (o *OmniStackClusterResource) GetById(id string) (*OmniStackCluster, error) {
	clusters, err := o.GetBy("id", id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(clusters) > 0 {
		cluster := clusters[0]
		return cluster, nil
	}

	return nil, errors.New("Resource doesn't exist")
}
