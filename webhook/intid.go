package webhook

import (
	"encoding/json"
	"strconv"
)

type IntID int

var (
	_ json.Unmarshaler = (*IntID)(nil)
)

func (id IntID) Int() int {
	return (int)(id)
}

func (id *IntID) UnmarshalJSON(data []byte) (err error) {
	var tmp int
	if err = json.Unmarshal(data, &tmp); err != nil {
		tmp, err = strconv.Atoi((string)(data[1 : len(data)-1]))
	}
	if err == nil {
		*id = IntID(tmp)
	}

	return err
}
