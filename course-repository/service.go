package main

import (
	"errors"
	"sync"
)

type DefaultCourseService struct {
	courseRepo CourseRepository
}

func (s *DefaultCourseService) SearchCoursesByCode(code string) (*Course, error) {
	return s.courseRepo.FindByCode(code)
}

func (s *DefaultCourseService) SearchCoursesByName(name string) ([]*Course, error) {
	return s.courseRepo.FindByNane(name)
}

type DefaultRegistrationService struct {
	courseRepo  CourseRepository
	studentRepo StudentRepository
	mu          sync.Mutex
}

func (s *DefaultRegistrationService) RegisterStudentForCourse(studentId int, courseCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	course, err := s.courseRepo.FindByCode(courseCode)
	if err != nil {
		return err
	}
	if course.Enrolled >= course.MaxCapacity {
		return errors.New("course is full")
	}
	student, err := s.studentRepo.FindById(studentId)
	if err != nil {
		return err
	}
	if student.RegisteredCourses == nil {
		student.RegisteredCourses = make(map[string]bool)
	}
	if _, exist := student.RegisteredCourses[courseCode]; exist {
		return errors.New("user is already registered")
	}
	student.RegisteredCourses[courseCode] = true
	course.Enrolled++

	s.courseRepo.Save(course)
	s.studentRepo.Save(student)
	return nil
}

func (s *DefaultRegistrationService) GetRegisteredCourses(studentId int) ([]*Course, error) {
	student, err := s.studentRepo.FindById(studentId)
	if err != nil {
		return nil, err
	}
	var courses []*Course
	for code := range student.RegisteredCourses {
		course, err := s.courseRepo.FindByCode(code)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}
