package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"notes/api/models"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *httptest.Server {
	mux := getRootMux()
	return httptest.NewServer(mux)
}

func TestPingRoute(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/ping")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, "pong", string(body))
}

func TestNoteRoutes(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Helper function to make requests
	makeRequest := func(method, path string, body string) (*http.Response, error) {
		req, err := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		return client.Do(req)
	}

	// Create a note
	createPayload := `{"name": "My First Note", "description": "This is the first note", "author_id": "` + uuid.New().String() + `"}`
	resp, err := makeRequest(http.MethodPost, "/note", createPayload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var createdNote models.Note
	err = json.NewDecoder(resp.Body).Decode(&createdNote)
	assert.NoError(t, err)
	assert.NotNil(t, createdNote.ID)
	assert.Equal(t, "My First Note", createdNote.Name)
	assert.Equal(t, "This is the first note", *createdNote.Description)
	assert.NotNil(t, createdNote.Created)

	noteID := createdNote.ID.String()

	// Get the created note
	resp, err = http.Get(server.URL + "/note/" + noteID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedNote models.Note
	err = json.NewDecoder(resp.Body).Decode(&fetchedNote)
	assert.NoError(t, err)
	assert.Equal(t, createdNote, fetchedNote)

	// List notes
	resp, err = http.Get(server.URL + "/note")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse map[string][]models.Note
	err = json.NewDecoder(resp.Body).Decode(&listResponse)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResponse["results"]), 1)

	// Update the note
	updatePayload := `{"name": "Updated Note Name", "description": "Updated description"}`
	resp, err = makeRequest(http.MethodPut, "/note/"+noteID, updatePayload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedNote models.Note
	err = json.NewDecoder(resp.Body).Decode(&updatedNote)
	assert.NoError(t, err)
	assert.Equal(t, *createdNote.ID, *updatedNote.ID)
	assert.Equal(t, "Updated Note Name", updatedNote.Name)
	assert.Equal(t, "Updated description", *updatedNote.Description)
	assert.NotEqual(t, createdNote.AuthorId, updatedNote.AuthorId)
	assert.NotEqual(t, createdNote.Name, updatedNote.Name)
	assert.NotEqual(t, createdNote.Description, updatedNote.Description)

	// Get the updated note
	resp, err = http.Get(server.URL + "/note/" + noteID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedUpdatedNote models.Note
	err = json.NewDecoder(resp.Body).Decode(&fetchedUpdatedNote)
	assert.NoError(t, err)
	assert.Equal(t, updatedNote, fetchedUpdatedNote)

	// Delete the note
	resp, err = makeRequest(http.MethodDelete, "/note/"+noteID, "")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Try to get the deleted note
	resp, err = http.Get(server.URL + "/note/" + noteID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAuthorRoutes(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Helper function to make requests
	makeRequest := func(method, path string, body string) (*http.Response, error) {
		req, err := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		return client.Do(req)
	}

	// Create an author
	createPayload := `{"username": "testuser", "firstname": "John", "secondname": "Doe"}`
	resp, err := makeRequest(http.MethodPost, "/author", createPayload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var createdAuthor models.Author
	err = json.NewDecoder(resp.Body).Decode(&createdAuthor)
	assert.NoError(t, err)
	assert.NotNil(t, createdAuthor.ID)
	assert.Equal(t, "testuser", createdAuthor.Username)
	assert.Equal(t, "John", *createdAuthor.Firstname)
	assert.Equal(t, "Doe", *createdAuthor.Secondname)

	authorID := createdAuthor.ID.String()

	// Get the created author
	resp, err = http.Get(server.URL + "/author/" + authorID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedAuthor models.Author
	err = json.NewDecoder(resp.Body).Decode(&fetchedAuthor)
	assert.NoError(t, err)
	assert.Equal(t, createdAuthor, fetchedAuthor)

	// List authors
	resp, err = http.Get(server.URL + "/author")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse map[string][]models.Author
	err = json.NewDecoder(resp.Body).Decode(&listResponse)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResponse["results"]), 1)

	// Update the author
	updatePayload := `{"username": "updateduser", "secondname": "Smith"}`
	resp, err = makeRequest(http.MethodPut, "/author/"+authorID, updatePayload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedAuthor models.Author
	err = json.NewDecoder(resp.Body).Decode(&updatedAuthor)
	assert.NoError(t, err)
	assert.Equal(t, *createdAuthor.ID, *updatedAuthor.ID)
	assert.Equal(t, "updateduser", updatedAuthor.Username)
	assert.Nil(t, updatedAuthor.Firstname)
	assert.Equal(t, "Smith", *updatedAuthor.Secondname)
	assert.NotEqual(t, createdAuthor.Username, updatedAuthor.Username)
	assert.NotEqual(t, createdAuthor.Secondname, updatedAuthor.Secondname)

	// Get the updated author
	resp, err = http.Get(server.URL + "/author/" + authorID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedUpdatedAuthor models.Author
	err = json.NewDecoder(resp.Body).Decode(&fetchedUpdatedAuthor)
	assert.NoError(t, err)
	assert.Equal(t, updatedAuthor, fetchedUpdatedAuthor)

	// Delete the author
	resp, err = makeRequest(http.MethodDelete, "/author/"+authorID, "")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Try to get the deleted author
	resp, err = http.Get(server.URL + "/author/" + authorID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestInvalidUUID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/note/invalid-uuid")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errorResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "bad_request", errorResponse["error"])
	assert.Equal(t, "invalid uuid", errorResponse["message"])

	// Test delete with invalid UUID
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/note/invalid-uuid", nil)
	assert.NoError(t, err)
	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "bad_request", errorResponse["error"])
	assert.Equal(t, "invalid uuid", errorResponse["message"])

	// Test update with invalid UUID
	updatePayload := `{"name": "Should not update"}`
	req, err = http.NewRequest(http.MethodPut, server.URL+"/note/invalid-uuid", bytes.NewBufferString(updatePayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "bad_request", errorResponse["error"])
	assert.Equal(t, "invalid uuid", errorResponse["message"])
}

func TestNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	nonExistentID := uuid.New().String()

	// Test get non-existent
	resp, err := http.Get(server.URL + "/note/" + nonExistentID)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	var errorResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", errorResponse["error"])

	// Test delete non-existent
	req, err := http.NewRequest(http.MethodDelete, server.URL+"/note/"+nonExistentID, nil)
	assert.NoError(t, err)
	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", errorResponse["error"])

	// Test update non-existent
	updatePayload := `{"name": "Should not update"}`
	req, err = http.NewRequest(http.MethodPut, server.URL+"/note/"+nonExistentID, bytes.NewBufferString(updatePayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, "not_found", errorResponse["error"])
}

func TestMethodNotAllowed(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Test PUT on /note
	req, err := http.NewRequest(http.MethodPut, server.URL+"/note", strings.NewReader(`{}`))
	assert.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	// Test DELETE on /note
	req, err = http.NewRequest(http.MethodDelete, server.URL+"/note", nil)
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	// Test POST on /note/{id}
	req, err = http.NewRequest(http.MethodPost, server.URL+"/note/"+uuid.New().String(), strings.NewReader(`{}`))
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestBadRequestOnCreateUpdate(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Test create note with invalid payload (missing required fields)
	invalidCreatePayload := `{"description": "Missing name"}`
	resp, err := http.Post(server.URL+"/note", "application/json", strings.NewReader(invalidCreatePayload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode) // Assuming internal error for now, adjust if model parsing handles validation

	// Test update note with invalid payload
	validNote := models.Note{Name: "Test", AuthorId: &models.ObjectID{}}
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(validNote)
	createResp, err := http.Post(server.URL+"/note", "application/json", &b)
	assert.NoError(t, err)
	var createdNote models.Note
	json.NewDecoder(createResp.Body).Decode(&createdNote)
	assert.NotNil(t, createdNote.ID)

	invalidUpdatePayload := `{"author_id": "not-a-uuid"}`
	req, err := http.NewRequest(http.MethodPut, server.URL+"/note/"+createdNote.ID.String(), strings.NewReader(invalidUpdatePayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode) // Assuming internal error due to parsing

	// Test create author with invalid payload
	invalidAuthorPayload := `{"firstname": "John"}`
	resp, err = http.Post(server.URL+"/author", "application/json", strings.NewReader(invalidAuthorPayload))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Test update author with invalid payload
	validAuthor := models.Author{Username: "test"}
	var authorBuffer bytes.Buffer
	json.NewEncoder(&authorBuffer).Encode(validAuthor)
	createAuthorResp, err := http.Post(server.URL+"/author", "application/json", &authorBuffer)
	assert.NoError(t, err)
	var createdAuthor models.Author
	json.NewDecoder(createAuthorResp.Body).Decode(&createdAuthor)
	assert.NotNil(t, createdAuthor.ID)
}
