package models

type Patient struct {
	User
	Appointments []Appointment
}
