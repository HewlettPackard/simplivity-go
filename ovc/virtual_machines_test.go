package ovc

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func mockGetVMById(apiHandler *http.ServeMux) {
	apiHandler.HandleFunc("/virtual_machines", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "virtual_machines":[{"name": "testvm", "id":"1"}]}`)
	})
}

func TestVMGetAllWithDefaultParameters(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{}

	apiHandler.HandleFunc("/virtual_machines", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "virtual_machines":[{"name": "testvm", "id":"1"}]}`)
	})

	vm_list, err := client.VirtualMachines.GetAll(params)
	if err != nil {
		t.Error(err)
	}

	vm := &VirtualMachine{Name: "testvm", Id: "1"}
	expected := &VirtualMachineList{Offset: 1,
		Limit:   1,
		Count:   1,
		Members: []*VirtualMachine{vm}}

	if !reflect.DeepEqual(vm_list, expected) {
		t.Errorf("Returned = %v, expected %v", vm_list, expected)
	}
}

func TestVMGetAllWithFilter(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()
	params := GetAllParams{Filters: map[string]string{"name": "testname"}}

	apiHandler.HandleFunc("/virtual_machines", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 1, "count": 1, "virtual_machines":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.VirtualMachines.GetAll(params)
	if err != nil {
		t.Error(err)
	}
}

func TestVMGetBy(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/virtual_machines", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "virtual_machines":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.VirtualMachines.GetBy("name", "testname")
	if err != nil {
		t.Error(err)
	}
}

func TestVMGetByName(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/virtual_machines", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "name": "testname"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "virtual_machines":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.VirtualMachines.GetByName("testname")
	if err != nil {
		t.Error(err)
	}
}

func TestVMGetById(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/virtual_machines", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmVal := formValues{"limit": "500", "offset": "0", "sort": "",
			"order": "", "fields": "", "case": "",
			"show_optional_fields": "false", "id": "123"}

		testFormValues(t, r, fmVal)
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "virtual_machines":[{"name": "testvm", "id":"1"}]}`)
	})

	_, err := client.VirtualMachines.GetById("123")
	if err != nil {
		t.Error(err)
	}
}

func TestSetPolicyForMultipleVMs(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	apiHandler.HandleFunc("/virtual_machines/set_policy", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestBody(t, r, `{"policy_id":"1","virtual_machine_id":["1"]}`+"\n")
		fmt.Fprint(w, `{"task":{"state": "IN_PROGRESS", "id": "1", "percent_complete": 1,
		    "affected_objects":[]}}`)
	})

	//Mock request to task endpoint
	mockTaskRequest(apiHandler)

	policy := &Policy{Id: "1"}
	vm := &VirtualMachine{Id: "1"}
	vms := []*VirtualMachine{vm}

	err := client.VirtualMachines.SetPolicyForMultipleVMs(policy, vms)
	if err != nil {
		t.Error(err)
	}
}

func TestSetPolicyForSingleVM(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	//Mock get by id request
	mockGetVMById(apiHandler)

	apiHandler.HandleFunc("/virtual_machines/1/set_policy", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestBody(t, r, `{"policy_id":"1"}`+"\n")
		fmt.Fprint(w, `{"task":{"state": "IN_PROGRESS", "id": "1", "percent_complete": 1,
		    "affected_objects":[]}}`)
	})

	//Mock request to task endpoint
	mockTaskRequest(apiHandler)

	vm, err := client.VirtualMachines.GetById("1")
	if err != nil {
		t.Error(err)
	}

	policy := &Policy{Id: "1"}

	err = vm.SetPolicy(policy)
	if err != nil {
		t.Error(err)
	}
}

func TestClone(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	//Mock get by id request
	mockGetVMById(apiHandler)

	apiHandler.HandleFunc("/virtual_machines/1/clone", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestBody(t, r, `{"app_consistent":false,"virtual_machine_name":"testvm"}`+"\n")
		fmt.Fprint(w, `{"task":{"state": "IN_PROGRESS", "id": "1", "percent_complete": 1, "affected_objects":[]}}`)
	})

	//Mock request to task endpoint
	mockTaskRequest(apiHandler)

	vm, err := client.VirtualMachines.GetById("1")
	if err != nil {
		t.Error(err)
	}

	_, err = vm.Clone("testvm", false)
	if err != nil {
		t.Error(err)
	}
}

func TestMove(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	//Mock get by id request
	mockGetVMById(apiHandler)

	apiHandler.HandleFunc("/virtual_machines/1/move", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestBody(t, r, `{"destination_datastore_id":"1","virtual_machine_name":"testvm"}`+"\n")
		fmt.Fprint(w, `{"task":{"state": "IN_PROGRESS", "id": "1", "percent_complete": 1, "affected_objects":[]}}`)
	})

	//Mock request to task endpoint
	mockTaskRequest(apiHandler)

	vm, err := client.VirtualMachines.GetById("1")
	if err != nil {
		t.Error(err)
	}

	d := &Datastore{Id: "1"}

	_, err = vm.Move("testvm", d)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateBackup(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	//Mock get by id request
	mockGetVMById(apiHandler)

	apiHandler.HandleFunc("/backups", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "backups":[{"name": "backup_name", "id":"1"}]}`)
	})

	apiHandler.HandleFunc("/virtual_machines/1/backup", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestBody(t, r, `{"backup_name":"backup_name","destination_id":"1"}`+"\n")
		fmt.Fprint(w, `{"task":{"state": "IN_PROGRESS", "id": "1", "percent_complete": 1, "affected_objects":[]}}`)
	})

	//Mock request to task endpoint
	mockTaskRequest(apiHandler)

	vm, err := client.VirtualMachines.GetById("1")
	if err != nil {
		t.Error(err)
	}

	cluster := &OmniStackCluster{Id: "1"}
	req := &CreateBackupRequest{Name: "backup_name"}
	backup, err := vm.CreateBackup(req, cluster)
	if err != nil {
		t.Error(err)
	}

	expected := &Backup{Name: "backup_name", Id: "1"}
	if !reflect.DeepEqual(backup, expected) {
		t.Errorf("Returned = %v, expected %v", backup, expected)
	}
}

