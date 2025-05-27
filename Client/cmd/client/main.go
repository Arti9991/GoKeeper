package main

import (
	"fmt"

	clientcli "github.com/Arti9991/GoKeeper/client/internal/clientCLI"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Client build version: %s\n", buildVersion)
	fmt.Printf("Client build date: %s\n", buildDate)
	fmt.Printf("Client build commit: %s\n", buildCommit)

	clientcli.StartCLI()
}
