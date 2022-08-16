package hanetai

import (
	"context"
)

type DeviceService service

type DeviceStatus struct {
	DeviceID string
	IsOnline bool
}

type ConnectionStatusRequest struct {
	DeviceIDs []string `json:"deviceIDs" url:"deviceIDs,comma"`
}

type ConnectionStatusResponse struct {
	Devices []DeviceStatus
}

func (s *DeviceService) GetConnectionStatus(ctx context.Context, data *ConnectionStatusRequest) (*ConnectionStatusResponse, error) {
	req, err := s.client.NewRequest("device/getConnectionStatus", urlencodeBody(data))
	if err != nil {
		return nil, err
	}

	var i map[string]bool
	_, err = s.client.Do(ctx, req, &i)
	if err != nil {
		return nil, err
	}

	devices := make([]DeviceStatus, 0, len(i))
	for id, isOnline := range i {
		devices = append(devices, DeviceStatus{
			DeviceID: id,
			IsOnline: isOnline,
		})
	}

	return &ConnectionStatusResponse{
		Devices: devices,
	}, nil
}

type SetDeviceMQTTRequest struct {
	DeviceID    string `json:"deviceID" url:"deviceID"`
	Enable      bool   `json:"enable" url:"enable,int"`
	URL         string `json:"url" url:"url"`
	Username    string `json:"username" url:"username"`
	Password    string `json:"password" url:"password"`
	Base64Image bool   `json:"image" url:"image,int"`
}

func (s *DeviceService) SetDeviceMQTT(ctx context.Context, data *SetDeviceMQTTRequest) error {
	req, err := s.client.NewRequest("device/setDeviceMQTT", urlencodeBody(data))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

type DeviceInfo struct {
	DeviceID   string `json:"deviceID"`
	DeviceName string `json:"deviceName"`
	Address    string `json:"address"`

	PlaceID   int    `json:"placeID"`
	PlaceName string `json:"placeName"`
}

type ListDevicesResponse struct {
	Devices []DeviceInfo
}

func (s *DeviceService) GetListDevices(ctx context.Context) (*ListDevicesResponse, error) {
	req, err := s.client.NewRequest("device/getListDevice", urlencodeBody(nil))
	if err != nil {
		return nil, err
	}

	var i []DeviceInfo
	_, err = s.client.Do(ctx, req, &i)
	if err != nil {
		return nil, err
	}

	return &ListDevicesResponse{
		Devices: i,
	}, nil
}

type ListDevicesByPlaceRequest struct {
	PlaceID int `url:"placeID"`
}

func (s *DeviceService) GetListDevicesByPlace(ctx context.Context, data *ListDevicesByPlaceRequest) (*ListDevicesResponse, error) {
	req, err := s.client.NewRequest("device/getListDeviceByPlace", urlencodeBody(data))
	if err != nil {
		return nil, err
	}

	var i []DeviceInfo
	_, err = s.client.Do(ctx, req, &i)
	if err != nil {
		return nil, err
	}

	return &ListDevicesResponse{
		Devices: i,
	}, nil
}

type UpdateDeviceRequest struct {
	DeviceID   string `json:"deviceID" url:"deviceID"`
	DeviceName string `json:"deviceName" url:"deviceName"`
}

func (s *DeviceService) UpdateDevice(ctx context.Context, data *UpdateDeviceRequest) error {
	req, err := s.client.NewRequest("device/updateDevice", urlencodeBody(data))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}
