package ovc

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestClusterGetAllWithDefaultParameters(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{}

	apiHandler.HandleFunc("/omnistack_clusters", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "omnistack_clusters":[{"name": "test", "id":"1"}]}`)
	})

	resource_list, err := client.OmniStackClusters.GetAll(params)
	if err != nil {
		t.Error(err)
	}

	resource := &OmniStackCluster{Name: "test", Id: "1"}
	expected := &OmniStackClusterList{Offset: 1,
		Limit:   1,
		Count:   1,
		Members: []*OmniStackCluster{resource}}

	if !reflect.DeepEqual(resource_list, expected) {
		t.Errorf("Returned = %v, expected %v", resource_list, expected)
	}
}

func TestClusterGetAllWithFilter(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{Filters: map[string]string{"name": "testname"}}

	apiHandler.HandleFunc("/omnistack_clusters", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "omnistack_clusters":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.OmniStackClusters.GetAll(params)
	if err != nil {
		t.Error(err)
	}
}

func TestClusterGetBy(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/omnistack_clusters", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "omnistack_clusters":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.OmniStackClusters.GetBy("name", "testname")
	if err != nil {
		t.Error(err)
	}
}

func TestClusterGetByName(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/omnistack_clusters", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "omnistack_clusters":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.OmniStackClusters.GetByName("testname")
	if err != nil {
		t.Error(err)
	}
}

func TestClusterGetById(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/omnistack_clusters", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "id": "123"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "omnistack_clusters":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.OmniStackClusters.GetById("123")
	if err != nil {
		t.Error(err)
	}
}
