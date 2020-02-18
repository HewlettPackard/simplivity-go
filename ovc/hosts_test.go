package ovc

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestHostGetAllWithDefaultParameters(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{}

	apiHandler.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "hosts":[{"name": "test", "id":"1"}]}`)
	})

	resource_list, err := client.Hosts.GetAll(params)
	if err != nil {
		t.Error(err)
	}

	resource := &Host{Name: "test", Id: "1"}
	expected := &HostList{Offset: 1,
		Limit:   1,
		Count:   1,
		Members: []*Host{resource}}

	if !reflect.DeepEqual(resource_list, expected) {
		t.Errorf("Returned = %v, expected %v", resource_list, expected)
	}
}

func TestHostGetAllWithFilter(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{Filters: map[string]string{"name": "testname"}}

	apiHandler.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "hosts":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Hosts.GetAll(params)
	if err != nil {
		t.Error(err)
	}
}

func TestHostGetBy(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "hosts":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Hosts.GetBy("name", "testname")
	if err != nil {
		t.Error(err)
	}
}

func TesHostGetByName(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "hosts":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Hosts.GetByName("testname")
	if err != nil {
		t.Error(err)
	}
}

func TestHostGetById(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "id": "123"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "hosts":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Hosts.GetById("123")
	if err != nil {
		t.Error(err)
	}
}
