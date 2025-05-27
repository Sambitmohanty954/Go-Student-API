package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sambitmohanty954/students-api-golang/internal/storage"
	"github.com/Sambitmohanty954/students-api-golang/internal/types"
	"github.com/Sambitmohanty954/students-api-golang/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

// Create student
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating new student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Request Validation
		validate := validator.New()
		err = validate.Struct(student)
		if err != nil {

			validationErrors := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationErrors(validationErrors))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		slog.Info("User Created Successfully", slog.String("user id ", fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		//response.WriteJson(w, http.StatusCreated, student)
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}
