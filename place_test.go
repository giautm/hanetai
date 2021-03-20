package hanetai

import (
	"context"
	"os"
	"reflect"
	"testing"

	"golang.org/x/oauth2"
)

var ts = oauth2.StaticTokenSource(&oauth2.Token{
	AccessToken: os.Getenv("HANET_ACCESS_TOKEN"),
})

func TestPlaceService_AddPlace(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx   context.Context
		place Place
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Place
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
				place: Place{
					Name:    "My Happy Case",
					Address: "Ù ú u",
				},
			},
			want: &Place{
				Name:    "My Happy Case",
				Address: "Ù ú u",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PlaceService{
				client: tt.fields.client,
			}
			got, err := s.AddPlace(tt.args.ctx, tt.args.place)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlaceService.AddPlace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlaceService.AddPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlaceService_UpdatePlace(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx   context.Context
		place Place
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
				place: Place{
					ID:      1790,
					Name:    "My Happy Case",
					Address: "Ù ú u",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PlaceService{
				client: tt.fields.client,
			}
			if err := s.UpdatePlace(tt.args.ctx, tt.args.place); (err != nil) != tt.wantErr {
				t.Errorf("PlaceService.UpdatePlace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlaceService_Places(t *testing.T) {
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
		want    []Place
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				client: NewClient(nil, ts),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []Place{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PlaceService{
				client: tt.fields.client,
			}
			got, err := s.Places(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlaceService.Places() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlaceService.Places() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlaceService_Remove(t *testing.T) {
	type fields struct {
		client *Client
	}
	type args struct {
		ctx   context.Context
		place Place
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
				place: Place{
					ID: 1542,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PlaceService{
				client: tt.fields.client,
			}
			if err := s.Remove(tt.args.ctx, tt.args.place); (err != nil) != tt.wantErr {
				t.Errorf("PlaceService.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
