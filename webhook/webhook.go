package webhook

import (
	"context"
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"
	"os"
)

type ActionType string

const (
	ActionAdd    = ActionType("add")
	ActionDetele = ActionType("delete")
	ActionUpdate = ActionType("update")
)

type DataType string

const (
	DataCheckinPicture = DataType("checkin_picture")
	DataDevice         = DataType("device")
	DataLog            = DataType("log")
	DataPerson         = DataType("person")
	DataPlace          = DataType("place")
)

type Data struct {
	DataType DataType `json:"data_type"`

	*EventData
	*DeviceData
	*PersonData
	*PlaceData
}

type EventData struct {
	ActionType ActionType `json:"action_type"`
	Date       string     `json:"date"`
	KeyCode    string     `json:"keycode"`

	Hash string `json:"hash"`
	ID   string `json:"id"`
	Time uint64 `json:"time"`
}

type PersonType string

const (
	PersonEmployee PersonType = "0"
	PersonCustomer PersonType = "1"
	PersonStranger PersonType = "2"
)

func (p PersonType) String() string {
	switch p {
	case PersonEmployee:
		return "Employee"
	case PersonCustomer:
		return "Customer"
	case PersonStranger:
		return "Stranger"
	}

	return fmt.Sprintf("Unknown(%s)", (string)(p))
}

type PersonData struct {
	DetectedImageURL string `json:"detected_image_url"`
	PersonID         string `json:"personID"`
	AliasID          string `json:"aliasID"`
	PersonName       string `json:"personName"`

	PersonType PersonType `json:"personType"`
}

type DeviceData struct {
	DeviceID   string `json:"deviceID"`
	DeviceName string `json:"deviceName"`
}

type PlaceData struct {
	PlaceID   IntID  `json:"placeID"`
	PlaceName string `json:"placeName"`
}

type Handler interface {
	ServeWebhook(context.Context, *Data) error
}

type HandlerFunc func(context.Context, *Data) error

func (fn HandlerFunc) ServeWebhook(ctx context.Context, data *Data) error {
	return fn(ctx, data)
}

type Options struct {
	OnError func(context.Context, error)
	Stats   bool
	Verify  func(*EventData) bool
}

type Option = func(*Options)

func WithSecretVerify(secret []byte) Option {
	if secret == nil {
		if s := os.Getenv("HANET_CLIENT_SECRET"); s == "" {
			panic("HANET_CLIENT_SECRET is not set")
		} else {
			secret = ([]byte)(s)
		}
	}

	return func(o *Options) {
		o.Verify = func(data *EventData) bool {
			return verifyHash(md5.New(), secret, data)
		}
	}
}

func WithStats() Option {
	return func(o *Options) {
		o.Stats = true
	}
}

func WithOnError(fn func(context.Context, error)) Option {
	return func(o *Options) {
		o.OnError = fn
	}
}

func NewHTTPHandler(fn Handler, opts ...Option) http.HandlerFunc {
	o := &Options{
		Stats: false,
		Verify: func(*EventData) bool {
			return true
		},
	}

	for _, opt := range opts {
		opt(o)
	}
	if o.Stats {
		fn = ReportStats(fn)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			// Hanet Check Request
			w.WriteHeader(http.StatusOK)
			return
		}

		ctx := r.Context()

		var data Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			if o.OnError != nil {
				o.OnError(ctx, err)
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !o.Verify(data.EventData) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if err := fn.ServeWebhook(ctx, &data); err != nil {
			if o.OnError != nil {
				o.OnError(ctx, err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func verifyHash(h hash.Hash, secret []byte, e *EventData) bool {
	h.Reset()
	h.Write(secret)
	h.Write([]byte(e.ID))
	s := h.Sum(nil)

	dst := make([]byte, hex.EncodedLen(len(s)))
	hex.Encode(dst, s)

	return subtle.ConstantTimeCompare([]byte(e.Hash), dst) == 1
}
