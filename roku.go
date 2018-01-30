package main

import (
	"fmt"
	"os"

	"github.com/kinghrothgar/roku/roku"
)

func main() {
	ip := os.Getenv("ROKU")
	err := roku.LaunchAppNameMatch(ip, "ple")
	fmt.Printf("%v", err)
}
