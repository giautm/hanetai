package main

import (
	"encoding/csv"
	"encoding/json"
	"strconv"
)

type ProfileMeCmd struct{}

func (l *ProfileMeCmd) Run(ctx *CliContext) error {
	c := ctx.NewClient()
	profile, err := c.Profile.Me(ctx.Context)
	if err != nil {
		return err
	}
	if ctx.JSON {
		return json.NewEncoder(ctx.Writer()).Encode(profile)
	}

	s := csv.NewWriter(ctx.Writer())
	defer s.Flush()

	if !ctx.NoHeader {
		err = s.Write([]string{
			"ID",
			"Name",
			"Email",
		})
		if err != nil {
			return err
		}
	}
	err = s.Write([]string{
		strconv.Itoa(profile.ID),
		profile.Name,
		profile.Email,
	})
	if err != nil {
		return err
	}
	return nil
}
