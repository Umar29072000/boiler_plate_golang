package main

import (
	"boiler_plate_be_golang/app/cmd"
	"boiler_plate_be_golang/pkg/redis"

	// Load tzdata
	_ "time/tzdata"
)

func main() {
	cmd.Execute()
	defer redis.Close()
}
