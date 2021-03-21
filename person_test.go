package hanetai

import (
	"context"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestPersonService_Register(t *testing.T) {
	f, err := os.Open("/workspace/avatar-xoay.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	type fields struct {
		client *Client
	}
	type args struct {
		ctx context.Context
		pu  PersonRegisterRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				pu: PersonRegisterRequest{
					PersonFaceUpdateRequest: &PersonFaceUpdateRequest{
						AliasID: "12311123",
						PlaceID: 1542,
						File:    f,
					},
					Name:  "Tui 123",
					Title: "!2312",
					Type:  "1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PersonService{
				client: tt.fields.client,
			}
			if _, err := s.Register(tt.args.ctx, tt.args.pu); (err != nil) != tt.wantErr {
				t.Errorf("PersonService.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersonService_UpdateByFaceImage(t *testing.T) {
	f, err := os.Open("/workspace/avatar-xoay.jpg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	type fields struct {
		client *Client
	}
	type args struct {
		ctx context.Context
		pu  PersonFaceUpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				pu: PersonFaceUpdateRequest{
					AliasID: "VCFL1231231",
					PlaceID: 1542,
					File:    f,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PersonService{
				client: tt.fields.client,
			}
			if err := s.UpdateByFaceImage(tt.args.ctx, tt.args.pu); (err != nil) != tt.wantErr {
				t.Errorf("PersonService.UpdateByFaceImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersonService_Remove(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx  context.Context
		data PersonRemoveRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				data: PersonRemoveRequest{
					AliasID: "VCFL1231231",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PersonService{
				client: tt.fields.client,
			}
			if err := s.Remove(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PersonService.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersonService_RemoveByPlace(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx  context.Context
		data PersonRemoveByPlaceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				data: PersonRemoveByPlaceRequest{
					AliasID: "VCFL1231231",
					PlaceID: 1542,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PersonService{
				client: tt.fields.client,
			}
			if err := s.RemoveByPlace(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PersonService.RemoveByPlace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersonService_Update(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx  context.Context
		data PersonUpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				data: PersonUpdateRequest{
					AliasID: "852576",
					PlaceID: 1542,
					Name:    "ahihi",
					Title:   "tui là ai?",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PersonService{
				client: tt.fields.client,
			}
			if err := s.Update(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PersonService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersonService_UpdateAliasID(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx  context.Context
		data PersonUpdateAliasRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				data: PersonUpdateAliasRequest{
					AliasID:  "852576",
					PersonID: "1858497629510868992",
				},
			},
			wantErr: false,
		},
	}

	//
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PersonService{
				client: tt.fields.client,
			}
			if err := s.UpdateAliasID(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("PersonService.UpdateAliasID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPersonService_ListByPlace(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx  context.Context
		data PersonListByPlaceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []PersonListItem
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				data: PersonListByPlaceRequest{
					PlaceID: 1542,
					Type:    "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PersonService{
				client: tt.fields.client,
			}
			got, err := s.ListByPlace(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonService.ListByPlace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PersonService.ListByPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAvaterSize_SetUrlValues(t *testing.T) {
	type fields struct {
		Height int
		Width  int
	}
	type args struct {
		setter url.Values
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Happy Case",
			fields: fields{
				Height: 300,
				Width:  400,
			},
			args: args{
				setter: url.Values{},
			},
			want: "height=300&width=400",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AvaterSize{
				Height: tt.fields.Height,
				Width:  tt.fields.Width,
			}
			s.SetUrlValues(tt.args.setter)
			if got := tt.args.setter.Encode(); got != tt.want {
				t.Errorf("SetUrlValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
