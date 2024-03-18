package main

import "vk-test-spring/internal/app"

// @title Films library API
// @version 1.0
// @description API Server for Films library

// @host localhost:8080
// @BasePath /films

// @SecurityDefinitions basicAuth
// @SecurityScheme basic
// @Security BasicAuth

const configPath = "../../configs/main"

func main() {
	app.Run(configPath)
}
