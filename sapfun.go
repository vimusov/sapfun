/*
	sapfun - Utility that takes control over your video card coolers to keep it cool and steady.
	Works with `amdgpu' kernel module only.

	Copyright (C) 2022 Vadim Kuznetsov <vimusov@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.
	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	MaxSpeedTemp     = 80
	ForceCoolingTemp = 70
	SlowCoolingTemp  = 65
	StopCoolingTemp  = 57
)

func readValue(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Unable to read content of file '%s', error '%s'.\n", filePath, err)
	}
	return string(bytes.TrimRight(content, "\n"))
}

func writeValue(filePath string, value uint64) {
	fileDesc, err := os.OpenFile(filePath, os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		log.Fatalf("Unable to open file '%s', error '%s'.\n", filePath, err)
	}

	defer func(fileDesc *os.File) {
		err := fileDesc.Close()
		if err != nil {
			log.Fatalf("Unable to close file '%s', error '%s'.\n", fileDesc.Name(), err)
		}
	}(fileDesc)

	content := fmt.Sprintf("%d\n", value)

	_, err = fileDesc.WriteString(content)
	if err != nil {
		log.Fatalf("Unable to write value '%q' to file '%s', error '%s'.\n", content, filePath, err)
	}
}

func findRootDir() string {
	for attempt := 1; attempt < 120; attempt++ {
		time.Sleep(1 * time.Second)

		dirsPaths, err := filepath.Glob("/sys/class/hwmon/hwmon?")
		if err != nil {
			log.Printf("Unable to find module directory, error '%v'\n.", err)
			continue
		}

		for _, dirPath := range dirsPaths {
			if readValue(filepath.Join(dirPath, "name")) == "amdgpu" {
				return dirPath
			}
		}
	}

	log.Fatal("Unable to find 'amdgpu' module.\n")
	return ""
}

func getCurTemp(rootDir string) uint64 {
	filesPaths, err := filepath.Glob(filepath.Join(rootDir, "temp?_label"))
	if err != nil {
		log.Fatalf("Unable to enumetate sensors in '%s', error '%s'.\n", rootDir, err)
	}

	for _, filePath := range filesPaths {
		if readValue(filePath) != "junction" {
			continue
		}

		fileName := filepath.Base(filePath)
		newName := strings.Replace(fileName, "_label", "_input", 1)
		valuePath := filepath.Join(filepath.Dir(filePath), newName)

		rawValue := readValue(valuePath)

		tempValue, err := strconv.ParseUint(rawValue, 10, 32)
		if err != nil {
			log.Fatalf("Invalid temperature value '%s'.\n", rawValue)
		}
		return tempValue / 1000
	}

	log.Fatalf("Unable to find sensor in '%s'.\n", rootDir)
	return 0
}

func findFileByMask(rootDir string, mask string) string {
	filesPaths, err := filepath.Glob(filepath.Join(rootDir, mask))
	if err != nil {
		log.Fatalf("Unable to find file in '%s' by mask '%s'.\n", rootDir, mask)
	}
	if len(filesPaths) != 1 {
		log.Fatalf("More than one file found by mask '%s' in %s: %v.\n", mask, rootDir, filesPaths)
	}
	return filesPaths[0]
}

func setManualMode(rootDir string) {
	writeValue(findFileByMask(rootDir, "pwm?_enable"), 1)
}

func setPWMValue(rootDir string, newValue uint64) {
	filePath := findFileByMask(rootDir, "pwm?")
	rawValue := readValue(filePath)

	curValue, err := strconv.ParseUint(rawValue, 10, 32)
	if err != nil {
		log.Fatalf("Invalid PWM value '%s'.\n", rawValue)
	}

	if curValue != newValue {
		writeValue(filePath, newValue)
	}
}

func adjustFanSpeed(rootDir string, forceCooling bool) bool {
	temp := getCurTemp(rootDir)
	if temp >= MaxSpeedTemp {
		setPWMValue(rootDir, 255) // ~3600 RPM
		return true
	}
	if temp > ForceCoolingTemp && forceCooling {
		setPWMValue(rootDir, 255) // ~3600 RPM
		return true
	}
	if temp >= SlowCoolingTemp {
		setPWMValue(rootDir, 127) // ~2000 RPM
		return false
	}
	if temp >= StopCoolingTemp {
		setPWMValue(rootDir, 64) // ~800 RPM
		return false
	}
	setPWMValue(rootDir, 0)
	return false
}

func regulateTemp() {
	rootDir := findRootDir()
	setManualMode(rootDir)
	forceCooling := getCurTemp(rootDir) > ForceCoolingTemp
	for {
		forceCooling = adjustFanSpeed(rootDir, forceCooling)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	if os.Getuid() != 0 {
		log.Fatalf("This program requires root privileges to run.")
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go regulateTemp()
	<-signals
}
