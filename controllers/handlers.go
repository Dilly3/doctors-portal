package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dilly3/doctors-portal/database"
	"github.com/dilly3/doctors-portal/models"
	"github.com/dilly3/doctors-portal/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB       database.DataStore
	Sessions *scs.SessionManager
	Logger   *zap.Logger
}

func NewHandler() *Handler {

	newHandler := &Handler{
		DB:       database.NewDB(),
		Sessions: scs.New(),
		Logger:   zap.NewExample(),
	}
	newHandler.Sessions.Lifetime = 24 * time.Hour

	return newHandler
}

// Indexhandler gets the homepage
func (h Handler) Indexhandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/index.htm")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

// RegisterPatientHandler gets Patient's SignUp page
func (h Handler) RegisterPatientHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/registerpatient.htm")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

// PostRegisterPatientHandler successfully register's patient's name in the db if valid
func (h Handler) PostRegisterPatientHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Patient
	r.ParseForm()
	name := models.Capitalise(r.FormValue("name"))
	ageString := r.FormValue("age")
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	age, _ := strconv.Atoi(ageString)
	user.ID = utils.GenerateRandomID()
	user.Name = name
	user.Age = uint(age)
	user.Email = email
	user.Username = username
	user.Password = password
	_, err := h.DB.FindUserByEmailandUserName(user.Email, user.Username)
	if err == nil {
		// this user already exists
		// return a message to the user
		t, e := template.ParseFiles("pages/registerpatient.htm")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "User already exists. Check Email or Username")
		if e != nil {
			fmt.Println(e)
			return
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.PasswordHash = string(hashedPassword)
	h.DB.CreatePatientInTable(user)
	file, err3 := template.ParseFiles("pages/registerpatient.htm")
	if err3 != nil {
		fmt.Println(err3)
	}
	file.Execute(w, name+" "+" is now a Registered Patient. \n You can Login")

}

func (h Handler) PatientLoginHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/patientlogin.htm")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

// ------------------------------PostPatientLoginHandler logs in doctor if valid-----------------------------------------------------
func (h Handler) PostLoginPatientdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.Doctor
		user.Username = strings.TrimSpace(r.FormValue("username"))
		user.Password = strings.TrimSpace(r.FormValue("password"))
		_, err := h.DB.AuthenticatePatient(user.Username, user.Password)
		if err != nil {
			t, e := template.ParseFiles("pages/patientlogin.htm")
			if e != nil {
				fmt.Println(e)
				return
			}
			e = t.Execute(w, "Check username or Password")
			if e != nil {
				fmt.Println(e)
				return
			}
			return
		}
		h.Sessions.Put(r.Context(), "username", user.Username)
		http.Redirect(w, r, "patientdashboard", http.StatusFound)
	}
}

// ------------------------------PatientDashboardHandler gets Patient's Dashboard page-----------------------------------------------
func (h Handler) PatientHomeHandler(w http.ResponseWriter, r *http.Request) {
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, e := template.ParseFiles("pages/patientdashboard.htm")
	if e != nil {
		fmt.Println(e)
		return
	}

	patient, err := h.DB.FindPatientByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}

	e = t.Execute(w, patient)
	if e != nil {
		fmt.Println(e)
		return
	}
}

// ------------------------------PatientLogoutHandler logsout ---------------------------------------------------------------------
func (h Handler) PatientLogoutHandler(w http.ResponseWriter, r *http.Request) {
	h.Sessions.Remove(r.Context(), "username")
	http.Redirect(w, r, "/", http.StatusFound)
}

// -------------------------RegisterDoctorHandler gets Doctor's SignUp page-----------------------------------------------
func (h Handler) RegisterDoctorHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorregister.htm")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

