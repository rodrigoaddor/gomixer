package audio

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

type DeviceType string

const (
	DeviceTypeOutput DeviceType = "output"
	DeviceTypeInput  DeviceType = "input"
	DeviceTypeAll    DeviceType = "all"
)

type Device struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Volume      float32    `json:"volume"`
	Type        DeviceType `json:"type"`
	Mute        bool       `json:"mute"`
}

type DeviceUpdate struct {
	Volume *float32 `json:"volume" binding:"min=0,max=1"`
	Mute   *bool    `json:"mute"`
}

func getIMMDevices(deviceType DeviceType) ([]*wca.IMMDevice, error) {
	var pEnumerator *wca.IMMDeviceEnumerator
	var pCollection *wca.IMMDeviceCollection

	var dataFlow uint32
	var err error
	dataFlow, err = deviceTypeToDataFlow(deviceType)
	if err != nil {
		return nil, err
	}

	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, ole.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &pEnumerator); err != nil {
		return nil, err
	}
	if err := pEnumerator.EnumAudioEndpoints(dataFlow, wca.DEVICE_STATE_ACTIVE, &pCollection); err != nil {
		return nil, err
	}
	defer pEnumerator.Release()
	defer pCollection.Release()

	var count uint32
	if err := pCollection.GetCount(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	devices := make([]*wca.IMMDevice, count)

	for i := uint32(0); i < count; i++ {
		if err := pCollection.Item(i, &devices[i]); err != nil {
			return nil, err
		}
	}

	return devices, nil
}

func getIMMDevice(id string) (*wca.IMMDevice, error) {
	devices, err := getIMMDevices(DeviceTypeAll)
	if err != nil {
		return nil, err
	}

	var deviceId string
	for _, device := range devices {
		if err := device.GetId(&deviceId); err != nil {
			return nil, err
		}
		if deviceId == id {
			return device, nil
		}
	}

	return nil, fmt.Errorf("device with id \"%s\" not found", id)
}

func getDevice(immDevice *wca.IMMDevice) (Device, error) {
	device := Device{}
	var endpoint *wca.IMMEndpoint
	var props *wca.IPropertyStore
	var volume *wca.IAudioEndpointVolume

	if err := immDevice.GetId(&device.Id); err != nil {
		return device, err
	}
	if err := immDevice.OpenPropertyStore(wca.STGM_READ, &props); err != nil {
		return device, err
	}
	defer props.Release()

	dispatch, err := immDevice.QueryInterface(wca.IID_IMMEndpoint)
	if err != nil {
		return device, err
	}

	endpoint = (*wca.IMMEndpoint)(dispatch)
	defer endpoint.Release()

	var dataFlow uint32
	if err := endpoint.GetDataFlow(&dataFlow); err != nil {
		return device, err
	}

	device.Type, err = dataFlowToDeviceType(dataFlow)
	if err != nil {
		return device, err
	}

	var prop wca.PROPVARIANT
	if err := props.GetValue(&wca.PKEY_Device_FriendlyName, &prop); err != nil {
		return device, err
	}
	device.Name = prop.String()
	if err := props.GetValue(&wca.PKEY_Device_DeviceDesc, &prop); err != nil {
		return device, err
	}
	device.Description = prop.String()

	if err := immDevice.Activate(wca.IID_IAudioEndpointVolume, ole.CLSCTX_ALL, nil, &volume); err != nil {
		return device, err
	}
	defer volume.Release()

	if err := volume.GetMasterVolumeLevelScalar(&device.Volume); err != nil {
		return device, err
	}
	if err := volume.GetMute(&device.Mute); err != nil {
		return device, err
	}

	return device, nil
}

func updateDevice(id string, data DeviceUpdate) error {
	var device *wca.IMMDevice
	var volume *wca.IAudioEndpointVolume
	var err error
	if device, err = getIMMDevice(id); err != nil {
		return err
	}

	if err := device.Activate(wca.IID_IAudioEndpointVolume, ole.CLSCTX_ALL, nil, &volume); err != nil {
		return err
	}

	if data.Volume != nil {
		if err := volume.SetMasterVolumeLevelScalar(*data.Volume, nil); err != nil {
			return err
		}
	}

	if data.Mute != nil {
		if err := volume.SetMute(*data.Mute, nil); err != nil {
			return err
		}
	}

	return nil
}

func dataFlowToDeviceType(dataFlow uint32) (DeviceType, error) {
	switch dataFlow {
	case 0:
		return DeviceTypeOutput, nil
	case 1:
		return DeviceTypeInput, nil
	case 2:
		return DeviceTypeAll, nil
	}
	return "", fmt.Errorf("unknown data flow (%d)", dataFlow)
}

func deviceTypeToDataFlow(deviceType DeviceType) (uint32, error) {
	switch deviceType {
	case "output":
		return 0, nil
	case "input":
		return 1, nil
	case "all":
		return 2, nil
	}
	return 0, fmt.Errorf("unknown device type (%s)", deviceType)
}

func ValidateDeviceType(deviceType DeviceType) error {
	if deviceType != DeviceTypeOutput && deviceType != DeviceTypeInput && deviceType != DeviceTypeAll {
		return fmt.Errorf("invalid device type \"%s\"", deviceType)
	}
	return nil
}
