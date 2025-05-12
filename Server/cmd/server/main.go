package main

import "fmt"

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Server build version: %s\n", buildVersion)
	fmt.Printf("Server build date: %s\n", buildDate)
	fmt.Printf("Server build commit: %s\n", buildCommit)
}
