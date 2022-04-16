package configs

import (
	"fmt"
	"os"
)

func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("No env. port")
	}
	return ":" + port
}
