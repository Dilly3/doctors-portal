package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ErrorResponse struct {
	Code    int         `json:"code"`
	Error   string      `json:"error"`
	Message interface{} `json:"message"`
}
type NoAppointment struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}
type SuccessResponse struct {
	Code    int         `json:"code"`
	Object  interface{} `json:"object"`
	Message interface{} `json:"message"`
}

func GetPortFromEnv() string {
	port := os.Getenv("ROUTER_PORT")
	if port == "" || len(port) < 1 {
		port = "8080"
	}
	return port
}
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func GenerateRandomID() string {
	b := make([]byte, 9)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s%x", "HTH-", b)
}
func GenerateRandomAppointmentID() string {
	b := make([]byte, 5)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s%x", "HTH-", b)
}