// -------------------PostRegisterDoctorHandler successfully registers doctor's name in the db if valid----------------------------
func (h Handler) PostRegisterDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Doctor
	ageString := r.FormValue("age")
	age, _ := strconv.Atoi(ageString)
	user.ID = utils.GenerateRandomID()
	user.Name = models.Capitalise(strings.TrimSpace(r.FormValue("name")))
	user.Age = uint(age)
	user.Email = strings.TrimSpace(r.FormValue("email"))
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	user.Specialty = models.Capitalise(strings.TrimSpace(r.FormValue("specialty")))
	_, err := h.DB.FindDocByEmailandUserName(user.Email, user.Username)
	if err == nil {
		t, e := template.ParseFiles("pages/doctorRegister.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "User already exists, confirm email or username")
		if e != nil {
			fmt.Println(e)
			return
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.PasswordHash = string(hashedPassword)
	h.DB.CreateDocInTable(user)
	temp, errt := template.ParseFiles("pages/doctorRegister.htm")
	if errt != nil {
		fmt.Println(errt)
	}
	temp.Execute(w, user.Name+" is now a registered Doctor. Login")

}

// ------------------------------DoctorLoginHandler gets Doctor's Login page---------------------------------------------------------
func (h Handler) DoctorLoginHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorLogin.htm")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func (h Handler) PostLoginDoctordHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Doctor
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	doc, err := h.DB.Authenticate(user.Username, user.Password)
	if err != nil {
		t, e := template.ParseFiles("pages/doctorlogin.htm")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "Check Username or Password")
		if e != nil {
			fmt.Println(e)
			return
		}
		return
	}
	h.Sessions.Put(r.Context(), "username", doc.Username)
	http.Redirect(w, r, "/doctordashboard", http.StatusFound)
}

// ------------------------------DoctorDashboardHandler gets Doctor's Dashboard page-----------------------------------------------
func (h Handler) DoctorHomeHandler(w http.ResponseWriter, r *http.Request) {
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, e := template.ParseFiles("pages/doctordashboard.htm")
	if e != nil {
		fmt.Println(e)
		return
	}

	doc, err := h.DB.FindDoctorByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	e = t.Execute(w, doc)
	if e != nil {
		fmt.Println(e)
		return
	}
}

// ------------------------------DoctorLogoutHandler logsout ---------------------------------------------------------------------
func (h Handler) DoctorLogoutHandler(w http.ResponseWriter, r *http.Request) {
	h.Sessions.Remove(r.Context(), "username")
	http.Redirect(w, r, "/", http.StatusFound)
}

// ------------------------------List of Doctors for booking Appointments ---------------------------------------------------------------------
func (h Handler) DoctorListHandler(w http.ResponseWriter, r *http.Request) {
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, err := template.ParseFiles("pages/doctorlist.htm")
	if err != nil {
		log.Println(err)
		return
	}

	err = t.Execute(w, h.DB.GetAllDoctors())
	if err != nil {
		log.Fatal(err)
	}
}

func (h Handler) DoctorWorkingHoursHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/workinghours.htm")
	if e != nil {
		fmt.Println("now", e)
		return
	}
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	doc, err := h.DB.FindDoctorByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	e = t.Execute(w, doc)
	if e != nil {
		fmt.Println("no way", e)
		return
	}
}

func (h Handler) ChooseHoursHandler(w http.ResponseWriter, r *http.Request) {

	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	doc, err := h.DB.FindDoctorByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}

	e := r.ParseForm()
	if e != nil {
		fmt.Println(e)
	}

	starttime := r.PostForm.Get("start")
	closetime := r.PostForm.Get("end")
	startInt, _ := strconv.Atoi(starttime)
	closeInt, _ := strconv.Atoi(closetime)

	if startInt > closeInt {
		file, err := template.ParseFiles("pages/appointmentTimeError.htm")
		if err != nil {
			fmt.Println(err)
		}
		file.Execute(w, "START TIME LATER THAN END TIME ")
		return
	}
	checkStart := startInt > 12
	noonStart := startInt == 12
	fmt.Println(closetime)
	if checkStart {
		starttime = strconv.Itoa(startInt-12) + ":" + "00" + "PM"
	} else if noonStart {
		starttime = strconv.Itoa(startInt) + ":" + "00" + "PM"
	} else if !checkStart {
		starttime = strconv.Itoa(startInt) + ":" + "00" + "AM"
	}
	checkEnd := closeInt > 12
	noonEnd := closeInt == 12
	if checkEnd {
		closetime = strconv.Itoa(closeInt-12) + ":" + "00" + "PM"
	} else if noonEnd {
		closetime = strconv.Itoa(closeInt) + ":" + "00" + "PM"
	} else {
		closetime = strconv.Itoa(closeInt) + ":" + "00" + "AM"
	}
	h.DB.UpdateDoctorWorkingHours(doc.ID, starttime, closetime)
	h.DB.UpdateDoctorWorkingHoursInt(doc.ID, startInt, closeInt)
	http.Redirect(w, r, "doctordashboard", http.StatusFound)
}

