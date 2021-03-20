package hanetai

import (
	"context"
)

type PlaceService service

type Place struct {
	ID      int    `json:"id" url:"placeID"`
	Name    string `json:"name" url:"name"`
	Address string `json:"address" url:"address"`
}

func (s *PlaceService) AddPlace(ctx context.Context, place Place) (*Place, error) {
	req, err := s.client.NewRequest("place/addPlace", urlencodeBody(place))
	if err != nil {
		return nil, err
	}

	var i Place
	_, err = s.client.Do(ctx, req, &i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (s *PlaceService) UpdatePlace(ctx context.Context, place Place) error {
	req, err := s.client.NewRequest("place/updatePlace", urlencodeBody(place))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *PlaceService) Places(ctx context.Context) ([]Place, error) {
	req, err := s.client.NewRequest("place/getPlaces", nil)
	if err != nil {
		return nil, err
	}

	var i []Place
	_, err = s.client.Do(ctx, req, &i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (s *PlaceService) Remove(ctx context.Context, place Place) error {
	req, err := s.client.NewRequest("place/removePlace", urlencodeBody(place))
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}
