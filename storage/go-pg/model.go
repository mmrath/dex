
package go_pg

import (
	"github.com/dexidp/dex/storage"
	"gopkg.in/square/go-jose.v2"
	"time"
)


type UaaAuthRefreshToken struct {
	tableName struct{} `sql:"uaa_auth_refresh_token"`

	ID string `sql:"id"`

	// A single token that's rotated every time the refresh token is refreshed.
	//
	// May be empty.
	Token string `sql:"token"`

	CreatedAt  time.Time `sql:"created_at"`
	LastUsedAt time.Time `sql:"last_used_at"`

	// Client this refresh token is valid for.
	ClientID string `sql:"client_id"`

	// Authentication data provided by an upstream source.
	ConnectorID        string   `sql:"connector_id"`
	ConnectorData      []byte   `sql:"connector_data"`
	Scopes             []string `sql:"scopes"`
	Nonce              string   `sql:"nonce"`
	ClaimUserID        string   `sql:"claim_user_id"`
	ClaimUsername      string   `sql:"claim_username"`
	ClaimName          string   `sql:"claim_name"`
	ClaimEmail         string   `sql:"claim_email"`
	ClaimEmailVerified bool     `sql:"claim_email_verified"`
	ClaimGroups        []string `sql:"claim_groups"`
}

func toAuthRefreshToken(r *storage.RefreshToken) UaaAuthRefreshToken {
	return UaaAuthRefreshToken{
		ID:                 r.ID,
		Token:              r.Token,
		CreatedAt:          r.CreatedAt,
		LastUsedAt:         r.LastUsed,
		ClientID:           r.ClientID,
		ConnectorID:        r.ConnectorID,
		ConnectorData:      r.ConnectorData,
		ClaimUserID:        r.Claims.UserID,
		ClaimUsername:      r.Claims.Username,
		ClaimName:          r.Claims.Name,
		ClaimEmail:         r.Claims.Email,
		ClaimEmailVerified: r.Claims.EmailVerified,
		ClaimGroups:        r.Claims.Groups,
		Scopes:             r.Scopes,
		Nonce:              r.Nonce,
	}
}

// OfflineSessions objects are sessions pertaining to users with refresh tokens.
type UaaAuthOfflineSession struct {
	// UserID of an end user who has logged in to the server.
	UserID string `sql:"user_id"`
	// The ID of the connector used to login the user.
	ConnectorID string    `sql:"connector_id"`
	ClientID    string    `sql:"client_id"`
	Token       string    `sql:"token"`
	CreatedAt   time.Time `sql:"created_at"`
	LastUsedAt  time.Time `sql:"last_used_at"`
}

func toAuthOfflineSessions(s *storage.OfflineSessions) []UaaAuthOfflineSession {

	length := len(s.Refresh)
	sessions := make([]UaaAuthOfflineSession, length)

	index := 0
	for _, v := range s.Refresh {
		session := UaaAuthOfflineSession{
			UserID:      s.UserID,
			ConnectorID: s.ConnID,
			ClientID:    v.ClientID,
			Token:       v.ID,
			CreatedAt:   v.CreatedAt,
			LastUsedAt:  v.LastUsed,
		}
		sessions[index] = session
		index = index + 1
	}

	return sessions
}

// AuthKeys hold encryption and signing keys.
type UaaAuthKey struct {
	tableName struct{} `sql:"uaa_auth_key"`

	ID           string           `sql:"id"`
	SigningKey   *jose.JSONWebKey `sql:"signing_key"`
	PublicKey    *jose.JSONWebKey `sql:"public_key"`
	Rotated      bool             `sql:"rotated"`
	NextRotation time.Time        `sql:"next_rotation"`
	ExpiresAt    time.Time        `sql:"expires_at"`
}
