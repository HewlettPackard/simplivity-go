package ovc

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBackupGetAllWithDefaultParameters(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{}

	apiHandler.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "backups":[{"name": "test", "id":"1"}]}`)
	})

	resource_list, err := client.Backups.GetAll(params)
	if err != nil {
		t.Error(err)
	}

	resource := &Backup{Name: "test", Id: "1"}
	expected := &BackupList{Offset: 1,
		Limit:   1,
		Count:   1,
		Members: []*Backup{resource}}

	if !reflect.DeepEqual(resource_list, expected) {
		t.Errorf("Returned = %v, expected %v", resource_list, expected)
	}
}

func TestBackupGetAllWithFilter(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{Filters: map[string]string{"name": "testname"}}

	apiHandler.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "backups":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Backups.GetAll(params)
	if err != nil {
		t.Error(err)
	}
}

func TestBackupGetBy(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "backups":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Backups.GetBy("name", "testname")
	if err != nil {
		t.Error(err)
	}
}

func TestBackupGetByName(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "backups":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Backups.GetByName("testname")
	if err != nil {
		t.Error(err)
	}
}

func TestBackupGetById(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "id": "123"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "backups":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.Backups.GetById("123")
	if err != nil {
		t.Error(err)
	}
}
