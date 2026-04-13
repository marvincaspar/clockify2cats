package report

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTestRepository(server *httptest.Server) Repository {
	return Repository{
		WorkspaceID: "ws1",
		UserID:      "user1",
		ApiKey:      "test-api-key",
		BaseURL:     server.URL,
		HTTPClient:  server.Client(),
	}
}

func TestRepository_FetchClockifyData_success(t *testing.T) {
	entries := []ClockifyTimeEntry{
		{
			Description: "Task A",
			TimeInterval: struct {
				Start    string `json:"start"`
				End      string `json:"end"`
				Duration string `json:"duration"`
			}{Start: "2022-01-03T08:00:00.000Z", End: "2022-01-03T09:00:00.000Z", Duration: "PT1H"},
			Project: struct {
				Name string `json:"name"`
			}{Name: "Project (123)"},
			Billable: true,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(entries)
	}))
	defer server.Close()

	repo := makeTestRepository(server)
	result, err := repo.FetchClockifyData("2022-01-03T00:00:00.000Z")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Task A", result[0].Description)
	assert.Equal(t, "Project (123)", result[0].Project.Name)
	assert.True(t, result[0].Billable)
}

func TestRepository_FetchClockifyData_setsApiKeyHeader(t *testing.T) {
	var capturedKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedKey = r.Header.Get("X-Api-Key")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	repo := makeTestRepository(server)
	_, err := repo.FetchClockifyData("2022-01-03T00:00:00.000Z")

	assert.NoError(t, err)
	assert.Equal(t, "test-api-key", capturedKey)
}

func TestRepository_FetchClockifyData_buildsCorrectURL(t *testing.T) {
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.String()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	repo := makeTestRepository(server)
	_, err := repo.FetchClockifyData("2022-01-03T00:00:00.000Z")

	assert.NoError(t, err)
	assert.Contains(t, capturedPath, "/api/v1/workspaces/ws1/user/user1/time-entries")
	assert.Contains(t, capturedPath, "hydrated=1")
	assert.Contains(t, capturedPath, "page-size=1000")
	assert.Contains(t, capturedPath, "start=2022-01-03T00:00:00Z")
	assert.Contains(t, capturedPath, "end=2022-01-09T23:59:59Z")
}

func TestRepository_FetchClockifyData_nonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	repo := makeTestRepository(server)
	_, err := repo.FetchClockifyData("2022-01-03T00:00:00.000Z")

	assert.EqualError(t, err, "Clockify API error: 401 Unauthorized")
}

func TestRepository_FetchClockifyData_networkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	repo := makeTestRepository(server)
	server.Close() // close before the request is made

	_, err := repo.FetchClockifyData("2022-01-03T00:00:00.000Z")
	assert.Error(t, err)
}

func TestRepository_FetchClockifyData_invalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not valid json"))
	}))
	defer server.Close()

	repo := makeTestRepository(server)
	_, err := repo.FetchClockifyData("2022-01-03T00:00:00.000Z")

	assert.ErrorContains(t, err, "error parsing Clockify response")
}
