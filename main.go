package main

import (
	"encoding/json"
	"fmt"
	"google-cloud-task-processor/config"
	"google-cloud-task-processor/server"
	"os"
)

func main() {
	c := initConfigFile()
	if c == nil {
		return
	}
	newServer := server.NewServer(c)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}
	fmt.Printf("Server is running on port %s \n", port)
	newServer.Start(":" + port)
}

func initConfigFile() *config.Config {
	file, _ := os.Open("config/config.json")
	decoder := json.NewDecoder(file)
	c := new(config.Config)
	err := decoder.Decode(&c)
	if err != nil {
		fmt.Printf("Cannot read config file: %s", err)
		return nil
	}
	return c
}