package main

import (
	"fmt"

	clientcli "github.com/Arti9991/GoKeeper/client/internal/clientCLI"
	"github.com/Arti9991/GoKeeper/client/internal/requseter"
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

	//clientcli.StartCLI(req)

	req := requseter.NewRequester(":8082")
	// err := req.TestLogin()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	clientcli.StartCLI(req)
}
