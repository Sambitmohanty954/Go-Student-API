package storage

import "github.com/Sambitmohanty954/students-api-golang/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	// Fpr Full body update
	UpdateStudentById(id int64, student types.Student) (int64, error)
	// NEW: For updating a single field
	UpdateStudentFieldById(id int64, field string, value any) (int64, error)
}
