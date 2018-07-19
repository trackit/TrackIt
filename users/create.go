//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package users

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"net/http"
	
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit2/config"
	"github.com/trackit/trackit2/db"
	"github.com/trackit/trackit2/models"
	"github.com/trackit/trackit2/routes"
)

const (
	passwordMaxLength = 12
)

var (
	ErrPasswordTooShort = errors.New(fmt.Sprintf("Password must be at least %u characters.", passwordMaxLength))
)

func init() {
	routes.MethodMuxer{
		http.MethodPost: routes.H(createUser).With(
			routes.RequestContentType{"application/json"},
			routes.RequestBody{createUserRequestBody{"example@example.com", "pa55w0rd"}},
			routes.Documentation{
				Summary:     "register a new user",
				Description: "Registers a new user using an e-mail and password, and responds with the user's data.",
			},
		),
		http.MethodPatch: routes.H(patchUser).With(
			RequireAuthenticatedUser{ViewerAsSelf},
			routes.RequestContentType{"application/json"},
			routes.RequestBody{createUserRequestBody{"example@example.com", "pa55w0rd"}},
			routes.Documentation{
				Summary:     "edit the current user",
				Description: "Edit the current user, and responds with the user's data.",
			},
		),
		http.MethodGet: routes.H(me).With(
			RequireAuthenticatedUser{ViewerAsSelf},
			routes.Documentation{
				Summary:     "get the current user",
				Description: "Responds with the currently authenticated user's data.",
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
		routes.Documentation{
			Summary: "register or get the user",
		},
	).Register("/user")

	routes.MethodMuxer{
		http.MethodPost: routes.H(createViewerUser).With(
			routes.RequestContentType{"application/json"},
			RequireAuthenticatedUser{ViewerCannot},
			routes.RequestBody{createViewerUserRequestBody{"example@example.com"}},
			routes.Documentation{
				Summary:     "register a new viewer user",
				Description: "Registers a new viewer user linked to the current user, which will only be able to view its parent user's data.",
			},
		),
		http.MethodGet: routes.H(getViewerUsers).With(
			RequireAuthenticatedUser{ViewerAsParent},
			routes.Documentation{
				Summary:     "list viewer users",
				Description: "Lists the viewer users registered for the current account.",
			},
		),
	}.H().With(
		db.RequestTransaction{db.Db},
	).Register("/user/viewer")
}

type createUserRequestBody struct {
	Email    string `json:"email"    req:"nonzero"`
	Password string `json:"password" req:"nonzero"`
}

func createUser(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body createUserRequestBody
	routes.MustRequestBody(a, &body)
	tx := a[db.Transaction].(*sql.Tx)
	code, resp := createUserWithValidBody(request, body, tx)
	// Add the default role to the new account. No error is returned in case of failure
	// The billing repository is not processed instantly
	if code == 200 && config.DefaultRole != "" && config.DefaultRoleName != "" &&
		config.DefaultRoleExternal != "" && config.DefaultRoleBucket != "" {
		addDefaultRole(request, resp.(User), tx)
	}
	return code, resp
}

func createUserWithValidBody(request *http.Request, body createUserRequestBody, tx *sql.Tx) (int, interface{}) {
	ctx := request.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	user, err := CreateUserWithPassword(ctx, tx, body.Email, body.Password)
	if err == nil {
		logger.Info("User created.", user)
		return 200, user
	} else {
		logger.Error(err.Error(), nil)
		errSplit := strings.Split(err.Error(), ":")
		if (len(errSplit) >= 1 && errSplit[0] == "Error 1062") {
			return 409, errors.New("Account already exists.")
		} else {
			return 500, errors.New("Failed to create user.")
		}
	}
}

type createViewerUserRequestBody struct {
	Email string `json:"email"    req:"nonzero"`
}

type createViewerUserResponseBody struct {
	User
	Password string `json:"password" req:"nonzero"`
}

func createViewerUser(request *http.Request, a routes.Arguments) (int, interface{}) {
	var body createViewerUserRequestBody
	routes.MustRequestBody(a, &body)
	currentUser := a[AuthenticatedUser].(User)
	tx := a[db.Transaction].(*sql.Tx)
	ctx := request.Context()
	viewerUser, viewerUserPassword, err := CreateUserWithParent(ctx, tx, body.Email, currentUser)
	if err != nil {
		errSplit := strings.Split(err.Error(), ":")
		if (len(errSplit) >= 1 && errSplit[0] == "Error 1062") {
			return 409, errors.New("Email already taken.")
		} else {
			return 500, errors.New("Failed to create viewer user.")
		}
	}
	response := createViewerUserResponseBody{
		User:     viewerUser,
		Password: viewerUserPassword,
	}
	return http.StatusOK, response
}

func getViewerUsers(request *http.Request, a routes.Arguments) (int, interface{}) {
	currentUser := a[AuthenticatedUser].(User)
	tx := a[db.Transaction].(*sql.Tx)
	ctx := request.Context()
	users, err := GetUsersByParent(ctx, tx, currentUser)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed to get viewer users.")
	}
	return http.StatusOK, users
}

func addDefaultRole(request *http.Request, user User, tx *sql.Tx) {
	ctx := request.Context()
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	accoundDB := models.AwsAccount{
		UserID:   user.Id,
		Pretty:   config.DefaultRoleName,
		RoleArn:  config.DefaultRole,
		External: config.DefaultRoleExternal,
	}
	err := accoundDB.Insert(tx)
	if err != nil {
		logger.Error("Failed to add default role", err)
	} else {
		brDB := models.AwsBillRepository{
			AwsAccountID: accoundDB.ID,
			Bucket:       config.DefaultRoleBucket,
			Prefix:       config.DefaultRoleBucketPrefix,
		}
		err = brDB.Insert(tx)
		if err != nil {
			logger.Error("Failed to add default bill repository", err)
		}
	}
}
