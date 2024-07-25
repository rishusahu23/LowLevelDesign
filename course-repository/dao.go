package main

import (
	"errors"
	"sync"
)

type CourseRepository interface {
	FindByCode(code string) (*Course, error)
	FindByNane(name string) ([]*Course, error)
	Save(course *Course) error
}

type StudentRepository interface {
	FindById(id int) (*Student, error)
	Save(student *Student) error
}

var (
	_ CourseRepository  = &InMemoryCourseRepository{}
	_ StudentRepository = &InMemoryStudentRepository{}
)

type InMemoryCourseRepository struct {
	Courses map[string]*Course
	mu      sync.Mutex
}

func (i *InMemoryCourseRepository) FindByCode(code string) (*Course, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	course, exist := i.Courses[code]
	if !exist {
		return nil, errors.New("course not found by course id")
	}
	return course, nil
}

func (i *InMemoryCourseRepository) FindByNane(name string) ([]*Course, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	var result []*Course
	for _, course := range i.Courses {
		if course.Nane == name {
			result = append(result, course)
		}
	}
	return result, nil
}

func (i *InMemoryCourseRepository) Save(course *Course) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Courses[course.Code] = course
	return nil
}

type InMemoryStudentRepository struct {
	Students map[int]*Student
	mu       sync.Mutex
}

func (i *InMemoryStudentRepository) FindById(id int) (*Student, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	student, exist := i.Students[id]
	if !exist {
		return nil, errors.New("no student found")
	}
	return student, nil
}

func (i *InMemoryStudentRepository) Save(student *Student) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Students[student.ID] = student
	return nil
}
