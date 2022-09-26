package hanetai

import (
	"context"
)

type ProfileService service

type Profile struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *ProfileService) Me(ctx context.Context) (*Profile, error) {
	req, err := s.client.NewRequest("profile/getProfile", urlencodeBody(nil))
	if err != nil {
		return nil, err
	}

	var a Profile
	_, err = s.client.Do(ctx, req, &a)
	if err != nil {
		return nil, err
	}
	return &a, err
}
