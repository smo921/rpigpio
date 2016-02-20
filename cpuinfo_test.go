package rpigpio

import (
	"os"
	"testing"
)

func TestGetCPUInfo(t *testing.T) {
	i := RpiInfo{}
	file, err := os.Open("./test/cpuinfo.pi")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	err = i.readCPUInfo(file)
	if err != nil {
		t.Error("Received error", err)
	}
	if i.hardware != "BCM2708" {
		t.Error("Hardware family not identified")
	}
	if i.model != "Model B" ||
		i.piRevision != 2 ||
		i.ramMB != 512 ||
		i.manufacturer != "Sony" ||
		i.processor != "BCM2835" {
		t.Error("Raspberry Pi revision detected incorrectly")
	}
}

func TestBadPi(t *testing.T) {
	i := RpiInfo{}
	file, err := os.Open("./test/cpuinfo.bad_pi")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	err = i.readCPUInfo(file)
	if err == nil {
		t.Error("Bad Pi should return an error")
	}
}
