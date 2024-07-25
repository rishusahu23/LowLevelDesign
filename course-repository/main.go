package main

import (
	"fmt"
)

func main() {
	courseRepo := &InMemoryCourseRepository{
		Courses: make(map[string]*Course),
	}
	studentRepo := &InMemoryStudentRepository{
		Students: make(map[int]*Student),
	}
	courseRepo.Save(&Course{Code: "CS101", Nane: "Intro to Computer Science", Instructor: "Prof. Smith", MaxCapacity: 2})
	courseRepo.Save(&Course{Code: "MATH101", Nane: "Calculus I", Instructor: "Prof. Johnson", MaxCapacity: 2})

	studentRepo.Save(&Student{ID: 1, Name: "Alice"})
	studentRepo.Save(&Student{ID: 2, Name: "Bob"})

	registrationSvc := DefaultRegistrationService{
		courseRepo:  courseRepo,
		studentRepo: studentRepo,
	}

	fmt.Println(registrationSvc.RegisterStudentForCourse(1, "CS101"))
	fmt.Println(registrationSvc.RegisterStudentForCourse(2, "CS101"))
	fmt.Println(registrationSvc.RegisterStudentForCourse(1, "CS101"))

	courses, _ := registrationSvc.GetRegisteredCourses(1)
	for _, course := range courses {
		fmt.Println("Registered courses:", course.Nane)
	}
}
