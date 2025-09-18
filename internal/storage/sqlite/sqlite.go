package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/codebymahmud24/student-api/internal/config"
	"github.com/codebymahmud24/student-api/internal/types"
	_ "modernc.org/sqlite"
)

type SqLite struct {
	Db *sql.DB
}

func InitialiazeDB(cfg config.Config) (*SqLite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// 🔑 Connection check
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect sqlite: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT, 
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &SqLite{Db: db}, nil
}

func (s *SqLite) CreateStudent(name string, email string, age int) (int64, error) {

	result, err := s.Db.Exec("INSERT INTO students(name, email, age) VALUES(?, ?, ?)", name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil

}

func (s *SqLite) GetStudentById(id string) (types.Student, error) {
	var student types.Student
	err := s.Db.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		return types.Student{}, err
	}
	return student, nil
}

func (s *SqLite) GetAllStudents() ([]types.Student, error) {
	var students []types.Student
	rows, err := s.Db.Query("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *SqLite) UpdateStudentById(id string, name string, email string, age int) error {
	_, err := s.Db.Exec("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?", name, email, age, id)
	if err != nil {
		return err
	}
	return nil
}
