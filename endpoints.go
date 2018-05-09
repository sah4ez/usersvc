package usersvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/satori/go.uuid"
)

type Endpoints struct {
	PostUserEndpoint  endpoint.Endpoint
	GetUserEndpoint   endpoint.Endpoint
	PatchUserEndpoint endpoint.Endpoint
	GetUsersEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostUserEndpoint:  MakePostUserEndpoint(s),
		GetUserEndpoint:   MakeGetUserEndpoint(s),
		PatchUserEndpoint: MakePatchUserEndpoint(s),
		GetUsersEndpoint:  MakeGetUsersEndpoint(s),
	}
}

func MakePostUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUserRequest)
		user := User{
			ID:       uuid.Must(uuid.NewV4(), nil).String(),
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}
		e := s.PostUser(ctx, user)
		if e != nil {
			return postUserResponse{Err: e}, nil
		}
		token, e := Generate(user)
		if e != nil {
			return postUserResponse{Err: e}, nil
		}
		return postUserResponse{Token: token, ID: user.ID}, nil

	}
}

func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		p, e := s.GetUser(ctx, req.ID)
		return getUserResponse{User: p, Err: e}, nil
	}
}

func MakePatchUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchUserRequest)
		e := s.PatchUser(ctx, req.ID, req.User)
		return patchUserResponse{Err: e, ID: req.ID}, nil
	}
}

func MakeGetUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, e := s.GetUsers(ctx)
		return getUsersResponse{Err: e, Users: users}, nil
	}
}

type postUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type postUserResponse struct {
	Err   error  `json:"err,omitempty"`
	Token string `json:"token,omitempty"`
	ID    string `json:"id,omitempty"`
}

func (r postUserResponse) error() error { return r.Err }

type getUserRequest struct {
	ID string
}

type getUserResponse struct {
	User User  `json:"user,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r getUserResponse) error() error { return r.Err }

type patchUserRequest struct {
	ID   string
	User User
}

type patchUserResponse struct {
	Err error  `json:"err,omitempty"`
	ID  string `json:"id,omitempty"`
}

func (r patchUserResponse) error() error { return r.Err }

type getUsersRequest struct{}

type getUsersResponse struct {
	Err   error  `json:"err,omitempty"`
	Users []User `json:"users"`
}

func (r getUsersResponse) error() error { return r.Err }
