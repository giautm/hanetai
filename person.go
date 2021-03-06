package hanetai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
)

type PersonService service

type Person struct {
	Name    string `json:"name"`
	AliasID string `json:"aliasID"`
	PlaceID int    `json:"placeID"`
	Title   string `json:"title"`
	Type    string `json:"type"`
}

type PersonFaceUpdateRequest struct {
	AliasID string `json:"aliasID"`
	PlaceID int    `json:"placeID"`
	File    io.Reader
}

type PersonRegisterRequest struct {
	*PersonFaceUpdateRequest
	Name  string `json:"name"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type PersonRegisterResponse struct {
	*Person
	ID   string `json:"personID"`
	File string `json:"file"`
}

const (
	fileName = "avatar.png"
)

// AvatarSize store avatar size in height*width
type AvatarSize struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type UrlValuesSetter interface {
	Set(field string, value string)
}

func (s AvatarSize) SetUrlValues(setter UrlValuesSetter) {
	setter.Set("height", strconv.Itoa(s.Height))
	setter.Set("width", strconv.Itoa(s.Width))
}

var DefaultAvatarSize = AvatarSize{
	Height: 736,
	Width:  1280,
}

func (s *PersonService) Register(ctx context.Context, pu PersonRegisterRequest) (*PersonRegisterResponse, error) {
	req, err := s.client.NewRequest("employee/register",
		multipartBody(pu.File, func(w *multipart.Writer) error {
			w.WriteField("name", pu.Name)
			w.WriteField("aliasID", pu.AliasID)
			w.WriteField("placeID", fmt.Sprintf("%d", pu.PlaceID))
			w.WriteField("title", pu.Title)
			w.WriteField("type", pu.Type)

			return nil
		}))

	if err != nil {
		return nil, err
	}

	var p PersonRegisterResponse
	_, err = s.client.Do(ctx, req, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PersonService) UpdateByFaceImage(ctx context.Context, pu PersonFaceUpdateRequest) error {
	req, err := s.client.NewRequest("person/updateByFaceImage",
		multipartBody(pu.File, func(w *multipart.Writer) error {
			w.WriteField("aliasID", pu.AliasID)
			w.WriteField("placeID", fmt.Sprintf("%d", pu.PlaceID))

			return nil
		}))

	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

type PersonRemoveRequest struct {
	AliasID string `url:"aliasID"`
}

func (s *PersonService) Remove(ctx context.Context, data PersonRemoveRequest) error {
	req, err := s.client.NewRequest("person/remove", urlencodeBody(data))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

type PersonRemoveByPlaceRequest struct {
	AliasID string `url:"aliasID"`
	PlaceID int    `url:"placeID"`
}

func (s *PersonService) RemoveByPlace(ctx context.Context, data PersonRemoveByPlaceRequest) error {
	req, err := s.client.NewRequest("person/removeByPlace", urlencodeBody(data))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

type PersonUpdateRequest struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	AliasID string `json:"-"`
	PlaceID int    `json:"-"`
}

func (s *PersonService) Update(ctx context.Context, data PersonUpdateRequest) error {
	updates, err := json.Marshal(data)
	if err != nil {
		return err
	}

	reqData := struct {
		AliasID string `url:"aliasID"`
		PlaceID int    `url:"placeID"`
		Updates string `url:"updates"`
	}{
		AliasID: data.AliasID,
		PlaceID: data.PlaceID,
		Updates: string(updates),
	}
	req, err := s.client.NewRequest("person/update", urlencodeBody(reqData))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

type PersonUpdateAliasRequest struct {
	// NOTE(giautm): persionID is typo from Hanet
	PersonID string `url:"persionID"`
	AliasID  string `url:"aliasID"`
}

func (s *PersonService) UpdateAliasID(ctx context.Context, data PersonUpdateAliasRequest) error {
	req, err := s.client.NewRequest("person/updateAliasID", urlencodeBody(data))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

type PersonListByPlaceRequest struct {
	PlaceID int    `url:"placeID"`
	Type    string `url:"type"`
}

type PersonListItem struct {
	Name     string `json:"name"`
	AliasID  string `json:"aliasID"`
	PersonID string `json:"personID"`
	Title    string `json:"title"`
	Avatar   string `json:"avatar"`
}

func (s *PersonService) ListByPlace(ctx context.Context, data PersonListByPlaceRequest) ([]PersonListItem, error) {
	req, err := s.client.NewRequest("person/get-list-by-place", urlencodeBody(data))
	if err != nil {
		return nil, err
	}

	var a []PersonListItem
	_, err = s.client.Do(ctx, req, &a)
	return a, err
}
