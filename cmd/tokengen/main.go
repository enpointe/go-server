// tokengen - generates JWT authorization tokens that
// can be used to access the protected enpoints of the server.
package main

import (
	"flag"
	"fmt"
	"os"

	server "github.com/enpointe/go-server"
)

func main() {
	username := flag.String("username", "guest",
		"The username to use in the generation of the JWT authorization token")
	time := flag.Int("time", 1200, "The amount of time in seconds before the token expires")
	admin := flag.Bool("admin", false, "Boolean flag to indicate whether the user is privilege admin user")
	configFile := flag.String("config", server.ConfigFilename, "The configuration file")
	quite := flag.Bool("quite", false, "No verbose output, only output token")
	flag.Parse()
	if !(*quite) {
		fmt.Printf("Reading %s ...\n", *configFile)
	}
	config, err := server.ReadConfig(*configFile)
	if err != nil {
		fmt.Printf("Failed to read configuration file: %s\n", err.Error())
		os.Exit(-1)
	}
	if !(*quite) {
		fmt.Println("Generating Token ...")
	}
	token, err := server.GenerateToken(*username, *admin, *time, []byte(config.JWTKey))
	if err != nil {
		fmt.Println("Failed to generate token")
		fmt.Print(err)
		os.Exit(-1)
	}
	if !(*quite) {
		fmt.Println("JWT Authentication Token:")
	}
	fmt.Println(token)
}
