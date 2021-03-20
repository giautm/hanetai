package webhook

import (
	"crypto/md5"
	"hash"
	"testing"
)

func Test_verifyHash(t *testing.T) {
	type args struct {
		h      hash.Hash
		secret []byte
		e      *EventData
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Happy Case, match hash",
			args: args{
				h:      md5.New(),
				secret: []byte(`946b9654dcfc55342c55e533805cdba6`),
				e: &EventData{
					ID:   "c75570bb-dc1a-4192-946c-ed09a34f7d77",
					Hash: "a173b27d031519da1e0cc5468eb7b9f3",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := verifyHash(tt.args.h, tt.args.secret, tt.args.e); got != tt.want {
				t.Errorf("verifyHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntID_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	var a IntID
	tests := []struct {
		name    string
		id      *IntID
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Int case",
			id:   &a,
			args: args{
				data: []byte(`1234`),
			},
			want: 1234,
		},
		{
			name: "String case",
			id:   &a,
			args: args{
				data: []byte(`"1234"`),
			},
			want: 1234,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.id.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("IntID.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (int)(*tt.id) != tt.want {
				t.Errorf("UnmarshalJSON() = %v, want %v", *tt.id, tt.want)
			}
		})
	}
}
