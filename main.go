package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func initializeRouter() {
	r := mux.NewRouter()

	http.HandleFunc("/login", Login)
	http.HandleFunc("/home", Home)

	r.HandleFunc("/principals", GetPrincipals).Methods("GET")
	r.HandleFunc("/principal/{principalid}", GetPrincipal).Methods("GET")
	r.HandleFunc("/principal", AddPrincipal).Methods("POST")

	// Teacher
	r.HandleFunc("/teachers", GetTeachers).Methods("GET")
	r.HandleFunc("/teacher/{teacherid}", GetTeacher).Methods("GET")
	r.HandleFunc("/teacher", AddTeacher).Methods("POST")
	r.HandleFunc("/punchinteacher/{teacherid}", PunchInTeacher).Methods("POST")
	r.HandleFunc("/punchoutteacher/{teacherid}", PunchOutTeacher).Methods("POST")
	// Principal Can see the attendance of teachers
	r.HandleFunc("/teacher/{teacherid}/{month}/{year}", GetTeacherAttendance).Methods("GET")

	// Student
	r.HandleFunc("/students", GetStudents).Methods("GET")
	r.HandleFunc("/student/{id}", GetStudent).Methods("GET")
	r.HandleFunc("/student", AddStudent).Methods("POST")
	// Punch_In and Punch_Out for Students
	r.HandleFunc("/punchinstudent/{studentid}", PunchInStudent).Methods("POST")
	r.HandleFunc("/punchoutstudent/{studentid}", PunchOutStudent).Methods("POST")
	r.HandleFunc("/student/{class}/{day}/{month}/{year}", GetStudentAttendanceByClass).Methods("GET")
	r.HandleFunc("/student/{studentid}/{month}/{year}", GetStudentAttendanceById).Methods("GET")

	// r.HandleFunc("/users", GetUsers).Methods("GET")
	// r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	// r.HandleFunc("/users", CreateUser).Methods("POST")
	// r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	// r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe((":9000"), r))
}

func main() {
	InitialMigration()
	initializeRouter()
	DB.Exec("ALTER TABLE teacher_attendance ADD FOREIGN KEY (teacher_id) REFERENCES teacher(id);")
	DB.Exec("ALTER TABLE student_attendance ADD FOREIGN KEY (student_id) REFERENCES student(id);")
	fmt.Println("Successfully connected!")
}
