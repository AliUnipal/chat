package usercontroller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

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
	Authenticate(ctx context.Context, username, password string) (usersvc.AuthenticateOutput, error)
}

type userAccessor interface {
	CurrentUserID(ctx context.Context) (uuid.UUID, bool)
}

func New(userService userService, userAcc userAccessor) *userController {
	return &userController{userAcc, userService}
}

type userController struct {
	userAcc     userAccessor
	userService userService
}

type validationErrors map[string]string

func (v validationErrors) Error() string {
	return "validation errors"
}

type (
	CreateUserRequest struct {
		ImageURL  string `json:"imageURL"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName,omitempty"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}
	CreateUserResponse struct {
		UserID uuid.UUID `json:"userID"`
	}
)

func (c *CreateUserRequest) decodeAndValidate(ctx context.Context, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		slog.ErrorContext(ctx, "failed to decode request CreateUserRequest", "error", err)
		return err
	}

	errs := validationErrors{}

	if c.ImageURL == "" {
		errs["imageURL"] = "required"
	}
	if c.FirstName == "" {
		errs["firstName"] = "required"
	}
	if c.Username == "" {
		errs["username"] = "required"
	}
	if c.Password == "" {
		errs["password"] = "required"
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

	err := req.decodeAndValidate(ctx, r)
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
		Password:  req.Password,
	})
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	httpsrv.RespondWithJSON(ctx, w, CreateUserResponse{userID})
}

type (
	GetUserWithIDRequest struct {
		ID         string `http_path:"id"`
		parsedUUID uuid.UUID
	}
	GetUserWithIDResponse struct {
		ID        uuid.UUID `json:"id"`
		ImageURL  string    `json:"imageURL"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Username  string    `json:"username"`
	}
)

func (r *GetUserWithIDRequest) validate(ctx context.Context) error {
	if r.ID == "" {
		slog.ErrorContext(ctx, "id is required", "error", "required")
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

func (u *userController) GetUserWithID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GetUserWithIDRequest{ID: r.PathValue("id")}

	if err := req.validate(ctx); err != nil {
		httpsrv.RespondWithBadRequestError(ctx, w, errcodes.InvalidUUID, err)
		return
	}

	usr, err := u.userService.GetUser(ctx, req.parsedUUID)
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	httpsrv.RespondWithJSON(ctx, w, GetUserWithIDResponse(usr))
}

type (
	GetUserRequest struct {
		ID         string `http_path:"id"`
		parsedUUID uuid.UUID
	}
	userResponse struct {
		ID        uuid.UUID `json:"id"`
		ImageURL  string    `json:"imageURL"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Username  string    `json:"username"`
	}
	GetUserResponse struct {
		Data userResponse `json:"data"`
	}
)

func (u *userController) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usrID, ok := u.userAcc.CurrentUserID(ctx)
	if !ok {
		httpsrv.RespondWithAuthenticationError(ctx, w, errcodes.Unauthenticated)
		return
	}

	usr, err := u.userService.GetUser(ctx, usrID)
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	httpsrv.RespondWithJSON(ctx, w, GetUserResponse{
		Data: userResponse(usr),
	})
}

type (
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	loginUserResponse struct {
		ID        uuid.UUID `json:"id"`
		ImageURL  string    `json:"imageURL"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Username  string    `json:"username"`
		Token     string    `json:"token"`
	}
	LoginResponse struct {
		Data loginUserResponse `json:"data"`
	}
)

func (l *LoginRequest) decodeAndValidate(ctx context.Context, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		slog.ErrorContext(ctx, "failed to decode request LoginRequest", "error", err)
		return err
	}
	errs := validationErrors{}

	if l.Username == "" {
		errs["username"] = "required"
	}
	if l.Password == "" {
		errs["password"] = "required"
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (u *userController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := LoginRequest{}

	if err := req.decodeAndValidate(ctx, r); err != nil {
		var fieldErrs validationErrors
		if errors.As(err, &fieldErrs) {
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

	out, err := u.userService.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		httpsrv.RespondWithError(ctx, w, errcodes.ServiceError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    out.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		MaxAge:   int(10 * time.Hour.Seconds()),
	})
	httpsrv.RespondWithJSON(ctx, w, LoginResponse{
		Data: loginUserResponse{
			ID:        out.User.ID,
			ImageURL:  out.User.ImageURL,
			FirstName: out.User.FirstName,
			LastName:  out.User.LastName,
			Username:  out.User.Username,
			Token:     out.Token,
		},
	})
}
