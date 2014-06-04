// +build darwin

package roomba_api

import "path/filepath"

func listAllPorts() ([]string, error) {
	return filepath.Glob("/dev/cu.*") //usbserial*")
}
