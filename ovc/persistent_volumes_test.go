package ovc

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func mockGetPVById(apiHandler *http.ServeMux) {
	apiHandler.HandleFunc("/persistent_volumes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "persistent_volumes":[{"name": "pvc-123_fcd", "id":"1"}]}`)
	})
}

func TestPVGetAllWithDefaultParameters(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{}

	apiHandler.HandleFunc("/persistent_volumes", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "persistent_volumes":[{"name": "pvc-123_fcd", "id":"1"}]}`)
	})

	pv_list, err := client.PersistentVolumes.GetAll(params)
	if err != nil {
		t.Error(err)
	}

	pv := &PersistentVolume{Name: "pvc-123_fcd", Id: "1"}
	expected := &PersistentVolumeList{Offset: 1,
		Limit:   1,
		Count:   1,
		Members: []*PersistentVolume{pv}}

	if !reflect.DeepEqual(pv_list, expected) {
		t.Errorf("Returned = %v, expected %v", pv_list, expected)
	}
}

func TestPVGetAllWithFilter(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{Filters: map[string]string{"name": "testname"}}

	apiHandler.HandleFunc("/persistent_volumes", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "persistent_volumes":[{"name": "pvc-123_fcd", "id":"1"}]}`)
	})

	_, err := client.PersistentVolumes.GetAll(params)
	if err != nil {
		t.Error(err)
	}
}

func TestPVGetBy(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/persistent_volumes", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "persistent_volumes":[{"name": "pvc-123_fcd", "id":"1"}]}`)
	})

	_, err := client.PersistentVolumes.GetBy("name", "testname")
	if err != nil {
		t.Error(err)
	}
}

func TestPVGetByName(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/persistent_volumes", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "persistent_volumes":[{"name": "pvc-123_fcd", "id":"1"}]}`)
	})

	_, err := client.PersistentVolumes.GetByName("testname")
	if err != nil {
		t.Error(err)
	}
}

func TestPVGetById(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/persistent_volumes", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "id": "123"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "persistent_volumes":[{"name": "pvc-123_fcd", "id":"1"}]}`)
	})

	_, err := client.PersistentVolumes.GetById("123")
	if err != nil {
		t.Error(err)
	}
}

func TestSetPolicyForMultiplePVs(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/persistent_volumes/set_policy", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestBody(t, r, `{"persistent_volume_id":["1"],"policy_id":"1"}`+"\n")
		fmt.Fprint(w, `{"task":{"state": "IN_PROGRESS", "id": "1", "percent_complete": 1,
		    "affected_objects":[]}}`)
	})

	//Mock request to task endpoint
	mockTaskRequest(apiHandler)

	policy := &Policy{Id: "1"}
	pv := &PersistentVolume{Id: "1"}
	pvs := []*PersistentVolume{pv}

	err := client.PersistentVolumes.SetPolicyForMultiplePVs(policy, pvs)
	if err != nil {
		t.Error(err)
	}

	// Negative test
	pvs = nil
	err = client.PersistentVolumes.SetPolicyForMultiplePVs(policy, pvs)
	if err.Error() != "Pass a list of PV resoures" {
		t.Error(err)
	}
}

func TestPVCreateBackup(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	//Mock get by id request
	mockGetPVById(apiHandler)

	apiHandler.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "backups":[{"name": "backup_name", "id":"1"}]}`)
	})

	apiHandler.HandleFunc("/persistent_volumes/1/backup", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestBody(t, r, `{"backup_name":"backup_name","destination_id":"1"}`+"\n")
		fmt.Fprint(w, `{"task":{"state": "IN_PROGRESS", "id": "1", "percent_complete": 1, "affected_objects":[]}}`)
	})

	//Mock request to task endpoint
	mockTaskRequest(apiHandler)

	pv, err := client.PersistentVolumes.GetById("1")
	if err != nil {
		t.Error(err)
	}

	cluster := &OmniStackCluster{Id: "1"}
	req := &CreateBackupRequest{Name: "backup_name"}
	backup, err := pv.CreateBackup(req, cluster)
	if err != nil {
		t.Error(err)
	}

	expected := &Backup{Name: "backup_name", Id: "1"}
	if !reflect.DeepEqual(backup, expected) {
		t.Errorf("Returned = %v, expected %v", backup, expected)
	}
}

func TestPVGetBackups(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	//Mock get by id request
	mockGetPVById(apiHandler)

	apiHandler.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "backups":[{"name": "backup_name", "id":"1"}]}`)
	})

	pv, err := client.PersistentVolumes.GetById("1")
	if err != nil {
		t.Error(err)
	}

	backups, err := pv.GetBackups()
	if err != nil {
		t.Error(err)
	}

	backup := &Backup{Name: "backup_name", Id: "1"}
	expected := &BackupList{Offset: 1,
		Limit:   500,
		Count:   1,
		Members: []*Backup{backup}}

	if !reflect.DeepEqual(backups, expected) {
		t.Errorf("Returned = %v, expected %v", backups, expected)
	}
}
