package main

import (
	"gomixer/pkg/audio"
	"log"
)

func main() {
	var err error
	var devices []audio.Device
	if err, devices = audio.ListDevices(); err != nil {
		log.Fatal(err)
	}

	log.Print(devices)
}
