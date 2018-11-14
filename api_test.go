package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleListRequests(t *testing.T) {
	{
		req := httptest.NewRequest(http.MethodGet, "/requests", nil)
		recorder := httptest.NewRecorder()

		HandleListRequests(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, recorder.Code)
		}

		var res []RequestResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &res)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}

	{
		req := httptest.NewRequest(
			http.MethodGet,
			"/requests?start=2018-11-11T22%3A43%3A04%2B01%3A00&end=2018-11-12T22%3A43%3A04%2B01%3A00&limit=1",
			nil,
		)
		recorder := httptest.NewRecorder()

		HandleListRequests(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("expected %d, got %d", http.StatusOK, recorder.Code)
		}
		var res []RequestResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &res)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}

	{
		req := httptest.NewRequest(
			http.MethodGet,
			"/requests?start=toto",
			nil,
		)
		recorder := httptest.NewRecorder()

		HandleListRequests(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	}

	{
		req := httptest.NewRequest(
			http.MethodGet,
			"/requests?end=toto",
			nil,
		)
		recorder := httptest.NewRecorder()

		HandleListRequests(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	}

	{
		req := httptest.NewRequest(
			http.MethodGet,
			"/requests?limit=toto",
			nil,
		)
		recorder := httptest.NewRecorder()

		HandleListRequests(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	}

	{
		req := httptest.NewRequest(
			http.MethodPost,
			"/requests",
			bytes.NewBuffer([]byte(`{"start":"2018-11-11T23:30:03+01:00","end":"2018-11-11T23:30:03+01:00"}`)),
		)
		recorder := httptest.NewRecorder()

		HandleListRequests(recorder, req)

		if recorder.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected %d, got %d", http.StatusMethodNotAllowed, recorder.Code)
		}
	}
}
