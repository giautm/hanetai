package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"giautm.dev/hanetai"
)

type PersonRmCmd struct {
	ID string `kong:"required,name='person-id',help:'The ID of person will delete'"`
}

func (r *PersonRmCmd) Run(ctx *CliContext) error {
	c := ctx.NewClient()
	err := c.Persons.RemovePersonByID(ctx.Context, hanetai.PersonRemovePersonByIDRequest{
		PersonID: r.ID,
	})
	if err == nil {
		fmt.Printf("Successfully removed %q\n", r.ID)
	}
	return err
}

type PersonRmByAliasCmd struct {
	PlaceID int    `kong:"required,name='place-id',help:'The place that person belong to'"`
	AliasID string `kong:"required,name='alias-id',help:'The alias ID of person will delete'"`
}

func (r *PersonRmByAliasCmd) Run(ctx *CliContext) error {
	c := ctx.NewClient()
	err := c.Persons.RemoveByPlace(ctx.Context, hanetai.PersonRemoveByPlaceRequest{
		PlaceID: r.PlaceID,
		AliasID: r.AliasID,
	})
	if err == nil {
		fmt.Printf("Successfully removed %q\n", r.AliasID)
	}
	return err
}

type PersonLsCmd struct {
	PlaceID    int    `kong:"required,name='place-id',help:'The place to list persons'"`
	PersonType string `kong:"optional,name='type',help:'The person type'"`
	JSON       bool   `kong:"optional,name='json',default:false"`
}

func (l *PersonLsCmd) Run(ctx *CliContext) error {
	w := os.Stdout

	c := ctx.NewClient()
	items, err := c.Persons.ListByPlace(ctx.Context, hanetai.PersonListByPlaceRequest{
		PlaceID: l.PlaceID,
		Type:    l.PersonType,
	})
	if err != nil {
		return err
	}
	if l.JSON {
		return json.NewEncoder(w).Encode(items)
	}

	s := csv.NewWriter(w)
	defer s.Flush()

	err = s.Write([]string{
		"PersonID",
		"AliasID",
		"Title",
		"Name",
		"Avatar",
	})
	if err != nil {
		return err
	}
	for _, i := range items {
		err = s.Write([]string{
			i.PersonID,
			i.AliasID,
			i.Title,
			i.Name,
			i.Avatar,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
