package usercontroller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AliUnipal/chat/internal/models/user"
	"github.com/AliUnipal/chat/internal/server"
	"github.com/AliUnipal/chat/internal/server/httpsrv"
	"github.com/AliUnipal/chat/internal/service/usersvc"
	"github.com/AliUnipal/chat/pkg/errcodes"
	"github.com/google/uuid"
)

type userService interface {
	CreateUser(ctx context.Context, in usersvc.CreateUserInput) (uuid.UUID, error)
	GetUser(ctx context.Context, id uuid.UUID) (user.User, error)
}

type userController struct {
	userService userService
}

func New(userService userService) *userController {
	return &userController{userService}
}

type (
	CreateUserRequest struct {
		ImageURL  string `json:"imageURL,omitempty"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Username  string `json:"username"`
	}
	validationErrors   map[string]string
	CreateUserResponse struct {
		UserID uuid.UUID `json:"userID"`
	}
)

func (v validationErrors) Error() string {
	return "validation errors"
}

func (c *CreateUserRequest) validateAndDecode(ctx context.Context, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&r); err != nil {
		slog.ErrorContext(ctx, "failed to decode request CreateUserRequest", "error", err)
		return err
	}

	errs := validationErrors{}

	if c.FirstName == "" {
		errs["firstName"] = "required"
	}
	if c.LastName == "" {
		errs["lastName"] = "required"
	}
	if c.Username == "" {
		errs["username"] = "required"
	}

	if len(errs) > 0 {
		slog.ErrorContext(ctx, "validation errors", "errors", errs)
		return errs
	}

	return nil
}

func (u *userController) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateUserRequest

	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	slog.ErrorContext(ctx, "failed to decode request CreateUserRequest", "error", err)
	//	httpsrv.RespondWithError(ctx, w, errcodes.InvalidInput, err)
	//	return
	//}

	err := req.validateAndDecode(ctx, r)
	if err != nil {
		if fieldErrs, ok := err.(validationErrors); ok {
			var kvs httpsrv.KVs

			for field, errMsg := range fieldErrs {
				kvs = append(kvs, httpsrv.KV(field, errMsg))
			}

			httpsrv.RespondWithValidationError(ctx, w, errcodes.InvalidInput, kvs...)
			return
		}

		httpsrv.RespondWithError(ctx, w, errcodes.InvalidInput, err)
		return
	}

	userID, err := u.userService.CreateUser(ctx, usersvc.CreateUserInput{
		ImageURL:  req.ImageURL,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
	})
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	httpsrv.RespondWithJSON(ctx, w, CreateUserResponse{userID})
}

type (
	GetUserRequest struct {
		ID         string `http_path:"id"`
		parsedUUID uuid.UUID
	}
	GetUserResponse user.User
)

func (r *GetUserRequest) validate(ctx context.Context) error {
	if r.ID == "" {
		slog.ErrorContext(ctx, "id is required", "error")
		return errors.New("id is required")
	}

	chatId, err := uuid.Parse(r.ID)
	if err != nil {
		slog.ErrorContext(ctx, "invalid uuid", "error", err)
		return server.ErrInvalidUUID
	}

	r.parsedUUID = chatId
	return nil
}

func (u *userController) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GetUserRequest{ID: r.PathValue("id")}

	if err := req.validate(ctx); err != nil {
		httpsrv.RespondWithBadRequestError(ctx, w, errcodes.InvalidUUID, err)
		return
	}

	usr, err := u.userService.GetUser(ctx, req.parsedUUID)
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	httpsrv.RespondWithJSON(ctx, w, GetUserResponse(usr))
}
