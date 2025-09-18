package storage

import "github.com/codebymahmud24/student-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id string) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	UpdateStudentById(id string, name string, email string, age int) error
}
