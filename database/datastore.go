package database

import "github.com/dilly3/doctors-portal/models"

type DataStore interface {
	FindUserByEmailandUserName(email string, username string) (*models.Patient, error)
	CreatePatientInTable(user models.Patient)
	FindDocByEmailandUserName(email string, username string) (*models.Doctor, error)
	CreateDocInTable(user models.Doctor)
	FindDoctorByID(id string) *models.Doctor
	Authenticate(username, password string) (*models.Doctor, error)
	AuthenticatePatient(username, password string) (*models.Patient, error)
	FindDoctorByUsername(username string) (*models.Doctor, error)
	FindPatientByUsername(username string) (*models.Patient, error)
	GetAllDoctors() []models.Doctor
	CreateAppointmentInTable(user models.Appointment)
	FindPatientAppointmentsByPatientID(id string) []models.Appointment
	DeleteAppointmentbyID(id string)
	CheckAppointmentIsValidWithDoctorID(id string, time string) bool
	FindDoctorAppointmentsByDoctorID(id string) []models.Appointment
	UpdateDoctorWorkingHours(docID string, starttime string, closetime string) (*models.Doctor, error)
	UpdateDoctorWorkingHoursInt(docID string, starttime int, closetime int) (*models.Doctor, error)
}
