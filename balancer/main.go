package main

import (
	"balancer/cmd"
)

func main() {
	cmd.Start()

	select {}
}
