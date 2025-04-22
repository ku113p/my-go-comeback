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
	"github.com/stretchr/testify/require"
)

func setupTestServer() *httptest.Server {
	mux := getRootMux()
	return httptest.NewServer(mux)
}

func makeRequest(method, url, body string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func decodeJSON(t *testing.T, body io.Reader, target interface{}) {
	err := json.NewDecoder(body).Decode(target)
	require.NoError(t, err)
}

func TestPingRoute(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/ping")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	require.NoError(t, err)

	assert.Equal(t, "pong", string(body))
}

func TestNoteRoutes(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	var createdNote models.Note
	t.Run("Create", func(t *testing.T) {
		payload := `{"name": "My First Note", "description": "This is the first note", "author_id": "` + uuid.New().String() + `"}`
		resp, err := makeRequest(http.MethodPost, server.URL+"/note", payload)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		decodeJSON(t, resp.Body, &createdNote)
		assert.NotNil(t, createdNote.ID)
	})

	t.Run("Get", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/note/" + createdNote.ID.String())
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var fetched models.Note
		decodeJSON(t, resp.Body, &fetched)
		assert.Equal(t, createdNote, fetched)
	})

	t.Run("List", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/note")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var list map[string][]models.Note
		decodeJSON(t, resp.Body, &list)
		assert.GreaterOrEqual(t, len(list["results"]), 1)
	})

	t.Run("Update", func(t *testing.T) {
		payload := `{"name": "Updated Note", "description": "Updated Desc"}`
		resp, err := makeRequest(http.MethodPut, server.URL+"/note/"+createdNote.ID.String(), payload)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated models.Note
		decodeJSON(t, resp.Body, &updated)
		assert.Equal(t, "Updated Note", updated.Name)
		assert.Equal(t, "Updated Desc", *updated.Description)
	})

	t.Run("Delete", func(t *testing.T) {
		resp, err := makeRequest(http.MethodDelete, server.URL+"/note/"+createdNote.ID.String(), "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		resp, err = http.Get(server.URL + "/note/" + createdNote.ID.String())
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestAuthorRoutes(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	var created models.Author
	t.Run("Create", func(t *testing.T) {
		payload := `{"username": "testuser", "firstname": "John", "secondname": "Doe"}`
		resp, err := makeRequest(http.MethodPost, server.URL+"/author", payload)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		decodeJSON(t, resp.Body, &created)
	})

	t.Run("Get", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/author/" + created.ID.String())
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var fetched models.Author
		decodeJSON(t, resp.Body, &fetched)
		assert.Equal(t, created, fetched)
	})

	t.Run("List", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/author")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var list map[string][]models.Author
		decodeJSON(t, resp.Body, &list)
		assert.GreaterOrEqual(t, len(list["results"]), 1)
	})

	t.Run("Update", func(t *testing.T) {
		payload := `{"username": "updateduser", "secondname": "Smith"}`
		resp, err := makeRequest(http.MethodPut, server.URL+"/author/"+created.ID.String(), payload)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated models.Author
		decodeJSON(t, resp.Body, &updated)
		assert.Equal(t, "updateduser", updated.Username)
		assert.Equal(t, "Smith", *updated.Secondname)
	})

	t.Run("Delete", func(t *testing.T) {
		resp, err := makeRequest(http.MethodDelete, server.URL+"/author/"+created.ID.String(), "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestInvalidUUID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	checkError := func(method, path string, body string) {
		resp, err := makeRequest(method, server.URL+path, body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp map[string]string
		decodeJSON(t, resp.Body, &errResp)
		assert.Equal(t, "bad_request", errResp["error"])
		assert.Equal(t, "invalid uuid", errResp["message"])
	}

	checkError(http.MethodGet, "/note/invalid-uuid", "")
	checkError(http.MethodDelete, "/note/invalid-uuid", "")
	checkError(http.MethodPut, "/note/invalid-uuid", `{"name": "test"}`)
}

func TestNotFound(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	id := uuid.New().String()

	checkNotFound := func(method, path, body string) {
		resp, err := makeRequest(method, server.URL+path, body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var errResp map[string]string
		decodeJSON(t, resp.Body, &errResp)
		assert.Equal(t, "not_found", errResp["error"])
	}

	checkNotFound(http.MethodGet, "/note/"+id, "")
	checkNotFound(http.MethodPut, "/note/"+id, `{"name": "test"}`)
	checkNotFound(http.MethodDelete, "/note/"+id, "")
}

func TestMethodNotAllowed(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	checkMethodNotAllowed := func(method, path string) {
		req, err := http.NewRequest(method, server.URL+path, strings.NewReader(`{}`))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	}

	checkMethodNotAllowed(http.MethodPut, "/note")
	checkMethodNotAllowed(http.MethodDelete, "/note")
	checkMethodNotAllowed(http.MethodPost, "/note/"+uuid.New().String())
}

func TestBadRequestOnCreateUpdate(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	t.Run("CreateNoteInvalidPayload", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/note", "application/json", strings.NewReader(`{"description":"Missing name"}`))
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("UpdateNoteInvalidPayload", func(t *testing.T) {
		// First, create a valid note
		valid := models.Note{Name: "Test", AuthorId: &models.ObjectID{}}
		var buf bytes.Buffer
		require.NoError(t, json.NewEncoder(&buf).Encode(valid))
		resp, err := http.Post(server.URL+"/note", "application/json", &buf)
		require.NoError(t, err)

		var note models.Note
		decodeJSON(t, resp.Body, &note)

		// Then, try to update it with an invalid payload
		resp, err = makeRequest(http.MethodPut, server.URL+"/note/"+note.ID.String(), `{"author_id":"not-a-uuid"}`)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("CreateAuthorInvalidPayload", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/author", "application/json", strings.NewReader(`{"firstname":"John"}`))
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("UpdateAuthorInvalidPayload", func(t *testing.T) {
		author := models.Author{Username: "test"}
		var buf bytes.Buffer
		require.NoError(t, json.NewEncoder(&buf).Encode(author))
		resp, err := http.Post(server.URL+"/author", "application/json", &buf)
		require.NoError(t, err)

		var created models.Author
		decodeJSON(t, resp.Body, &created)

		resp, err = makeRequest(http.MethodPut, server.URL+"/author/"+created.ID.String(), `{"username":123}`)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
