package go_pg

import (
	"github.com/dexidp/dex/storage"
	"gopkg.in/square/go-jose.v2"
	"time"
)

type UaaAuthRequest struct {
	tableName struct{} `sql:"uaa_auth_request"`

	// ID used to identify the authorization request.
	ID string `sql:"id"`

	// ID of the client requesting authorization from a user.
	ClientID string `sql:"client_id"`

	// Values parsed from the initial request. These describe the resources the client is
	// requesting as well as values describing the form of the response.
	ResponseTypes []string `sql:"response_types,array"`
	Scopes        []string `sql:"scopes,array"`
	RedirectURI   string   `sql:"redirect_uri"`
	Nonce         string   `sql:"nonce"`
	State         string   `sql:"state"`

	// The client has indicated that the end user must be shown an approval prompt
	// on all requests. The server cannot cache their initial action for subsequent
	// attempts.
	ForceApprovalPrompt bool `sql:"force_approval_prompt"`

	Expiry time.Time `sql:"expiry"`

	// Has the user proved their identity through a backing identity provider?
	//
	// If false, the following fields are invalid.
	LoggedIn bool `sql:"logged_in"`

	// The identity of the end user. Generally nil until the user authenticates
	// with a backend.
	ClaimUserID        string `sql:"claim_user_id"`
	ClaimUsername      string `sql:"claim_username"`
	ClaimName          string `sql:"claim_name"`
	ClaimEmail         string `sql:"claim_email"`
	ClaimEmailVerified bool   `sql:"claim_email_verified"`

	ClaimGroups []string `sql:"claim_groups,array"`

	// The connector used to login the user and any data the connector wishes to persists.
	// Set when the user authenticates.
	ConnectorID   string `sql:"connector_id"`
	ConnectorData []byte
}

func toAuthRequest(r *storage.AuthRequest) UaaAuthRequest {
	request := UaaAuthRequest{
		// ID used to identify the authorization request.
		ID: r.ID,

		// ID of the client requesting authorization from a user.
		ClientID: r.ClientID,

		// Values parsed from the initial request. These describe the resources the client is
		// requesting as well as values describing the form of the response.
		ResponseTypes:       r.ResponseTypes,
		Scopes:              r.Scopes,
		RedirectURI:         r.RedirectURI,
		Nonce:               r.Nonce,
		State:               r.State,
		ForceApprovalPrompt: r.ForceApprovalPrompt,
		Expiry:              r.Expiry,
		LoggedIn:            r.LoggedIn,

		ClaimUserID:        r.Claims.UserID,
		ClaimUsername:      r.Claims.Username,
		ClaimName:          r.Claims.Name,
		ClaimEmail:         r.Claims.Email,
		ClaimEmailVerified: r.Claims.EmailVerified,
		ClaimGroups:        r.Claims.Groups,

		// The connector used to login the user and any data the connector wishes to persists.
		// Set when the user authenticates.
		ConnectorID:   r.ConnectorID,
		ConnectorData: r.ConnectorData,
	}

	return request
}

func toDexAuthRequest(req *UaaAuthRequest) storage.AuthRequest {
	return storage.AuthRequest{
		ID:                  req.ID,
		ClientID:            req.ClientID,
		ResponseTypes:       req.ResponseTypes,
		Scopes:              req.Scopes,
		RedirectURI:         req.RedirectURI,
		Nonce:               req.Nonce,
		State:               req.State,
		ForceApprovalPrompt: req.ForceApprovalPrompt,
		Expiry:              req.Expiry,
		LoggedIn:            req.LoggedIn,
		Claims: storage.Claims{
			UserID:        req.ClaimUserID,
			Username:      req.ClaimUsername,
			Name:          req.ClaimName,
			Email:         req.ClaimEmail,
			EmailVerified: req.ClaimEmailVerified,
			Groups:        req.ClaimGroups,
		},
		ConnectorID:   req.ConnectorID,
		ConnectorData: req.ConnectorData,
	}
}

type UaaAuthClient struct {
	tableName struct{} `sql:"uaa_auth_client"`

	// Client ID and secret used to identify the client.
	ID     string `sql:"id"`
	Secret string `sql:"secret"`

	// A registered set of redirect URIs. When redirecting from dex to the client, the URI
	// requested to redirect to MUST match one of these values, unless the client is "public".
	RedirectURIs []string `sql:"redirect_uris,array"`

	// TrustedPeers are a list of peers which can issue tokens on this client's behalf using
	// the dynamic "oauth2:server:client_id:(client_id)" scope. If a peer makes such a request,
	// this client's ID will appear as the ID Token's audience.
	//
	// Clients inherently trust themselves.
	TrustedPeers []string `sql:"trusted_peers,array"`

	// Public clients must use either use a redirectURL 127.0.0.1:X or "urn:ietf:wg:oauth:2.0:oob"
	Public bool `sql:"public"`

	// Name and LogoURL used when displaying this client to the end user.
	Name    string `sql:"name"`
	LogoURL string `sql:"logo_url"`
}

