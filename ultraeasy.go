/*
Go bindings/Library for interact with the Third party driver for LifeScan ultraeasy devices.
Author: Jorge Niedbalski <jnr@metaklass.org>
License: GPL 3
*/

package ultraeasy

/*
#cgo CFLAGS: -Wall -Werror -I/usr/local/include
#cgo LDFLAGS: -lultraeasy

#include <ultraeasy/ultraeasy.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"path"
	"path/filepath"
	"time"
)

const (
	DEFAULT_SYS_CLASS_PATTERN = "/sys/class/tty/ttyUSB**"
)

type UltraEasyReading struct {
	time  time.Time
	value uint
}

type UltraEasyDevice struct {
	path    string
	serial  string
	version string
	handler *UltraEasy
}

type UltraEasy struct {
	handler *C.struct_ultraeasy
}

func (r *UltraEasyReading) GetValue() uint {
	return r.value
}

func (r *UltraEasyReading) GetTime() time.Time {
	return r.time
}

func (r *UltraEasyReading) GetTimeRFC3339() string {
	return r.time.Format(time.RFC3339)
}

func (u *UltraEasy) GetReadingsByTimeCallback(t time.Time, cb func(u *UltraEasyReading)) {
	u.GetAllReadingsCallback(func(u *UltraEasyReading) {
		if u.time.Equal(t) {
			cb(u)
		}
	})
}

func (u *UltraEasy) GetAllReadingsCallback(cb func(u *UltraEasyReading)) error {
	var x C.struct_ultraeasy_record

	numReadings := C.ultraeasy_num_records(u.handler)
	if numReadings == 0 {
		return fmt.Errorf("No records found on device")
	}

	for n := 0; n <= int(numReadings); n++ {
		ret := C.ultraeasy_get_record(u.handler, C.uint(n), &x)
		if ret == 0 {
			cb(&UltraEasyReading{
				time:  time.Unix(int64(x.raw.date), 0),
				value: uint(x.raw.reading),
			})
		}
	}
	return nil
}

func (u *UltraEasy) GetReadingsByTime(t time.Time) ([]UltraEasyReading, error) {
	var ret []UltraEasyReading

	records, err := u.GetAllReadings()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record.time.Equal(t) {
			ret = append(ret, record)
		}
	}

	return ret, nil

}

func (u *UltraEasy) GetVersion() (string, error) {
	version, err := C.ultraeasy_read_version(u.handler)
	if err != nil {
		return "", err
	}

	return C.GoString(version), nil
}

func (u *UltraEasy) GetSerial() (string, error) {
	serial, err := C.ultraeasy_read_serial(u.handler)
	if err != nil {
		return "", err
	}

	return C.GoString(serial), nil
}

func (u *UltraEasy) GetAllReadings() ([]UltraEasyReading, error) {
	var x C.struct_ultraeasy_record
	var records []UltraEasyReading

	numReadings := C.ultraeasy_num_records(u.handler)

	if numReadings == 0 {
		return nil, fmt.Errorf("No records found on device")
	}

	for n := 0; n <= int(numReadings); n++ {
		ret := C.ultraeasy_get_record(u.handler, C.uint(n), &x)
		if ret == 0 {
			records = append(records, UltraEasyReading{
				time:  time.Unix(int64(x.raw.date), 0),
				value: uint(x.raw.reading),
			})
		}
	}

	return records, nil
}

func (u *UltraEasy) Close() {
	C.ultraeasy_close(u.handler)
}

func getDevicePath(deviceClass string) string {
	return fmt.Sprintf("/dev/%s", path.Base(deviceClass))
}

func (d *UltraEasyDevice) GetVersion() string {
	return d.version
}

func (d *UltraEasyDevice) GetPath() string {
	return d.path
}

func (d *UltraEasyDevice) GetSerial() string {
	return d.serial
}

func (d *UltraEasyDevice) GetHandler() *UltraEasy {
	return d.handler
}

func GetDevices() ([]UltraEasyDevice, error) {
	var devices []UltraEasyDevice

	classes, err := filepath.Glob(DEFAULT_SYS_CLASS_PATTERN)
	if err != nil {
		return nil, err
	}

	for _, class := range classes {
		devicePath := getDevicePath(class)

		handler, err := NewUltraEasy(devicePath)
		if err != nil {
			continue
		}

		serial, err := handler.GetSerial()
		if err != nil {
			continue
		}

		version, err := handler.GetVersion()
		if err != nil {
			continue
		}

		devices = append(devices, UltraEasyDevice{
			path:    devicePath,
			serial:  serial,
			version: version,
			handler: handler,
		})
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("No UltraEasy Device found on the system")
	}

	return devices, nil
}

func NewUltraEasy(device string) (*UltraEasy, error) {
	var err error

	ultraeasy := &UltraEasy{}
	if ultraeasy.handler, err = C.ultraeasy_open(C.CString(device)); err != nil {
		return nil, err
	}

	return ultraeasy, err
}
