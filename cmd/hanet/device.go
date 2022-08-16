package main

import (
	"encoding/csv"
	"encoding/json"
	"strconv"

	"giautm.dev/hanetai"
)

type DeviceLsCmd struct {
	PlaceID int `kong:"optional,name='place-id',help:'The place to list devices'"`
}

func (r *DeviceLsCmd) Run(ctx *CliContext) (err error) {
	c := ctx.NewClient()

	var items []hanetai.DeviceInfo
	if r.PlaceID != 0 {
		data, err := c.Devices.GetListDevicesByPlace(ctx.Context, &hanetai.ListDevicesByPlaceRequest{
			PlaceID: r.PlaceID,
		})
		if err != nil {
			return err
		}
		items = data.Devices
	} else {
		data, err := c.Devices.GetListDevices(ctx.Context)
		if err != nil {
			return err
		}
		items = data.Devices
	}

	if ctx.JSON {
		return json.NewEncoder(ctx.Writer()).Encode(items)
	}

	s := csv.NewWriter(ctx.Writer())
	defer s.Flush()
	if !ctx.NoHeader {
		err = s.Write([]string{
			"PlaceID",
			"PlaceName",
			"Address",
			"DeviceID",
			"DeviceName",
		})
		if err != nil {
			return err
		}
	}
	for _, i := range items {
		err = s.Write([]string{
			strconv.FormatInt(int64(i.PlaceID), 10),
			i.PlaceName,
			i.Address,
			i.DeviceID,
			i.DeviceName,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

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
