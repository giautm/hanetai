package hanetai

import (
	"context"
	"reflect"
	"testing"
)

func TestDeviceService_GetConnectionStatus(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx  context.Context
		data *ConnectionStatusRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ConnectionStatusResponse
		wantErr bool
	}{
		{
			name: "Happy Case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				data: &ConnectionStatusRequest{
					DeviceIDs: []string{"C21024B155"},
				},
			},
			want: &ConnectionStatusResponse{
				Devices: []DeviceStatus{
					{
						DeviceID: "C21024B155",
						IsOnline: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DeviceService{
				client: tt.fields.client,
			}
			got, err := s.GetConnectionStatus(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeviceService.GetConnectionStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeviceService.GetConnectionStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceService_GetListDevices(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListDevicesResponse
		wantErr bool
	}{
		{
			name: "Happy Case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DeviceService{
				client: tt.fields.client,
			}
			got, err := s.GetListDevices(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeviceService.GetListDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeviceService.GetListDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceService_GetListDevicesByPlace(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx  context.Context
		data *ListDevicesByPlaceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListDevicesResponse
		wantErr bool
	}{
		{
			name: "Happy Case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx:  context.Background(),
				data: &ListDevicesByPlaceRequest{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DeviceService{
				client: tt.fields.client,
			}
			got, err := s.GetListDevicesByPlace(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeviceService.GetListDevicesByPlace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeviceService.GetListDevicesByPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceService_UpdateDevice(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx  context.Context
		data *UpdateDeviceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx:  context.Background(),
				data: &UpdateDeviceRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DeviceService{
				client: tt.fields.client,
			}
			if err := s.UpdateDevice(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DeviceService.UpdateDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