func (h Handler) BookByIdHandler(w http.ResponseWriter, r *http.Request) {

	//This points to the html location
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, err := template.ParseFiles("pages/appointments.htm")
	if err != nil {
		fmt.Println("now", err)
		return
	}

	params := mux.Vars(r)
	ID := params["ID"]

	doctor := h.DB.FindDoctorByID(ID)
	workinghrs := doctor.SetWorkingHours()
	dean := models.Dean{
		ID:           doctor.ID,
		WorkingHours: workinghrs,
	}

	err = t.Execute(w, dean)
	if err != nil {
		fmt.Println("now", err)
		return
	}
}

func (h Handler) PostBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment

	userName := h.Sessions.GetString(r.Context(), "username")
	patient, err := h.DB.FindPatientByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	e := r.ParseForm()
	if e != nil {
		fmt.Println(e)
	}

	appointment.ID = utils.GenerateRandomAppointmentID()
	appointment.AppointmentHour = r.PostForm.Get("time")
	appointment.Purpose = r.PostForm.Get("purpose")
	params := mux.Vars(r)
	appointment.DoctorID = params["ID"]
	fmt.Println(appointment.DoctorID)
	f := h.DB.FindDoctorByID(appointment.DoctorID)
	fmt.Println(appointment.AppointmentHour)
	valid := h.DB.CheckAppointmentIsValidWithDoctorID(appointment.DoctorID, appointment.AppointmentHour)
	appointment.DoctorName = f.Name
	appointment.Date = time.Now().String()
	appointment.PatientName = patient.Name
	appointment.PatientID = patient.ID
	if valid {
		h.DB.CreateAppointmentInTable(appointment)
		http.Redirect(w, r, "/patientdashboard", http.StatusFound)
	} else {
		file, err := template.ParseFiles("pages/appointmenttimeerr.htm")
		if err != nil {
			fmt.Println(err)
		}
		file.Execute(w, utils.NoAppointment{
			Message: "Appointment time is not not Available",
			Name:    patient.Name,
		})

	}

}

func (h Handler) CheckPatientAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, e := template.ParseFiles("pages/patientcheckappointments.htm")
	if e != nil {
		fmt.Println(e)
		return
	}

	patient, err := h.DB.FindPatientByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	appointment := h.DB.FindPatientAppointmentsByPatientID(patient.ID)
	if len(appointment) < 1 {
		file, err := template.ParseFiles("pages/patientnoappointments.htm")
		if err != nil {
			fmt.Println(err)
		}

		e = file.Execute(w, utils.NoAppointment{
			Message: "You Have No Appointments",
			Name:    patient.Name})

		if e != nil {
			fmt.Println(e)
			return
		}
	}

	e = t.Execute(w, appointment)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func (h Handler) PatientDeleteAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	params := mux.Vars(r)
	ID := params["ID"]
	h.DB.DeleteAppointmentbyID(ID)
	http.Redirect(w, r, "/checkappointments", http.StatusFound)
}

func (h Handler) CheckDoctorAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, e := template.ParseFiles("pages/viewdocappointments.htm")
	if e != nil {
		fmt.Println(e)
		return
	}

	doctor, err := h.DB.FindDoctorByUsername(userName)
	if err != nil {
		log.Println(err)
		return
	}
	appointment := h.DB.FindDoctorAppointmentsByDoctorID(doctor.ID)
	if len(appointment) < 1 {
		file, err := template.ParseFiles("pages/noappointments.htm")
		if err != nil {
			fmt.Println(err)
		}

		e = file.Execute(w, utils.NoAppointment{
			Message: "You Have No Appointments",
			Name:    doctor.Name})

		if e != nil {
			fmt.Println(e)
			return
		}
	}

	e = t.Execute(w, appointment)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func (h Handler) DeleteDoctorAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	userName := h.Sessions.GetString(r.Context(), "username")
	if len(userName) < 1 {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	params := mux.Vars(r)
	ID := params["ID"]
	h.DB.DeleteAppointmentbyID(ID)
	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/viewdocappointments", http.StatusFound)
}
