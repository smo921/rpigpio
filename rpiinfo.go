package rpigpio

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// HWInfo interface
type HWInfo interface {
	GetCPUInfo() error
}

// RpiInfo holds system specs/info
type RpiInfo struct {
	piRevision   float32
	hardware     string
	manufacturer string
	model        string
	processor    string
	ramMB        int
	revision     string
}

func (info *RpiInfo) getRevision() error {
	rev := info.revision[len(info.revision)-4:]
	switch string(rev) {
	case "0002", "0003":
		info.model = "Model B"
		info.piRevision = 1.0
		info.ramMB = 256
		info.processor = "BCM2835"
	case "0004":
		info.model = "Model B"
		info.piRevision = 2.0
		info.ramMB = 256
		info.manufacturer = "Sony"
		info.processor = "BCM2835"
	case "0005":
		info.model = "Model B"
		info.piRevision = 2.0
		info.ramMB = 256
		info.manufacturer = "Qisda"
		info.processor = "BCM2835"
	case "0006":
		info.model = "Model B"
		info.piRevision = 2.0
		info.ramMB = 256
		info.manufacturer = "Egoman"
		info.processor = "BCM2835"
	case "0007":
		info.model = "Model A"
		info.piRevision = 2.0
		info.ramMB = 256
		info.manufacturer = "Egoman"
		info.processor = "BCM2835"
	case "0008":
		info.model = "Model A"
		info.piRevision = 2.0
		info.ramMB = 256
		info.manufacturer = "Sony"
		info.processor = "BCM2835"
	case "0009":
		info.model = "Model A"
		info.piRevision = 2.0
		info.ramMB = 256
		info.manufacturer = "Qisda"
		info.processor = "BCM2835"
	case "000d":
		info.model = "Model B"
		info.piRevision = 2.0
		info.ramMB = 512
		info.manufacturer = "Egoman"
		info.processor = "BCM2835"
	case "000e":
		info.model = "Model B"
		info.piRevision = 2.0
		info.ramMB = 512
		info.manufacturer = "Sony"
		info.processor = "BCM2835"
	case "000f":
		info.model = "Model B"
		info.piRevision = 2.0
		info.ramMB = 512
		info.manufacturer = "Qisda"
		info.processor = "BCM2835"
	case "0010":
		info.model = "Model B+"
		info.piRevision = 1.0
		info.ramMB = 512
		info.manufacturer = "Sony"
		info.processor = "BCM2835"
	case "0011":
		info.model = "Compute Module"
		info.piRevision = 1.0
		info.ramMB = 512
		info.manufacturer = "Sony"
		info.processor = "BCM2835"
	case "0012":
		info.model = "Model A+"
		info.piRevision = 1.0
		info.ramMB = 256
		info.manufacturer = "Sony"
		info.processor = "BCM2835"
	case "0013":
		info.model = "Model B+"
		info.piRevision = 1.2
		info.ramMB = 512
		info.processor = "BCM2835"
	default:
		info.model = "Unknown"
		info.piRevision = -1
		info.ramMB = -1
		info.processor = "Unknown"
		return errors.New("cpuinfo: unknown device")
	}
	return nil
}

func (info *RpiInfo) readCPUInfo(f *os.File) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(scanner.Text()) < 1 {
			continue
		}
		fields := strings.Split(scanner.Text(), ":")
		if len(fields) == 2 {
			key := strings.TrimSpace(fields[0])
			val := strings.TrimSpace(fields[1])
			switch key {
			case "Hardware":
				info.hardware = val
			case "Revision":
				info.revision = val
				err := info.getRevision()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GetCPUInfo reads the cpu info from the system
func (info *RpiInfo) GetCPUInfo() error {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return err
	}
	defer file.Close()
	info.readCPUInfo(file)
	return nil
}
