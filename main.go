package main

import (
	"net/http"

	"github.com/dilly3/doctors-portal/controllers"
	_ "github.com/dilly3/doctors-portal/utils"
)

// start
func main() {

	done := make(chan error)

	mux := controllers.SetupRouter()
	server, pgHandler, port := controllers.StartServer(mux)
	go controllers.GracefulShutdown(done, server)
	http.ListenAndServe(port, pgHandler.Sessions.LoadAndSave(mux))
	<-done
}
