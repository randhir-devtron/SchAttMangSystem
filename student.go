package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Class     string `json:"class"`
}

// Student_Attendance Table
type Student_Attendance struct {
	gorm.Model
	Student_Id   string    `json:"studentid"`
	PunchInTime  time.Time `json:"punchintime" gorm:"default:CURRENT_TIMESTAMP"`
	PunchOutTime time.Time `json:"punchouttime" gorm:"default:CURRENT_TIMESTAMP"`
	Day          int       `json:"day"`
	Month        int       `json:"month"`
	Year         int       `json:"year"`
	DutyTime     time.Time `json:"dutytime"`
	Student      Student   `gorm:"foreignKey:Student_Id"`
}

// Get Students
func GetStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var student []Student
	DB.Find(&student)
	json.NewEncoder(w).Encode(student)
}

// Get Student
func GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var student Student
	params := mux.Vars(r)
	DB.First(&student, params["id"])
	json.NewEncoder(w).Encode(student)
}

// Add a student into Student table
func AddStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	DB.Create(&student)
	username := student.FirstName + strconv.Itoa(int(student.ID))
	password := username
	credentials, err := AddCredentials(username, password, "Student")
	if err != nil {
		log.Printf("Error while fetching credentials %v", err)
		return
		// http.Error(w, "Error while adding credential from AddStudent table", err)
	}
	json.NewEncoder(w).Encode(student)
	json.NewEncoder(w).Encode(credentials)
}

// Delete a Student from Student Table
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var student Student
	DB.Delete(&student, params["studentid"])
	json.NewEncoder(w).Encode("The student is deleted successfully")
	json.NewEncoder(w).Encode(student)
}

// PunchIn handler for students for the very first time
func PunchInStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["studentid"]
	day := time.Now().Day()
	month := int(time.Now().Month())
	year := time.Now().Year()
	var studentAttendance Student_Attendance
	result := DB.Where("Student_Id = ? AND Day = ? AND Month = ? AND Year = ?", id, day, month, year).First(&studentAttendance)
	if result.Error == gorm.ErrRecordNotFound {
		// Create a new instance of Student_Attendance
		studentAttendance := Student_Attendance{
			Student_Id:   id,
			PunchInTime:  time.Now(),
			PunchOutTime: time.Now(),
			Day:          day,
			Month:        month,
			Year:         year,
		}
		// var studentAttendance Student_Attendance
		json.NewDecoder(r.Body).Decode(&studentAttendance)
		DB.Unscoped().Create(&studentAttendance)
		json.NewEncoder(w).Encode(studentAttendance)
		return
	}
	if result.Error != nil {
		http.Error(w, "Error while Fetching Student Attendance ", http.StatusBadRequest)
		return
	}
	if studentAttendance.PunchOutTime.After(studentAttendance.PunchInTime) {
		studentAttendance.PunchInTime = time.Now()
		result := DB.Save(&studentAttendance)
		if result != nil {
			http.Error(w, "Error while saving the student attendance", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(studentAttendance)
		return
	}
	http.Error(w, "Can not Punch In again Without Punching Out", http.StatusBadRequest)
}

// PunchOut handler for students for the very first time
func PunchOutStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	StudentId := params["studentid"]
	var studentAttendance Student_Attendance
	day := time.Now().Day()
	month := int(time.Now().Month())
	year := time.Now().Year()
	result := DB.Where("Student_Id = ? AND Day = ? AND Month = ? AND Year = ?", StudentId, day, month, year).First(&studentAttendance)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "You can not Punch out without Punching In", http.StatusBadRequest)
		return
	}
	if result.Error != nil {
		http.Error(w, "Can not fetch data from student attendance", http.StatusBadRequest)
		return
	}
	if studentAttendance.PunchOutTime.After(studentAttendance.PunchInTime) {
		http.Error(w, "You can not Punch Out without Punching in again", http.StatusBadRequest)
		return
	}
	studentAttendance.PunchOutTime = time.Now()
	DB.Save(&studentAttendance)
	json.NewEncoder(w).Encode(studentAttendance)
}

// He/She is able to see a class attendance by entering class, day, month, and year.
// (eg:- when a teacher enters class, day, month and year then that would fetch the class's attendance for that particular day month and year.)
func GetStudentAttendanceByClass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract class, day, month, and year from the request URL parameters
	params := mux.Vars(r)
	class := params["class"]
	day, err := strconv.Atoi(params["day"])
	if err != nil {
		// Handle error if day is not a valid integer
		http.Error(w, "Invalid day", http.StatusBadRequest)
		return
	}
	month, err := strconv.Atoi(params["month"])
	if err != nil {
		// Handle error if month is not a valid integer
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(params["year"])
	if err != nil {
		// Handle error if year is not a valid integer
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	// Query student attendance records with matching class, day, month, and year
	var studentAttendances []Student_Attendance
	DB.Joins("JOIN students ON student_attendances.StudentId = students.StudentId").
		Where("students.class = ? AND student_attendances.day = ? AND student_attendances.month = ? AND student_attendances.year = ?", class, day, month, year).
		Find(&studentAttendances)

	// Encode the student attendance records into JSON and send it as the response
	json.NewEncoder(w).Encode(studentAttendances)
}

// He/She is able to see his/her attendance by just entering their respective student ID, month and year.
// So that would fetch the daily attendance log for that particular month and year for that student ID.
func GetStudentAttendanceById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	month := params["month"]
	year := params["year"]
	var studentattendance []Student_Attendance
	// fmt.Println("Hey")
	result := DB.Where("StudentId = ? AND Month = ? AND Year = ?", id, month, year).Find(&studentattendance)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Encode studentattendance slice into JSON format and send it as the response
	err := json.NewEncoder(w).Encode(studentattendance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json.NewEncoder(w).Encode(studentattendance)
	// fmt.Println("Hey", studentattendance)
}
