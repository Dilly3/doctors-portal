package controllers

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dilly3/doctors-portal/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var PgHandler = NewHandler()

func SetupRouter() *mux.Router {

	var dir string
	router := mux.NewRouter()

	flag.StringVar(&dir, "dir", ".", "pages/static/")
	flag.Parse()

	router.HandleFunc("/", PgHandler.Indexhandler).Methods("GET")
	// Patirnts routes
	{
		router.HandleFunc("/registerpatient", PgHandler.RegisterPatientHandler).Methods("GET")
		router.HandleFunc("/postregisterpatient", PgHandler.PostRegisterPatientHandler).Methods("POST")
		router.HandleFunc("/patientlogin", PgHandler.PatientLoginHandler).Methods("GET")
		router.HandleFunc("/postpatientlogin", PgHandler.PostLoginPatientdHandler()).Methods("POST")
		router.HandleFunc("/patientdashboard", PgHandler.PatientHomeHandler).Methods("GET")
		router.HandleFunc("/patientlogout", PgHandler.PatientLogoutHandler).Methods("GET")
		router.HandleFunc("/doctorlist", PgHandler.DoctorListHandler).Methods("GET")
		router.HandleFunc("/checkappointments", PgHandler.CheckPatientAppointmentHandler).Methods("GET")
		router.HandleFunc("/doctorappointment/{ID}", PgHandler.BookByIdHandler).Methods("GET")
		router.HandleFunc("/bookappointment/{ID}", PgHandler.PostBookByIdHandler).Methods("POST")
		router.HandleFunc("/cancelappointment/{ID}", PgHandler.PatientDeleteAppointmentHandler).Methods("GET")
	}
	//doctor routes
	{
		router.HandleFunc("/workinghours", PgHandler.DoctorWorkingHoursHandler).Methods("GET")
		router.HandleFunc("/postworkinghours", PgHandler.ChooseHoursHandler).Methods("POST")
		router.HandleFunc("/viewdoctorappointments", PgHandler.CheckDoctorAppointmentHandler).Methods("GET")
		router.HandleFunc("/canceldoc/{ID}", PgHandler.DeleteDoctorAppointmentHandler).Methods("GET")
		router.HandleFunc("/postdoctorlogin", PgHandler.PostLoginDoctordHandler).Methods("POST")
		router.HandleFunc("/doctorlogout", PgHandler.DoctorLogoutHandler).Methods("GET")
		router.HandleFunc("/doctordashboard", PgHandler.DoctorHomeHandler).Methods("GET")
		router.HandleFunc("/registerdoctor", PgHandler.RegisterDoctorHandler).Methods("GET")
		router.HandleFunc("/postregisterdoctor", PgHandler.PostRegisterDoctorHandler).Methods("POST")
		router.HandleFunc("/doctorlogin", PgHandler.DoctorLoginHandler).Methods("GET")

	}
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("pages/static/"))))

	//router.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./pages/static/")))
	PgHandler.Sessions.LoadAndSave(router)

	return router
}

func StartServer(r *mux.Router) (*http.Server, *Handler, string) {
	PgHandler1 := PgHandler
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	time.Sleep(time.Millisecond * 1600)
	port := utils.GetPortFromEnv()
	if port == "" {
		port = ":8080"
	}
	server := &http.Server{
		Addr:         "127.0.0.1" + port,
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("server started on port " + port)

	return server, PgHandler1, port
}

func GracefulShutdown(done chan error, server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nshutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()
	fmt.Print("B")
	time.Sleep(time.Millisecond * 500)
	fmt.Print("Y")
	time.Sleep(time.Millisecond * 500)
	fmt.Print("E")
	time.Sleep(time.Millisecond * 500)
	fmt.Print("!")
	time.Sleep(time.Millisecond * 500)
	fmt.Print(" ")
	time.Sleep(time.Millisecond * 500)
	fmt.Print("\n")
	done <- server.Shutdown(ctx)

}
