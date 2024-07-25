package main

type Student struct {
	ID                int
	Name              string
	RegisteredCourses map[string]bool
}

type Course struct {
	Code        string
	Nane        string
	Instructor  string
	MaxCapacity int
	Enrolled    int
}
