/*
   GoToSocial
   Copyright (C) 2021 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package federation

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/superseriousbusiness/activity/streams"
	"github.com/superseriousbusiness/gotosocial/internal/db"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

func (p *processor) GetStatus(ctx context.Context, requestedUsername string, requestedStatusID string, requestURL *url.URL) (interface{}, gtserror.WithCode) {
	// get the account the request is referring to
	requestedAccount, err := p.db.GetLocalAccountByUsername(ctx, requestedUsername)
	if err != nil {
		return nil, gtserror.NewErrorNotFound(fmt.Errorf("database error getting account with username %s: %s", requestedUsername, err))
	}

	// authenticate the request
	requestingAccountURI, authenticated, err := p.federator.AuthenticateFederatedRequest(ctx, requestedUsername)
	if err != nil || !authenticated {
		return nil, gtserror.NewErrorNotAuthorized(errors.New("not authorized"), "not authorized")
	}

	requestingAccount, _, err := p.federator.GetRemoteAccount(ctx, requestedUsername, requestingAccountURI, false)
	if err != nil {
		return nil, gtserror.NewErrorNotAuthorized(err)
	}

	// authorize the request:
	// 1. check if a block exists between the requester and the requestee
	blocked, err := p.db.IsBlocked(ctx, requestedAccount.ID, requestingAccount.ID, true)
	if err != nil {
		return nil, gtserror.NewErrorInternalError(err)
	}

	if blocked {
		return nil, gtserror.NewErrorNotAuthorized(fmt.Errorf("block exists between accounts %s and %s", requestedAccount.ID, requestingAccount.ID))
	}

	// get the status out of the database here
	s := &gtsmodel.Status{}
	if err := p.db.GetWhere(ctx, []db.Where{
		{Key: "id", Value: requestedStatusID},
		{Key: "account_id", Value: requestedAccount.ID},
	}, s); err != nil {
		return nil, gtserror.NewErrorNotFound(fmt.Errorf("database error getting status with id %s and account id %s: %s", requestedStatusID, requestedAccount.ID, err))
	}

	visible, err := p.filter.StatusVisible(ctx, s, requestingAccount)
	if err != nil {
		return nil, gtserror.NewErrorInternalError(err)
	}
	if !visible {
		return nil, gtserror.NewErrorNotFound(fmt.Errorf("status with id %s not visible to user with id %s", s.ID, requestingAccount.ID))
	}

	// requester is authorized to view the status, so convert it to AP representation and serialize it
	asStatus, err := p.tc.StatusToAS(ctx, s)
	if err != nil {
		return nil, gtserror.NewErrorInternalError(err)
	}

	data, err := streams.Serialize(asStatus)
	if err != nil {
		return nil, gtserror.NewErrorInternalError(err)
	}

	return data, nil
}