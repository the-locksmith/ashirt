// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type UpdateUserProfileInput struct {
	UserSlug  string
	FirstName string
	LastName  string
	Email     string
}

func UpdateUserProfile(ctx context.Context, db *database.Connection, i UpdateUserProfileInput) error {
	var userID int64
	var err error

	if userID, err = selfOrSlugToUserID(ctx, db, i.UserSlug); err != nil {
		return backend.DatabaseErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyUser{UserID: userID}); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	err = db.Update(sq.Update("users").
		SetMap(map[string]interface{}{
			"first_name": i.FirstName,
			"last_name":  i.LastName,
			"email":      i.Email,
		}).
		Where(sq.Eq{"id": userID}))

	if err != nil {
		return backend.DatabaseErr(err)
	}
	return nil
}
