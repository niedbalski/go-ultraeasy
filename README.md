Description
===========

This is a GO binding library (plus some additions) for the
[[https://github.com/daniel-thompson/ultraeasy|Libultraeasy]] library.

This library is a third party library for accessing the
Lifescan OneTouch UltraEasy devices.

Usage
======

You will need to have libultraeasy installed on the system

```bash
$ sudo add-apt-repository ppa:niedbalski/ultraeasy-driver
$ sudo apt-get update
$ sudo apt-get install libultraeasy libultreasy-dev
```

* Example code:

```go
package main

import (
	"fmt"
	"github.com/niedbalski/go-ultraeasy"
)

func main() {
	devices, err := ultraeasy.GetDevices()
	if err != nil {
		fmt.Println(err)
	}

	for _, device := range devices {
			fmt.Println("Serial:", device.GetSerial(), "Version:", device.GetVersion())

			device.GetHandler().GetAllReadingsCallback(func(reading *ultraeasy.UltraEasyReading) {
				fmt.Println(reading.GetValue(), reading.GetTime())
			})
		}
	}
}
```
