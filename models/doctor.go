package models

import (
	"strconv"
)

type Doctor struct {
	User
	Specialty   string `json:"specialty"`
	StartTime   int    `json:"starttime"`
	CloseTime   int    `json:"closetime"`
	StringStart string `json:"stringstart"`
	StringClose string `json:"stringclose"`
	Bookings    []Appointment
}

func (d *Doctor) SetWorkingHours() map[int]string {
	workinghrs := make(map[int]string)
	var i int

	for i = d.StartTime; i < d.CloseTime; i++ {
		if i < 12 {

			workinghrs[i] = strconv.Itoa(i) + ":" + "AM"
		}
		if i == 12 {
			workinghrs[i] = strconv.Itoa(i) + ":" + "PM"
		}
		if i > 12 {
			workinghrs[i] = strconv.Itoa(i-12) + ":" + "PM"
		}

	}
	return workinghrs
}
