package pdf

import (
	"testing"
	"time"

	"go-service/pkg/models"
)

// TestNewGenerator tests the creation of a new PDF generator
func TestNewGenerator(t *testing.T) {
	generator := NewGenerator()
	
	if generator == nil {
		t.Error("Expected generator to be created, got nil")
	}
	
	if generator.pdf == nil {
		t.Error("Expected PDF instance to be created, got nil")
	}
}

// TestGenerateStudentReport tests PDF generation with sample data
func TestGenerateStudentReport(t *testing.T) {
	generator := NewGenerator()
	
	// Create sample student data
	student := &models.Student{
		ID:               1,
		Name:             "Test Student",
		Email:            "test@example.com",
		Phone:            "555-1234",
		Gender:           "Male",
		DOB:              time.Date(2005, 1, 15, 0, 0, 0, 0, time.UTC),
		Class:            "Grade 10",
		Section:          "A",
		Roll:             1,
		FatherName:       "Test Father",
		FatherPhone:      "555-5678",
		MotherName:       "Test Mother",
		MotherPhone:      "555-9012",
		GuardianName:     "Test Guardian",
		GuardianPhone:    "555-3456",
		RelationOfGuardian: "Father",
		CurrentAddress:   "123 Test Street",
		PermanentAddress: "123 Test Street",
		AdmissionDate:    time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
		ReporterName:     "Test Teacher",
		SystemAccess:     true,
	}
	
	pdfBytes, err := generator.GenerateStudentReport(student)
	if err != nil {
		t.Errorf("Expected PDF generation to succeed, got error: %v", err)
	}
	
	if len(pdfBytes) == 0 {
		t.Error("Expected PDF bytes to be generated, got empty slice")
	}
	
	// Basic validation - PDF files start with "%PDF"
	if len(pdfBytes) < 4 || string(pdfBytes[:4]) != "%PDF" {
		t.Error("Generated content does not appear to be a valid PDF")
	}
} 