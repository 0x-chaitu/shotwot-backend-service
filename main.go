package main

import "shotwot_backend/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
