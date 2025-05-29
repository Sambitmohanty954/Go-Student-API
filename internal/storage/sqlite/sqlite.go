package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/Sambitmohanty954/students-api-golang/internal/config"
	"github.com/Sambitmohanty954/students-api-golang/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	// we can pass postgres as a driver name
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (	
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
 )`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{db}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student not found %d", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("error getting student by id: %v", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return []types.Student{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []types.Student{}, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student

		rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) UpdateStudentFieldById(id int64, field string, value any) (int64, error) {
	// Validate column names to prevent SQL injection
	validFields := map[string]bool{
		"name":  true,
		"email": true,
		"age":   true,
	}
	if !validFields[field] {
		return 0, fmt.Errorf("invalid field name: %s", field)
	}

	query := fmt.Sprintf("UPDATE students SET %s = ? WHERE id = ?", field)

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(value, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Sqlite) UpdateStudentById(id int64, student types.Student) (int64, error) {
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(student.Name, student.Email, student.Age, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s *Sqlite) DeleteStudentById(id int64) (int64, error) {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
