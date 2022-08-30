package models

type Appointment struct {
	ID              string `json:"id" gorm:"primaryKey"`
	Purpose         string `json:"purpose"`
	PatientID       string `json:"patientID" gorm:"size:50"`
	DoctorID        string `json:"doctorID" gorm:"size:50"`
	Date            string `json:"date"`
	AppointmentHour string `json:"appointmentHour"`
	DoctorName      string `json:"doctorName"`
	PatientName     string `json:"patientName"`
}

func (a Appointment) Add(val int) int {
	return val + 1
}
