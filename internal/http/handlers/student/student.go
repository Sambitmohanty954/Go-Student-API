package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sambitmohanty954/students-api-golang/internal/types"
	"github.com/Sambitmohanty954/students-api-golang/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

// Create student
func New() http.HandlerFunc {
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

		response.WriteJson(w, http.StatusCreated, student)
	}
}
