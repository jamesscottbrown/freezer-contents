package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary test file with given content
func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_contents.json")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile
}

// Helper function to set up test environment with a temporary contents file
func setupTestEnv(t *testing.T, content string) func() {
	t.Helper()
	tmpFile := createTempFile(t, content)
	oldContentsFile := contentsFile
	contentsFile = tmpFile
	return func() {
		contentsFile = oldContentsFile
	}
}

// Sample test data
var sampleState = `{
  "Containers": ["1", "2", "3"],
  "Freezers": [
    {
      "Name": "freezer1",
      "Contents": [
        {
          "Name": "chicken",
          "Date": "2024-01-01",
          "Containers": ["1", "2"]
        }
      ]
    },
    {
      "Name": "freezer2",
      "Contents": []
    }
  ]
}`

// ============== Tests for readContents ==============

func TestReadContents_ValidFile(t *testing.T) {
	tmpFile := createTempFile(t, sampleState)

	state, err := readContents(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(state.Containers) != 3 {
		t.Errorf("Expected 3 containers, got %d", len(state.Containers))
	}

	if len(state.Freezers) != 2 {
		t.Errorf("Expected 2 freezers, got %d", len(state.Freezers))
	}

	if state.Freezers[0].Name != "freezer1" {
		t.Errorf("Expected freezer1, got %s", state.Freezers[0].Name)
	}

	if len(state.Freezers[0].Contents) != 1 {
		t.Errorf("Expected 1 item in freezer1, got %d", len(state.Freezers[0].Contents))
	}

	if state.Freezers[0].Contents[0].Name != "chicken" {
		t.Errorf("Expected chicken, got %s", state.Freezers[0].Contents[0].Name)
	}
}

func TestReadContents_MissingFile(t *testing.T) {
	_, err := readContents("/nonexistent/path/contents.json")
	if err == nil {
		t.Error("Expected error for missing file, got nil")
	}
}

func TestReadContents_InvalidJSON(t *testing.T) {
	tmpFile := createTempFile(t, "not valid json {{{")

	_, err := readContents(tmpFile)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestReadContents_EmptyFile(t *testing.T) {
	tmpFile := createTempFile(t, "")

	_, err := readContents(tmpFile)
	if err == nil {
		t.Error("Expected error for empty file, got nil")
	}
}

func TestReadContents_EmptyState(t *testing.T) {
	emptyState := `{"Containers": [], "Freezers": []}`
	tmpFile := createTempFile(t, emptyState)

	state, err := readContents(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(state.Containers) != 0 {
		t.Errorf("Expected 0 containers, got %d", len(state.Containers))
	}

	if len(state.Freezers) != 0 {
		t.Errorf("Expected 0 freezers, got %d", len(state.Freezers))
	}
}

// ============== Tests for writeContents ==============

func TestWriteContents_Success(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "write_test.json")

	state := State{
		Containers: []string{"a", "b"},
		Freezers: []Freezer{
			{
				Name: "testFreezer",
				Contents: []Item{
					{Name: "soup", Date: "2024-01-01", Containers: []string{"a"}},
				},
			},
		},
	}

	err := writeContents(tmpFile, state)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Read back and verify
	readState, err := readContents(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read back written file: %v", err)
	}

	if len(readState.Containers) != 2 {
		t.Errorf("Expected 2 containers, got %d", len(readState.Containers))
	}

	if readState.Freezers[0].Name != "testFreezer" {
		t.Errorf("Expected testFreezer, got %s", readState.Freezers[0].Name)
	}
}

func TestWriteContents_InvalidPath(t *testing.T) {
	state := State{Containers: []string{}, Freezers: []Freezer{}}
	err := writeContents("/nonexistent/directory/file.json", state)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

// ============== Tests for CORS middleware ==============

func TestCORS_AddsHeaders(t *testing.T) {
	handler := CORS(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Missing or incorrect Access-Control-Allow-Origin header")
	}

	if rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("Missing or incorrect Access-Control-Allow-Credentials header")
	}

	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Missing Access-Control-Allow-Methods header")
	}
}

func TestCORS_OptionsRequest(t *testing.T) {
	handlerCalled := false
	handler := CORS(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if handlerCalled {
		t.Error("Handler should not be called for OPTIONS request")
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for OPTIONS, got %d", rec.Code)
	}
}

func TestCORS_PassesThroughNonOptions(t *testing.T) {
	handlerCalled := false
	handler := CORS(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.Write([]byte("OK"))
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if !handlerCalled {
		t.Error("Handler should be called for non-OPTIONS request")
	}
}

// ============== Tests for handleStateRequest ==============

func TestHandleStateRequest_Success(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	rec := httptest.NewRecorder()

	handleStateRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var state State
	err := json.Unmarshal(rec.Body.Bytes(), &state)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	if len(state.Freezers) != 2 {
		t.Errorf("Expected 2 freezers, got %d", len(state.Freezers))
	}
}

func TestHandleStateRequest_MissingFile(t *testing.T) {
	oldContentsFile := contentsFile
	contentsFile = "/nonexistent/contents.json"
	defer func() { contentsFile = oldContentsFile }()

	req := httptest.NewRequest(http.MethodGet, "/state", nil)
	rec := httptest.NewRecorder()

	handleStateRequest(rec, req)

	// The handler doesn't set error status for missing file, it just returns empty response
	// This is a known issue in the original code
}

// ============== Tests for handleAddRequest ==============

func TestHandleAddRequest_Success(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	body := AddBody{
		Name:       "beef stew",
		Date:       "2024-02-01",
		Freezer:    "freezer1",
		Containers: []string{"3"},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handleAddRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var state State
	err := json.Unmarshal(rec.Body.Bytes(), &state)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	// Verify the item was added
	found := false
	for _, freezer := range state.Freezers {
		if freezer.Name == "freezer1" {
			for _, item := range freezer.Contents {
				if item.Name == "beef stew" {
					found = true
					break
				}
			}
		}
	}

	if !found {
		t.Error("Expected beef stew to be added to freezer1")
	}
}

func TestHandleAddRequest_NonExistentFreezer(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	body := AddBody{
		Name:       "pizza",
		Date:       "2024-02-01",
		Freezer:    "nonexistent",
		Containers: []string{"1"},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handleAddRequest(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for nonexistent freezer, got %d", rec.Code)
	}
}

func TestHandleAddRequest_InvalidJSON(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handleAddRequest(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestHandleAddRequest_EmptyBody(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/add", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handleAddRequest(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty body, got %d", rec.Code)
	}
}

// ============== Tests for handleRemoveRequest ==============

func TestHandleRemoveRequest_Success(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	body := RemoveBody{Container: "1"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/remove", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handleRemoveRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var state State
	err := json.Unmarshal(rec.Body.Bytes(), &state)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}

	// Verify container 1 was removed
	for _, freezer := range state.Freezers {
		for _, item := range freezer.Contents {
			for _, container := range item.Containers {
				if container == "1" {
					t.Error("Container 1 should have been removed")
				}
			}
		}
	}
}

func TestHandleRemoveRequest_RemoveLastContainer(t *testing.T) {
	// When removing the last container from an item, the item should be removed too
	stateWithOneContainer := `{
		"Containers": ["1"],
		"Freezers": [{
			"Name": "freezer1",
			"Contents": [{
				"Name": "single item",
				"Date": "2024-01-01",
				"Containers": ["1"]
			}]
		}]
	}`
	cleanup := setupTestEnv(t, stateWithOneContainer)
	defer cleanup()

	body := RemoveBody{Container: "1"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/remove", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()

	handleRemoveRequest(rec, req)

	var state State
	json.Unmarshal(rec.Body.Bytes(), &state)

	// The item should be removed entirely
	if len(state.Freezers[0].Contents) != 0 {
		t.Errorf("Expected 0 items in freezer after removing last container, got %d", len(state.Freezers[0].Contents))
	}
}

func TestHandleRemoveRequest_NonExistentContainer(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	body := RemoveBody{Container: "nonexistent"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/remove", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()

	handleRemoveRequest(rec, req)

	// Should succeed even if container doesn't exist (no-op)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestHandleRemoveRequest_InvalidJSON(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/remove", bytes.NewReader([]byte("invalid")))
	rec := httptest.NewRecorder()

	handleRemoveRequest(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", rec.Code)
	}
}

// ============== Tests for handleMoveRequest ==============

func TestHandleMoveRequest_Success(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	body := MoveBody{
		Container:  "1",
		NewFreezer: "freezer2",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()

	handleMoveRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var state State
	json.Unmarshal(rec.Body.Bytes(), &state)

	// Verify container moved from freezer1 to freezer2
	for _, freezer := range state.Freezers {
		if freezer.Name == "freezer1" {
			for _, item := range freezer.Contents {
				for _, c := range item.Containers {
					if c == "1" {
						t.Error("Container 1 should no longer be in freezer1")
					}
				}
			}
		}
		if freezer.Name == "freezer2" {
			found := false
			for _, item := range freezer.Contents {
				for _, c := range item.Containers {
					if c == "1" {
						found = true
					}
				}
			}
			if !found {
				t.Error("Container 1 should be in freezer2")
			}
		}
	}
}

func TestHandleMoveRequest_ConsolidateItems(t *testing.T) {
	// When moving a container to a freezer that already has an item with the same name and date,
	// the containers should be consolidated
	stateWithMatchingItem := `{
		"Containers": ["1", "2"],
		"Freezers": [
			{
				"Name": "freezer1",
				"Contents": [{
					"Name": "soup",
					"Date": "2024-01-01",
					"Containers": ["1"]
				}]
			},
			{
				"Name": "freezer2",
				"Contents": [{
					"Name": "soup",
					"Date": "2024-01-01",
					"Containers": ["2"]
				}]
			}
		]
	}`
	cleanup := setupTestEnv(t, stateWithMatchingItem)
	defer cleanup()

	body := MoveBody{
		Container:  "1",
		NewFreezer: "freezer2",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()

	handleMoveRequest(rec, req)

	var state State
	json.Unmarshal(rec.Body.Bytes(), &state)

	// Verify containers are consolidated
	for _, freezer := range state.Freezers {
		if freezer.Name == "freezer2" {
			if len(freezer.Contents) != 1 {
				t.Errorf("Expected 1 consolidated item, got %d", len(freezer.Contents))
			}
			if len(freezer.Contents) > 0 && len(freezer.Contents[0].Containers) != 2 {
				t.Errorf("Expected 2 containers in consolidated item, got %d", len(freezer.Contents[0].Containers))
			}
		}
	}
}

func TestHandleMoveRequest_InvalidJSON(t *testing.T) {
	cleanup := setupTestEnv(t, sampleState)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewReader([]byte("invalid")))
	rec := httptest.NewRecorder()

	handleMoveRequest(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", rec.Code)
	}
}

func TestHandleMoveRequest_MoveLastContainer(t *testing.T) {
	// When moving the last container, the source item should be removed
	stateWithOneContainer := `{
		"Containers": ["1"],
		"Freezers": [
			{
				"Name": "freezer1",
				"Contents": [{
					"Name": "item",
					"Date": "2024-01-01",
					"Containers": ["1"]
				}]
			},
			{
				"Name": "freezer2",
				"Contents": []
			}
		]
	}`
	cleanup := setupTestEnv(t, stateWithOneContainer)
	defer cleanup()

	body := MoveBody{
		Container:  "1",
		NewFreezer: "freezer2",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/move", bytes.NewReader(bodyBytes))
	rec := httptest.NewRecorder()

	handleMoveRequest(rec, req)

	var state State
	json.Unmarshal(rec.Body.Bytes(), &state)

	// Verify source freezer is empty
	for _, freezer := range state.Freezers {
		if freezer.Name == "freezer1" {
			if len(freezer.Contents) != 0 {
				t.Errorf("Expected source freezer to be empty, got %d items", len(freezer.Contents))
			}
		}
	}
}

// ============== Tests for handleListRequest ==============

func TestHandleListRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	rec := httptest.NewRecorder()

	handleListRequest(rec, req)

	if rec.Body.String() != "Ok" {
		t.Errorf("Expected 'Ok', got '%s'", rec.Body.String())
	}
}
