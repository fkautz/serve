// Copyright (c) 2014-2025 Frederick F. Kautz IV
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/handlers"
)

func TestFileServer(t *testing.T) {
	// Create a temporary directory with a test file
	tmpDir := t.TempDir()
	testContent := []byte("Hello, World!")
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create the file server handler
	server := handlers.CompressHandler(http.FileServer(http.Dir(tmpDir)))

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/test.txt", nil)
	w := httptest.NewRecorder()

	// Serve the request
	server.ServeHTTP(w, req)

	// Check the response
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if string(body) != string(testContent) {
		t.Errorf("expected body %q, got %q", testContent, body)
	}
}

func TestFileServerNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	server := handlers.CompressHandler(http.FileServer(http.Dir(tmpDir)))

	req := httptest.NewRequest(http.MethodGet, "/nonexistent.txt", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestFileServerWithCompression(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a larger file to trigger compression
	testContent := make([]byte, 1024)
	for i := range testContent {
		testContent[i] = 'a'
	}
	testFile := filepath.Join(tmpDir, "large.txt")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	server := handlers.CompressHandler(http.FileServer(http.Dir(tmpDir)))

	req := httptest.NewRequest(http.MethodGet, "/large.txt", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check that compression was applied
	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("expected gzip encoding, got %q", resp.Header.Get("Content-Encoding"))
	}
}
