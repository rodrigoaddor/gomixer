package audio

import (
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

func ListDevices(deviceType DeviceType) ([]Device, error) {
	if err := ole.CoInitialize(0); err != nil {
		return nil, err
	}
	defer ole.CoUninitialize()

	immDevices, err := getIMMDevices(deviceType)
	if err != nil {
		return nil, err
	}

	devices := make([]Device, len(immDevices))
	for i, v := range immDevices {
		devices[i], err = getDevice(v)
		if err != nil {
			return nil, nil
		}
	}

	return devices, nil
}

func GetDevice(id string) (Device, error) {
	if err := ole.CoInitialize(0); err != nil {
		return Device{}, err
	}
	defer ole.CoUninitialize()

	var immDevice *wca.IMMDevice
	var device Device
	var err error

	immDevice, err = getIMMDevice(id)
	if err != nil {
		return Device{}, err
	}

	device, err = getDevice(immDevice)
	if err != nil {
		return Device{}, err
	}

	return device, nil
}

func UpdateDevice(id string, data DeviceUpdate) (Device, error) {
	if err := ole.CoInitialize(0); err != nil {
		return Device{}, err
	}
	defer ole.CoUninitialize()

	var immDevice *wca.IMMDevice
	var device Device
	var err error

	if err := updateDevice(id, data); err != nil {
		return Device{}, err
	}

	immDevice, err = getIMMDevice(id)
	if err != nil {
		return Device{}, err
	}

	device, err = getDevice(immDevice)
	if err != nil {
		return Device{}, err
	}

	return device, nil
}
