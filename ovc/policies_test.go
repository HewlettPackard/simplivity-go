package ovc

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPolicyGetAllWithDefaultParameters(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{}

	apiHandler.HandleFunc("/policies", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "policies":[{"name": "test", "id":"1"}]}`)
	})

	resource_list, err := client.Policies.GetAll(params)
	if err != nil {
		t.Error(err)
	}

	resource := &Policy{Name: "test", Id: "1"}
	expected := &PolicyList{Offset: 1,
		Limit:   1,
		Count:   1,
		Members: []*Policy{resource}}

	if !reflect.DeepEqual(resource_list, expected) {
		t.Errorf("Returned = %v, expected %v", resource_list, expected)
	}
}

func TestPolicyGetAllWithFilter(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{Filters: map[string]string{"name": "testname"}}

	apiHandler.HandleFunc("/policies", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "policies":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Policies.GetAll(params)
	if err != nil {
		t.Error(err)
	}
}

func TestPolicyGetBy(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/policies", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "policies":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Policies.GetBy("name", "testname")
	if err != nil {
		t.Error(err)
	}
}

func TestPolicyGetByName(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/policies", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "policies":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Policies.GetByName("testname")
	if err != nil {
		t.Error(err)
	}
}

func TestPolicyGetById(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/policies", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "id": "123"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "policies":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Policies.GetById("123")
	if err != nil {
		t.Error(err)
	}
}
