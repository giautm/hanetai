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
	Page       int    `kong:"optional,name='page',help:'The page number'"`
	Size       int    `kong:"optional,name='size',help:'Number of items per page'"`
}

func (l *PersonLsCmd) Run(ctx *CliContext) error {
	c := ctx.NewClient()
	items, err := c.Persons.ListByPlace(ctx.Context, hanetai.PersonListByPlaceRequest{
		PlaceID: l.PlaceID,
		Type:    l.PersonType,
		Page:    l.Page,
		Size:    l.Size,
	})
	if err != nil {
		return err
	}
	if ctx.JSON {
		return json.NewEncoder(ctx.Writer()).Encode(items)
	}

	s := csv.NewWriter(ctx.Writer())
	defer s.Flush()

	if !ctx.NoHeader {
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

type PersonRegisterCmd struct {
	PlaceID int      `kong:"required,name='place-id',help:'The place that person belong to'"`
	AliasID string   `kong:"required,name='alias-id',help:'The alias ID of person will register'"`
	Photo   *os.File `kong:"required,name='photo',help:'The photo of person'"`
	Name    string   `kong:"required,name='name',help:'The name of person'"`

	PersonType string `kong:"optional,name='avatar',help:'The avatar of person',default:'0'"`
	Title      string `kong:"optional,name='title',help:'The title of person',default:'Nhân viên'"`
}

func (r *PersonRegisterCmd) Run(ctx *CliContext) error {
	faceReq := &hanetai.PersonFaceUpdateRequest{
		AliasID: r.AliasID,
		PlaceID: r.PlaceID,
		File:    r.Photo,
	}
	defer r.Photo.Close()

	c := ctx.NewClient()
	person, err := c.Persons.Register(ctx.Context, hanetai.PersonRegisterRequest{
		PersonFaceUpdateRequest: faceReq,

		Name:  r.Name,
		Title: r.Title,
		Type:  r.PersonType,
	})
	if err != nil {
		return err
	}

	if ctx.JSON {
		return json.NewEncoder(ctx.Writer()).Encode(person)
	}
	fmt.Printf("Successfully register %q\n", person.ID)
	return nil
}
