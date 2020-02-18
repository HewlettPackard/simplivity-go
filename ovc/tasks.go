package ovc

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// TaskResource handles communications with the SimpliVity task resource.
//
// SimpliVity API docs: https://developer.hpe.com/api/simplivity/endpoint?&path=%2Ftasks
type TaskResource resourceClient

// Task fields
type Task struct {
	State             string              `json:"state,omitempty"`
	Id                string              `json:"id,omitempty"`
	Progress          int                 `json:"percent_complete,omitempty"`
	AffectedResources []*AffectedResource `json:"affected_objects,omitempty"`
	ErrorCode         int                 `json:"error_code,omitempty"`
}

// List of affected resources in task response
type AffectedResource struct {
	ObjectType string `json:"object_type,omitempty"`
	ObjectId   string `json:"object_id,omitempty"`
}

// Task endpoint response
type TaskResp struct {
	Task *Task `json:"task,omitempty"`
}

// WaitForTask waits for the task to complete
// Makes continous calls to the server using CheckProgress method 
// and checks the status of the task
func (s *TaskResource) WaitForTask(resp []byte) (*Task, error) {
	var (
		taskResp TaskResp
		task     *Task
		err      error
	)

	for {
		err = json.Unmarshal(resp, &taskResp)
		if err != nil {
			break
		}

		log.Println(taskResp.Task.State)
		if taskResp.Task.State != "IN_PROGRESS" {
			task = taskResp.Task
			break
		}
                 
		// Sleep for two seconds if the request is in progress
		// To avoid hitting the server continously
		time.Sleep(2 * time.Second)

		resp, err = s.CheckProgress(taskResp.Task)
		if err != nil {
			break
		}
	}

	return task, err
}

// CheckProgress makes call to the server for task status.
func (s *TaskResource) CheckProgress(task *Task) ([]byte, error) {
	var (
		path = fmt.Sprintf("/tasks/%s", task.Id)
	)

	resp, err := s.client.DoRequest("GET", path, "", nil, nil)

	return resp, err
}
