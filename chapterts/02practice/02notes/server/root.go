package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"notes/api/models"
	"os"

	"github.com/google/uuid"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

func getRootMux() *http.ServeMux {
	mux := http.NewServeMux()

	repository := models.NewModelsRepository(models.NewUuidGenerator(), models.ModelsToRegister)

	mux.HandleFunc("/ping", ping)

	for _, name := range models.ModelsToRegister {
		getPath := func(s string) string {
			return fmt.Sprintf("/%s%s", name, s)
		}

		op := models.NewRepositoryOperation(name, repository)

		mux.HandleFunc(getPath("/{id}"), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			id := r.PathValue("id")
			switch r.Method {
			case http.MethodGet:
				get(w, id, op)
			case http.MethodDelete:
				mDelete(w, id, op)
			case http.MethodPut:
				update(w, r, id, op)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})

		mux.HandleFunc(getPath(""), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodGet:
				list(w, op)
			case http.MethodPost:
				create(w, r, op)
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})
	}

	return mux
}

func internalError(w http.ResponseWriter, err error) {
	logger.Warn("Failed response", "error", err)

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "internal_error",
	})
}

func badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "bad_request",
		"message": message,
	})
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "not_found",
	})
}

func list(w http.ResponseWriter, op *models.RepositoryOperation) {
	objects, err := op.List()
	if err != nil {
		internalError(w, err)
		return
	}

	data := map[string]any{
		"results": objects,
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		internalError(w, err)
		return
	}
}

func get(w http.ResponseWriter, id string, op *models.RepositoryOperation) {
	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "invalid uuid")
		return
	}

	obj, err := op.Get(models.ObjectID(uid))
	if err != nil {
		var notExists *models.NotExistsError
		if errors.As(err, &notExists) {
			notFound(w)
			return
		}
		internalError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		internalError(w, err)
		return
	}
}

func mDelete(w http.ResponseWriter, id string, op *models.RepositoryOperation) {
	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "invalid uuid")
		return
	}

	err = op.Delete(models.ObjectID(uid))
	if err != nil {
		var notExists *models.NotExistsError
		if errors.As(err, &notExists) {
			notFound(w)
			return
		}
		internalError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func create(w http.ResponseWriter, r *http.Request, op *models.RepositoryOperation) {
	obj, err := models.ModelParsers[op.Name](r.Body)
	defer r.Body.Close()

	if err != nil {
		internalError(w, err)
		return
	}

	if err = op.Create(obj); err != nil {
		internalError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		internalError(w, err)
		return
	}
}

func update(w http.ResponseWriter, r *http.Request, id string, op *models.RepositoryOperation) {
	uid, err := uuid.Parse(id)
	if err != nil {
		badRequest(w, "invalid uuid")
		return
	}

	obj, err := models.ModelParsers[op.Name](r.Body)
	defer r.Body.Close()

	if err != nil {
		internalError(w, err)
		return
	}
	oid := models.ObjectID(uid)
	obj.SetID(&oid)

	if err = op.Update(obj); err != nil {
		var notExists *models.NotExistsError
		if errors.As(err, &notExists) {
			notFound(w)
			return
		}
		internalError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		internalError(w, err)
		return
	}
}
