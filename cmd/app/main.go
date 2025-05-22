//go:build app
// +build app

package main

import "github.com/redblood-pixel/pastebin/internal/app"

const configPath = "/app/configs/dev.yaml"

func main() {
	app.Run(configPath)
}
