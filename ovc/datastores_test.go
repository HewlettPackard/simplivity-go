package ovc

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDatastoreGetAllWithDefaultParameters(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{}

	apiHandler.HandleFunc("/datastores", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "datastores":[{"name": "test", "id":"1"}]}`)
	})

	resource_list, err := client.Datastores.GetAll(params)
	if err != nil {
		t.Error(err)
	}

	resource := &Datastore{Name: "test", Id: "1"}
	expected := &DatastoreList{Offset: 1,
		Limit:   1,
		Count:   1,
		Members: []*Datastore{resource}}

	if !reflect.DeepEqual(resource_list, expected) {
		t.Errorf("Returned = %v, expected %v", resource_list, expected)
	}
}

func TestDatastoreGetAllWithFilter(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{Filters: map[string]string{"name": "testname"}}

	apiHandler.HandleFunc("/datastores", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "datastores":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Datastores.GetAll(params)
	if err != nil {
		t.Error(err)
	}
}

func TestDatastoreGetBy(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/datastores", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "datastores":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Datastores.GetBy("name", "testname")
	if err != nil {
		t.Error(err)
	}
}

func TestDatastoreGetByName(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/datastores", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "datastores":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Datastores.GetByName("testname")
	if err != nil {
		t.Error(err)
	}
}

func TestDatastoreGetById(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/datastores", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "id": "123"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "datastores":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Datastores.GetById("123")
	if err != nil {
		t.Error(err)
	}
}
