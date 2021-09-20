package main

import (
	"awesomeProject/pkg/transport"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)



func main() {
	var port = os.Getenv("PORT")

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.WithFields(log.Fields{"url": port}).Info("server start")
	r := transport.Router()
	fmt.Println(http.ListenAndServe(":"+port, r))

}
