package ovc

import (
	"encoding/json"
	"errors"
	"log"
)

// PolicyResource handles communications with the the Policy resource methods
//
// SimpliVity API docs: https://developer.hpe.com/api/simplivity/endpoint?&path=%2Fpolicies
type PolicyResource resourceClient

type PolicyList struct {
	Offset  int       `json:"offset,omitempty"`
	Count   int       `json:"count,omitempty"`
	Limit   int       `json:"limit,omitempty"`
	Members []*Policy `json:"policies,omitempty"`
}

// Policy represents a SimpliVity Policy resource.
type Policy struct {
	Name            string        `json:"name,omitempty"`
	Id              string        `json:"id,omitempty"`
	ClusterGroupIds []string      `json:"cluster_group_ids,omitempty"`
	Rules           []interface{} `json:"rules,omitempty"`
}

// GetAll returns all the policies filtered by the query parameters.
// Filters:
//   id: The unique identifier (UID) of the policy
//     Accepts: Single value, comma-separated list
//   name:The name of the policy
//     Accepts: Single value, comma-separated list
func (p *PolicyResource) GetAll(params GetAllParams) (*PolicyList, error) {
	var (
		path       = "/policies"
		policyList PolicyList
	)

	qrStr := params.QueryString()

	resp, err := p.client.DoRequest("GET", path, qrStr, nil, nil)
	if err != nil {
		return &policyList, err
	}

	err = json.Unmarshal(resp, &policyList)
	if err != nil {
		return &policyList, err
	}

	return &policyList, nil
}

// GetBy searches for Policies with single filter.
func (p *PolicyResource) GetBy(field string, value string) ([]*Policy, error) {
	filters := map[string]string{field: value}
	policyList, err := p.GetAll(GetAllParams{Filters: filters})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	policies := policyList.Members

	return policies, nil
}

// GetByName searches for a Policy resource by its name.
func (p *PolicyResource) GetByName(name string) (*Policy, error) {
	policies, err := p.GetBy("name", name)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(policies) > 0 {
		policy := policies[0]
		return policy, nil
	}

	return nil, errors.New("Resource doesn't exist")
}

// GetById searches for a Policy resource by its id.
func (p *PolicyResource) GetById(id string) (*Policy, error) {
	policies, err := p.GetBy("id", id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(policies) > 0 {
		policy := policies[0]
		return policy, nil
	}

	return nil, errors.New("Resource doesn't exist")
}