func TestGetBackups(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	//Mock get by id request
	mockGetVMById(apiHandler)

	apiHandler.HandleFunc("/virtual_machines/1/backups", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "GET")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		fmt.Fprint(w, `{"offset": 1, "limit": 500, "count": 1, "backups":[{"name": "testbackup", "id":"1"}]}`)
	})

	vm, err := client.VirtualMachines.GetById("1")
	if err != nil {
		t.Error(err)
	}

	backups, err := vm.GetBackups()
	if err != nil {
		t.Error(err)
	}

	backup := &Backup{Name: "testbackup", Id: "1"}
	expected := &BackupList{Offset: 1,
		Limit:   500,
		Count:   1,
		Members: []*Backup{backup}}

	if !reflect.DeepEqual(backups, expected) {
		t.Errorf("Returned = %v, expected %v", backups, expected)
	}
}

func TestSetBackupParameters(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	//Mock get by id request
	mockGetVMById(apiHandler)

	apiHandler.HandleFunc("/virtual_machines/1/backup_parameters", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestHeader(t, r, "Content-Type", "application/vnd.simplivity.v1.11+json")
		testRequestBody(t, r, `{"guest_username":"username","guest_password":"password","app_aware_type":"NONE"}`+"\n")
		fmt.Fprint(w, `{"task":{"state": "IN_PROGRESS", "id": "1", "percent_complete": 1, "affected_objects":[]}}`)
	})

	//Mock request to task endpoint
	mockTaskRequest(apiHandler)

	vm, err := client.VirtualMachines.GetById("1")
	if err != nil {
		t.Error(err)
	}

	req := &SetBackupParametersRequest{Username: "username", Password: "password", AppAwareType: "NONE"}
	err = vm.SetBackupParameters(req)
	if err != nil {
		t.Error(err)
	}
}
