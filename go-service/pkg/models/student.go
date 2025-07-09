package models

import "time"

// Student represents the student data structure returned by the Node.js API
type Student struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	SystemAccess     bool      `json:"systemAccess"`
	Phone            string    `json:"phone"`
	Gender           string    `json:"gender"`
	DOB              time.Time `json:"dob"`
	Class            string    `json:"class"`
	Section          string    `json:"section"`
	Roll             int       `json:"roll"`
	FatherName       string    `json:"fatherName"`
	FatherPhone      string    `json:"fatherPhone"`
	MotherName       string    `json:"motherName"`
	MotherPhone      string    `json:"motherPhone"`
	GuardianName     string    `json:"guardianName"`
	GuardianPhone    string    `json:"guardianPhone"`
	RelationOfGuardian string  `json:"relationOfGuardian"`
	CurrentAddress   string    `json:"currentAddress"`
	PermanentAddress string    `json:"permanentAddress"`
	AdmissionDate    time.Time `json:"admissionDate"`
	ReporterName     string    `json:"reporterName"`
}

// StudentList represents a list of students from the API
type StudentList []Student 