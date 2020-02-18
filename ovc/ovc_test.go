package ovc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func setup() (client *Client, mux *http.ServeMux, teardown func()) {
	apiHandler := http.NewServeMux()

	apiHandler.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"access_token": "12345"}`)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	client, err := NewClient("user", "pass", server.URL, "")
	if err != nil {
		fmt.Println("Test setup Error", err)
	}

	return client, apiHandler, server.Close
}

func mockTaskRequest(apiHandler *http.ServeMux) {
	apiHandler.HandleFunc("/tasks/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"task": {"state": "COMPLETE", "id": "1", "percent_complete": 1, "affected_objects":[{"object_type":"", "object_id":"1"}]}}`)
	})
}

func testRequestMethod(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	if got := r.Method; got != expected {
		t.Errorf("Request method: %v, expected %v", got, expected)
	}
}

type formValues map[string]string

func testFormValues(t *testing.T, r *http.Request, values formValues) {
	t.Helper()
	expected := url.Values{}
	for k, v := range values {
		expected.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, expected) {
		t.Errorf("Request parameters: %v, expected %v", got, expected)
	}
}

func testRequestHeader(t *testing.T, r *http.Request, header string, expected string) {
	t.Helper()
	if got := r.Header.Get(header); got != expected {
		t.Errorf("Header.Get(%q) returned %q, expected %q", header, got, expected)
	}
}

func testRequestBody(t *testing.T, r *http.Request, expected string) {
	t.Helper()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != expected {
		t.Errorf("request Body is %s, expected %s", got, expected)
	}
}

func TestSetAccessToken(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	client.SetAccessToken()
	if client.AccessToken != "12345" {
		t.Errorf("AccessToken is not set")
	}
}

func TestDoRequest(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	type testBody struct {
		Test string
	}
	body := new(testBody)
	headers := map[string]string{"Test": "test"}
	queryStr := "a=1"

	apiHandler.HandleFunc("/test/path", func(w http.ResponseWriter, r *http.Request) {
		testRequestMethod(t, r, "POST")
		testRequestHeader(t, r, "Authorization", "Bearer 12345")
		testRequestBody(t, r, `{"Test":""}`+"\n")
		testFormValues(t, r, formValues{"a": "1"})
		fmt.Fprint(w, `{"Test":"a"}`)
	})

	resp, err := client.DoRequest("POST", "/test/path", queryStr, body, headers)
	if err != nil {
		t.Error(err)
	}

	json.Unmarshal(resp, &body)
	expected := &testBody{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestDoRequestTokenExpiry(t *testing.T) {
	client, apiHandler, teardown := setup()
	defer teardown()

	headers := map[string]string{"Test": "test"}
	apiHandler.HandleFunc("/test/path", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": "invalid_token","message": "Access token expired due to inactivity: 4a1da31f-5405-4a93-af5c-799403ea70d6"}`)
	})

	_, err := client.DoRequest("POST", "/test/path", "", "", headers)
	errorMessage := "Error: Status code: - Access token expired due to inactivity: 4a1da31f-5405-4a93-af5c-799403ea70d6"
	if err != nil && err.Error() != errorMessage {
		t.Error("Call should return invalid_token error")
	}
}

func TestCommonOVCClient(t *testing.T) {
	client, _, teardown := setup()
	defer teardown()

	if client.VirtualMachines.client != client.Policies.client {
		t.Errorf("Resource clients are not using a common OVC client.")
	}
}
