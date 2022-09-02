package main

import (
	"net/http"
	"os"
	"time"

	"github.com/dilly3/doctors-portal/controllers"
)

// start
func main() {
	time.Sleep(time.Millisecond * 500)
	done := make(chan error)

	mux := controllers.SetupRouter()
	server, pgHandler, port := controllers.StartServer(mux)
	go controllers.GracefulShutdown(done, server)
	http.ListenAndServe(port, pgHandler.Sessions.LoadAndSave(mux))
	<-done
	os.Exit(0)
}
