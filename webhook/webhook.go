package webhook

import (
	"context"
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"hash"
	"net/http"
)

type ActionType string

const (
	ActionAdd    = ActionType("add")
	ActionDetele = ActionType("delete")
	ActionUpdate = ActionType("update")
)

type DataType string

const (
	DataDevice = DataType("device")
	DataLog    = DataType("log")
	DataPerson = DataType("person")
	DataPlace  = DataType("place")
)

type Webhook struct {
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

type PersonData struct {
	DetectedImageURL string `json:"detected_image_url"`
	PersonID         string `json:"personID"`
	AliasID          string `json:"aliasID"`
	PersonName       string `json:"personName"`
	PersonType       string `json:"personType"`
}

type DeviceData struct {
	DeviceID   string `json:"deviceID"`
	DeviceName string `json:"deviceName"`
}

type PlaceData struct {
	PlaceID   IntID  `json:"placeID"`
	PlaceName string `json:"placeName"`
}

type WebhookFn func(context.Context, *Webhook) error

func Handler(clientSecret []byte, fn WebhookFn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			// Hanet Check Request
			w.WriteHeader(http.StatusOK)
			return
		}

		var data Webhook

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !verifyHash(md5.New(), clientSecret, data.EventData) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if err := fn(r.Context(), &data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func verifyHash(h hash.Hash, secret []byte, e *EventData) bool {
	h.Write(secret)
	h.Write([]byte(e.ID))

	src := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return subtle.ConstantTimeCompare([]byte(e.Hash), dst) == 1
}
