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
	Type    int    `json:"type"`
}

type PersonFaceUpdateRequest struct {
	AliasID string `json:"aliasID"`
	PlaceID int    `json:"placeID"`
	File    io.Reader
}

type PersonFaceURLUpdateRequest struct {
	AliasID string `json:"aliasID"`
	PlaceID int    `json:"placeID"`
	FileURL string `json:"url"`
}

type PersonRegisterRequest struct {
	*PersonFaceUpdateRequest
	Name  string `json:"name"`
	Title string `json:"title"`
	Type  string `json:"type"`
}
type PersonRegisterURLRequest struct {
	*PersonFaceURLUpdateRequest
	Name  string `json:"name"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type PersonRegisterResponse struct {
	*Person
	ID   string `json:"personID"`
	File string `json:"file"`
}

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
	req, err := s.client.NewRequest("person/register",
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
		if serr, ok := err.(*ServerError); ok && serr.Code == errCodeDuplicatedImage {
			serr.Person = p.Person
			return nil, err
		}
		return nil, err
	}

	return &p, nil
}

func (s *PersonService) RegisterByURL(ctx context.Context, pu PersonRegisterURLRequest) (*PersonRegisterResponse, error) {
	req, err := s.client.NewRequest("person/registerByUrl",
		multipartBody(nil, func(w *multipart.Writer) error {
			w.WriteField("name", pu.Name)
			w.WriteField("url", pu.FileURL)
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
		if serr, ok := err.(*ServerError); ok && serr.Code == errCodeDuplicatedImage {
			serr.Person = p.Person
			return nil, err
		}
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

	var p PersonRegisterResponse
	_, err = s.client.Do(ctx, req, &p)
	if err != nil {
		if serr, ok := err.(*ServerError); ok && serr.Code == errCodeDuplicatedImage {
			serr.Person = p.Person
		}
	}
	return err
}

func (s *PersonService) UpdateByFaceURL(ctx context.Context, pu PersonFaceURLUpdateRequest) error {
	req, err := s.client.NewRequest("person/updateByFaceUrl",
		multipartBody(nil, func(w *multipart.Writer) error {
			w.WriteField("url", pu.FileURL)
			w.WriteField("aliasID", pu.AliasID)
			w.WriteField("placeID", fmt.Sprintf("%d", pu.PlaceID))

			return nil
		}))

	if err != nil {
		return err
	}

	var p PersonRegisterResponse
	_, err = s.client.Do(ctx, req, &p)
	if err != nil {
		if serr, ok := err.(*ServerError); ok && serr.Code == errCodeDuplicatedImage {
			serr.Person = p.Person
		}
	}
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

type PersonRemoveByListAliasIDRequest struct {
	AliasIDs []string `url:"aliasIDs,comma"`
	PlaceIDs []int    `url:"placeIDs,comma"`
}

func (s *PersonService) RemoveByListAliasID(ctx context.Context, data PersonRemoveByListAliasIDRequest) error {
	req, err := s.client.NewRequest("person/removePersonByListAliasID", urlencodeBody(data))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

type PersonRemoveByIDRequest struct {
	PersonID string `url:"personID"`
}

func (s *PersonService) RemoveByID(ctx context.Context, data PersonRemoveByIDRequest) error {
	req, err := s.client.NewRequest("person/removePersonByID", urlencodeBody(data))
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
	Page    int    `url:"page"`
	Size    int    `url:"size"`
}

type PersonListItem struct {
	Name     string `json:"name"`
	AliasID  string `json:"aliasID"`
	PersonID string `json:"personID"`
	Title    string `json:"title"`
	Avatar   string `json:"avatar"`
}

func (s *PersonService) ListByPlace(ctx context.Context, data PersonListByPlaceRequest) ([]PersonListItem, error) {
	req, err := s.client.NewRequest("person/getListByPlace", urlencodeBody(data))
	if err != nil {
		return nil, err
	}

	var a []PersonListItem
	_, err = s.client.Do(ctx, req, &a)
	return a, err
}

type ListByAliasIDAllPlaceRequest struct {
	AliasID string `url:"aliasID"`
}

type PersonListItemWithPlace struct {
	PersonListItem
	PlaceID int `json:"placeID"`
}

func (s *PersonService) ListByAliasIDAllPlace(ctx context.Context, data ListByAliasIDAllPlaceRequest) ([]PersonListItemWithPlace, error) {
	req, err := s.client.NewRequest("person/getListByAliasIDAllPlace", urlencodeBody(data))
	if err != nil {
		return nil, err
	}

	var a []PersonListItemWithPlace
	_, err = s.client.Do(ctx, req, &a)
	return a, err
}

type UserInfoByAliasIDRequest struct {
	AliasID string `url:"aliasID"`
}

func (s *PersonService) UserInfoByAliasID(ctx context.Context, data UserInfoByAliasIDRequest) ([]PersonListItemWithPlace, error) {
	req, err := s.client.NewRequest("person/getUserInfoByAliasID", urlencodeBody(data))
	if err != nil {
		return nil, err
	}

	var a []PersonListItemWithPlace
	_, err = s.client.Do(ctx, req, &a)
	return a, err
}

type TakeFacePictureRequest struct {
	DeviceID string `url:"deviceID"`
}

func (s *PersonService) TakeFacePicture(ctx context.Context, data TakeFacePictureRequest) error {
	req, err := s.client.NewRequest("person/takeFacePicture", urlencodeBody(data))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}
