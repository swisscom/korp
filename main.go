package main

import (
	"github.com/swisscom/korp/app"
	"github.com/swisscom/korp/korp_utils"
)

// main -
func main() {

	korp_utils.SetLogLevel("default")

	app.Create().Start()
}
