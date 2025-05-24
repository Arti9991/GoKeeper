package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Arti9991/GoKeeper/server/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Server build version: %s\n", buildVersion)
	fmt.Printf("Server build date: %s\n", buildDate)
	fmt.Printf("Server build commit: %s\n", buildCommit)

	time.Sleep(1 * time.Second)

	err := server.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
