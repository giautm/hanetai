package main

import (
	"encoding/csv"
	"encoding/json"
	"strconv"

	"giautm.dev/hanetai"
)

type DeviceConnectionStatusCmd struct {
	DeviceIDs []string `kong:"required,name='device-ids',help:'The ID of device to get status'"`
}

func (r *DeviceConnectionStatusCmd) Run(ctx *CliContext) error {
	c := ctx.NewClient()
	data, err := c.Devices.GetConnectionStatus(ctx.Context, &hanetai.ConnectionStatusRequest{
		DeviceIDs: r.DeviceIDs,
	})
	if err != nil {
		return err
	}
	if ctx.JSON {
		return json.NewEncoder(ctx.Writer()).Encode(data.Devices)
	}

	s := csv.NewWriter(ctx.Writer())
	defer s.Flush()

	if !ctx.NoHeader {
		err = s.Write([]string{
			"Device",
			"IsOnline",
		})
		if err != nil {
			return err
		}
	}
	for _, i := range data.Devices {
		err = s.Write([]string{
			i.DeviceID,
			strconv.FormatBool(i.IsOnline),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
