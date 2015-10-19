package main

import (
	"fmt"
	"github.com/niedbalski/ultraeasy"
)

func main() {
	devices, err := ultraeasy.GetDevices()
	if err != nil {
		fmt.Println(err)
	}

	for _, device := range devices {
		if device.GetSerial() == "BWX843BAR" {
			fmt.Println("Serial:", device.GetSerial(), "Version:", device.GetVersion())

			device.GetHandler().GetAllReadingsCallback(func(reading *ultraeasy.UltraEasyReading) {
				fmt.Println(reading.GetValue())
			})
		}
	}

}
