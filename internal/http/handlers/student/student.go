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
	"strconv"
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

// Getting All Students
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting student by id ", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid student id", slog.String("id", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, e := storage.GetStudentById(intId)
		if e != nil {
			slog.Error("Error getting student by id", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(e))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting All student ")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

// Update one students
func UpdateStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Updating student by id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid student id", slog.Any("error", err))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		var student types.Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			slog.Error("Failed to decode request body", slog.Any("error", err))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Check which field is provided (only 1 allowed)
		var field string
		var value any

		switch {
		case student.Name != "":
			field = "name"
			value = student.Name
		case student.Email != "":
			field = "email"
			value = student.Email
		case student.Age != 0:
			field = "age"
			value = student.Age
		default:
			response.WriteJson(w, http.StatusBadRequest, fmt.Errorf("No updatable field provided"))
			return
		}

		rowsAffected, err := storage.UpdateStudentFieldById(intId, field, value)
		if err != nil {
			slog.Error("Error updating student", slog.String("id", id), slog.Any("error", err))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if rowsAffected == 0 {
			response.WriteJson(w, http.StatusNotFound, fmt.Errorf("Student not found"))
			return
		}

		updatedStudent, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, updatedStudent)
	}
}
