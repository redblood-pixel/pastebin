package main

import "github.com/redblood-pixel/pastebin/internal/app"

const configPath = "configs/dev.yaml"

func main() {
	app.Run(configPath)
}
