package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"sandbox/internal/dto"
	"sandbox/internal/handler"
	"sandbox/internal/response"
	"sandbox/internal/router"
	"sandbox/internal/sandbox/mock"
	"sandbox/internal/service"
	"sandbox/internal/workspace"
)

func TestSandboxAPIWorkflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Setup temporary directory for host workspaces
	tempDir := t.TempDir()

	// 2. Initialize dependencies with Mock Sandbox Manager
	workspaceManager := workspace.NewManager(tempDir)
	mockSandboxManager := mock.NewManager()
	sandboxService := service.NewSandboxService(mockSandboxManager, workspaceManager, tempDir)
	sandboxHandler := handler.NewSandboxHandler(sandboxService)

	// 3. Initialize router
	r := router.NewRouter(sandboxHandler)

	var sandboxID string

	// --- Test Scenario 1: Create Sandbox ---
	t.Run("Create Sandbox", func(t *testing.T) {
		reqBody := dto.SandboxCreateRequest{
			Scenario: "python:3.10-alpine",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/sandbox", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
		}

		var createResp response.SandboxCreateResponse
		if err := json.Unmarshal(resp.Body.Bytes(), &createResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if createResp.SandboxID == "" {
			t.Fatal("Expected SandboxID to be non-empty")
		}
		sandboxID = createResp.SandboxID
	})

	if sandboxID == "" {
		t.Fatal("Abort subsequent tests: sandboxID was not generated")
	}

	// --- Test Scenario 2: Write File ---
	t.Run("Write File", func(t *testing.T) {
		reqBody := dto.FileWriteRequest{
			Filename: "test.py",
			Content:  "print('hello from api test')",
		}
		jsonBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/sandbox/"+sandboxID+"/files", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", resp.Code, resp.Body.String())
		}
	})

	// --- Test Scenario 3: List Files ---
	t.Run("List Files", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/sandbox/"+sandboxID+"/files", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.Code)
		}

		var listResp response.FileListResponse
		if err := json.Unmarshal(resp.Body.Bytes(), &listResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// By default "main.py" is created, and we added "test.py"
		foundTest := false
		foundMain := false
		for _, f := range listResp.Files {
			if f == "test.py" {
				foundTest = true
			}
			if f == "main.py" {
				foundMain = true
			}
		}

		if !foundTest {
			t.Error("Expected test.py to be in file list")
		}
		if !foundMain {
			t.Error("Expected main.py to be in file list")
		}
	})

	// --- Test Scenario 4: Read File Content ---
	t.Run("Read File Content", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/sandbox/"+sandboxID+"/files/content?filename=test.py", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.Code)
		}

		var fileResp response.FileContentResponse
		if err := json.Unmarshal(resp.Body.Bytes(), &fileResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if fileResp.Filename != "test.py" {
			t.Errorf("Expected filename test.py, got %s", fileResp.Filename)
		}
		if fileResp.Content != "print('hello from api test')" {
			t.Errorf("Expected content match, got %s", fileResp.Content)
		}
	})

	// --- Test Scenario 5: Run Command ---
	t.Run("Run Command", func(t *testing.T) {
		reqBody := dto.SandboxRunRequest{
			Cmd: []string{"python", "main.py"},
		}
		jsonBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/sandbox/"+sandboxID+"/run", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.Code)
		}

		var runResp response.SandboxRunResponse
		if err := json.Unmarshal(resp.Body.Bytes(), &runResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if runResp.Output == "" {
			t.Error("Expected non-empty output from command run")
		}
	})

	// --- Test Scenario 6: Get Logs ---
	t.Run("Get Logs", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/sandbox/"+sandboxID+"/logs", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.Code)
		}

		var logsResp response.SandboxRunResponse
		if err := json.Unmarshal(resp.Body.Bytes(), &logsResp); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if logsResp.Output == "" {
			t.Error("Expected logs output to be non-empty")
		}
	})

	// --- Test Scenario 7: Stop Sandbox ---
	t.Run("Stop Sandbox", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/sandbox/"+sandboxID+"/stop", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.Code)
		}
	})

	// --- Test Scenario 8: Delete File ---
	t.Run("Delete File", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/sandbox/"+sandboxID+"/files?filename=test.py", nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.Code)
		}

		// Double check it's deleted
		checkReq := httptest.NewRequest("GET", "/sandbox/"+sandboxID+"/files", nil)
		checkResp := httptest.NewRecorder()
		r.ServeHTTP(checkResp, checkReq)

		var listResp response.FileListResponse
		_ = json.Unmarshal(checkResp.Body.Bytes(), &listResp)
		for _, f := range listResp.Files {
			if f == "test.py" {
				t.Error("Expected test.py to be deleted from workspaces")
			}
		}
	})

	// --- Test Scenario 9: Delete Sandbox ---
	t.Run("Delete Sandbox", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/sandbox/"+sandboxID, nil)
		resp := httptest.NewRecorder()

		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.Code)
		}
	})
}