func toAuthClient(c *storage.Client) UaaAuthClient {
	return UaaAuthClient{
		ID:           c.ID,
		Secret:       c.Secret,
		RedirectURIs: c.RedirectURIs,
		TrustedPeers: c.TrustedPeers,
		Public:       c.Public,
		Name:         c.Name,
		LogoURL:      c.LogoURL,
	}
}

type UaaAuthCode struct {
	tableName struct{} `sql:"uaa_auth_code"`

	// Actual string returned as the "code" value.
	ID string `sql:"id"`

	// The client this code value is valid for. When exchanging the code for a
	// token response, the client must use its client_secret to authenticate.
	ClientID string `sql:"client_id"`

	// As part of the OAuth2 spec when a client makes a token request it MUST
	// present the same redirect_uri as the initial redirect. This values is saved
	// to make this check.
	//
	// https://tools.ietf.org/html/rfc6749#section-4.1.3
	RedirectURI string `sql:"redirect_uri"`

	// If provided by the client in the initial request, the provider MUST create
	// a ID Token with this nonce in the JWT payload.
	Nonce string `sql:"nonce"`

	// Scopes authorized by the end user for the client.
	Scopes []string `sql:"scopes,array"`

	// Authentication data provided by an upstream source.
	ConnectorID   string `sql:"connector_id"`
	ConnectorData []byte `sql:"connector_data"`

	ClaimUserID        string `sql:"claim_user_id"`
	ClaimUsername      string `sql:"claim_username"`
	ClaimName          string `sql:"claim_name"`
	ClaimEmail         string `sql:"claim_email"`
	ClaimEmailVerified bool   `sql:"claim_email_verified"`

	Expiry time.Time
}

func toAuthCode(a *storage.AuthCode) UaaAuthCode {
	return UaaAuthCode{
		ID:                 a.ID,
		ClientID:           a.ClientID,
		RedirectURI:        a.RedirectURI,
		Nonce:              a.Nonce,
		Scopes:             a.Scopes,
		ConnectorID:        a.ConnectorID,
		ConnectorData:      a.ConnectorData,
		ClaimUserID:        a.Claims.UserID,
		ClaimUsername:      a.Claims.Username,
		ClaimName:          a.Claims.Name,
		ClaimEmail:         a.Claims.Email,
		ClaimEmailVerified: a.Claims.EmailVerified,
		Expiry:             a.Expiry,
	}
}

func toDexAuthCode(code *UaaAuthCode) storage.AuthCode {
	return storage.AuthCode{
		ID:            code.ID,
		ClientID:      code.ClientID,
		RedirectURI:   code.RedirectURI,
		Nonce:         code.Nonce,
		Scopes:        code.Scopes,
		Expiry:        code.Expiry,
		ConnectorID:   code.ConnectorID,
		ConnectorData: code.ConnectorData,
		Claims: storage.Claims{
			UserID:        code.ClaimUserID,
			Username:      code.ClaimUsername,
			Name:          code.ClaimName,
			Email:         code.ClaimEmail,
			EmailVerified: code.ClaimEmailVerified,
		},
	}
}

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

type UaaAuthPassword struct {
	tableName struct{} `sql:"uaa_auth_password"`

	// Email and identifying name of the password. Emails are assumed to be valid and
	// determining that an end-user controls the address is left to an outside application.
	//
	// Emails are case insensitive and should be standardized by the storage.
	//
	// Storages that don't support an extended character set for IDs, such as '.' and '@'
	// (cough cough, kubernetes), must map this value appropriately.
	Email string `sql:"email"`

	// Bcrypt encoded hash of the password. This package enforces a min cost value of 10
	Hash []byte `json:"hash"`

	// Optional username to display. NOT used during login.
	Name string `sql:"name"`

	// Randomly generated user ID. This is NOT the primary ID of the Password object.
	UserID string `sql:"user_id"`
}

func toAuthPassword(p *storage.Password) UaaAuthPassword {
	return UaaAuthPassword{
		Email:  p.Email,
		Hash:   p.Hash,
		Name:   p.Username,
		UserID: p.UserID,
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

// OfflineSessions objects are sessions pertaining to users with refresh tokens.
type UaaAuthConnector struct {
	tableName struct{} `sql:"uaa_auth_connector"`

	// ID that will uniquely identify the connector object.
	ID string `sql:"id"`
	// The Type of the connector. E.g. 'oidc' or 'ldap'
	Type string `sql:"type"`
	// The Name of the connector that is used when displaying it to the end user.
	Name string `sql:"name"`
	// ResourceVersion is the static versioning used to keep track of dynamic configuration
	// changes to the connector object made by the API calls.
	ResourceVersion string `sql:"resource_version"`

	Config []byte `json:"config"`
}

func toAuthConnector(c *storage.Connector) UaaAuthConnector {
	return UaaAuthConnector{
		ID:              c.ID,
		Type:            c.Type,
		Name:            c.Name,
		ResourceVersion: c.ResourceVersion,
		Config:          c.Config,
	}
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
