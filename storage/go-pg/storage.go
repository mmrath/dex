package go_pg

import (
	"fmt"
	"github.com/dexidp/dex/storage"
	"github.com/go-pg/pg"
	"strconv"
	"time"
)

type GoPgStorage struct {
	db pg.DB
}

func (s *GoPgStorage) Close() error {
	return s.db.Close()
}

func (s *GoPgStorage) CreateAuthRequest(a storage.AuthRequest) error {
	req := toAuthRequest(&a)
	err := s.db.Insert(req)
	if err != nil {
		return fmt.Errorf("insert auth request: %v", err)
	}
	return nil
}
func (s *GoPgStorage) CreateClient(c storage.Client) error {
	authClient := toAuthClient(&c)
	err := s.db.Insert(authClient)
	if err != nil {
		return fmt.Errorf("insert client: %v", err)
	}
	return nil
}
func (s *GoPgStorage) CreateAuthCode(c storage.AuthCode) error {
	authCode := toAuthCode(&c)
	err := s.db.Insert(authCode)
	if err != nil {
		return fmt.Errorf("insert auth code: %v", err)
	}
	return nil
}
func (s *GoPgStorage) CreateRefresh(r storage.RefreshToken) error {
	rt := toAuthRefreshToken(&r)
	err := s.db.Insert(rt)
	if err != nil {
		return fmt.Errorf("insert refresh token: %v", err)
	}
	return nil
}

func (s *GoPgStorage) CreatePassword(p storage.Password) error {
	pwd := toAuthPassword(&p)
	err := s.db.Insert(pwd)
	if err != nil {
		return fmt.Errorf("insert password: %v", err)
	}
	return nil
}
func (s *GoPgStorage) CreateOfflineSessions(os storage.OfflineSessions) error {
	sessions := toAuthOfflineSessions(&os)
	err := s.db.Insert(sessions)
	if err != nil {
		return fmt.Errorf("insert offline sessions: %v", err)
	}
	return nil
}
func (s *GoPgStorage) CreateConnector(c storage.Connector) error {
	connector := toAuthConnector(&c)
	err := s.db.Insert(connector)
	if err != nil {
		return fmt.Errorf("insert connector: %v", err)
	}
	return nil
}
func (s *GoPgStorage) GetAuthRequest(id string) (req storage.AuthRequest, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		req, err = getAuthRequest(tx, id)
		return err
	})
	return
}

func getAuthRequest(tx *pg.Tx, id string) (storage.AuthRequest, error) {
	au := &UaaAuthRequest{ID: id}
	err := tx.Select(au)

	if err != nil {
		return storage.AuthRequest{}, err
	}
	return toDexAuthRequest(au), nil
}

func (s *GoPgStorage) GetAuthCode(id string) (ac storage.AuthCode, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		ac, err = getAuthCode(tx, id)
		return err
	})
	return
}

func getAuthCode(tx *pg.Tx, id string) (storage.AuthCode, error) {
	au := &UaaAuthCode{ID: id}
	err := tx.Select(au)

	if err != nil {
		return storage.AuthCode{}, err
	}
	return toDexAuthCode(au), nil
}

func (s *GoPgStorage) GetClient(id string) (client storage.Client, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		client, err = getClient(tx, id)
		return err
	})
	return client, err
}

func getClient(tx *pg.Tx, id string) (storage.Client, error) {
	ac := &UaaAuthClient{ID: id}
	err := tx.Select(ac)

	if err != nil {
		return storage.Client{}, err
	}
	return toDexClient(ac), nil
}

func toDexClient(client *UaaAuthClient) storage.Client {
	return storage.Client{
		ID:           client.ID,
		Secret:       client.Secret,
		RedirectURIs: client.RedirectURIs,
		TrustedPeers: client.TrustedPeers,
		Public:       client.Public,
		Name:         client.Name,
		LogoURL:      client.LogoURL,
	}
}

func (s *GoPgStorage) GetKeys() (keys storage.Keys, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		keys, err = getKeys(tx)
		return err
	})
	return
}

func getKeys(tx *pg.Tx) (storage.Keys, error) {
	var authKeys []UaaAuthKey
	err := tx.Select(&authKeys)

	if err != nil {
		return storage.Keys{}, err
	}
	verificationKeys := make([]storage.VerificationKey, len(authKeys)-1)
	idx := 0
	keys := storage.Keys{}
	for _, key := range authKeys {
		if key.Rotated {
			verificationKeys[idx] = storage.VerificationKey{
				PublicKey: key.PublicKey,
				Expiry:    key.ExpiresAt,
			}
			idx = idx + 1
		} else {
			keys.SigningKey = key.SigningKey
			keys.NextRotation = key.NextRotation
			keys.SigningKeyPub = key.PublicKey
			keys.VerificationKeys = verificationKeys
		}
	}
	if idx != len(authKeys) {
		return storage.Keys{}, fmt.Errorf("keys in an inconsistent state")
	}

	return keys, nil
}

func (s *GoPgStorage) GetRefresh(id string) (rt storage.RefreshToken, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		rt, err = getRefresh(tx, id)
		return err
	})
	return
}

func getRefresh(tx *pg.Tx, id string) (storage.RefreshToken, error) {
	rt := UaaAuthRefreshToken{ID: id}
	err := tx.Select(&rt)
	if err != nil {
		return storage.RefreshToken{}, err
	}
	return toDexRefreshToken(&rt), err
}

func toDexRefreshToken(rt *UaaAuthRefreshToken) storage.RefreshToken {
	return storage.RefreshToken{
		ID:    rt.ID,
		Token: rt.Token,

		CreatedAt:     rt.CreatedAt,
		LastUsed:      rt.LastUsedAt,
		ClientID:      rt.ClientID,
		ConnectorID:   rt.ConnectorID,
		ConnectorData: rt.ConnectorData,
		Claims: storage.Claims{
			UserID:        rt.ClaimUserID,
			Username:      rt.ClaimUsername,
			Name:          rt.ClaimName,
			Email:         rt.ClaimEmail,
			EmailVerified: rt.ClaimEmailVerified,
			Groups:        rt.ClaimGroups,
		},
		Scopes: rt.Scopes,
		Nonce:  rt.Nonce,
	}
}
func (s *GoPgStorage) GetPassword(email string) (password storage.Password, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		password, err = getPassword(tx, email)
		return err
	})
	return
}

func getPassword(tx *pg.Tx, email string) (storage.Password, error) {
	password := UaaAuthPassword{}
	err := tx.Model(&password).Where("email = ?", email).Select()
	if err != nil {
		return storage.Password{}, err
	}
	return toDexPassword(password), nil
}

func toDexPassword(password UaaAuthPassword) storage.Password {
	return storage.Password{
		Email:    password.Email,
		Hash:     password.Hash,
		Username: password.Name,
		UserID:   password.UserID,
	}
}
func (s *GoPgStorage) GetOfflineSessions(userID string, connID string) (sessions storage.OfflineSessions, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		sessions, err = getOfflineSessions(tx, userID, connID)
		return err
	})
	return
}

func getOfflineSessions(tx *pg.Tx, userID string, connID string) (storage.OfflineSessions, error) {
	var offlineSessions []UaaAuthOfflineSession
	err := tx.Model(&offlineSessions).Where("user_id = ? and connector_id = ?", userID, connID).Select()
	if err != nil {
		return storage.OfflineSessions{}, err
	}
	return toDexOfflineSessions(connID, connID, offlineSessions), nil
}

func toDexOfflineSessions(userID string, connID string, sessions []UaaAuthOfflineSession) storage.OfflineSessions {
	os := storage.OfflineSessions{
		UserID: userID,
		ConnID: connID,
	}
	os.Refresh = map[string]*storage.RefreshTokenRef{}
	for _, session := range sessions {
		os.Refresh[session.ClientID] = &storage.RefreshTokenRef{
			ID:        session.Token,
			ClientID:  session.ClientID,
			CreatedAt: session.CreatedAt,
			LastUsed:  session.LastUsedAt,
		}
	}
	return os
}
func (s *GoPgStorage) GetConnector(id string) (connector storage.Connector, err error) {
	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		connector, err = getConnector(tx, id)
		return err
	})
	return
}

func getConnector(tx *pg.Tx, id string) (storage.Connector, error) {
	authConn := UaaAuthConnector{ID: id}
	err := tx.Select(&authConn)

	if err != nil {
		return storage.Connector{}, err
	}
	return toDexConnector(authConn), nil
}

func toDexConnector(connector UaaAuthConnector) storage.Connector {
	return storage.Connector{
		ID:              connector.ID,
		Type:            connector.Type,
		Name:            connector.Name,
		ResourceVersion: connector.ResourceVersion,
		Config:          connector.Config,
	}
}

func (s *GoPgStorage) ListClients() ([]storage.Client, error) {
	var authClients []UaaAuthClient
	err := s.db.Select(authClients)
	return toDexClients(authClients), err
}

func toDexClients(authClients []UaaAuthClient) []storage.Client {
	clients := make([]storage.Client, len(authClients))
	for idx, authClient := range authClients {
		clients[idx] = toDexClient(&authClient)
	}
	return clients
}

func (s *GoPgStorage) ListRefreshTokens() ([]storage.RefreshToken, error) {
	var authRefreshTokens []UaaAuthRefreshToken

	err := s.db.Select(authRefreshTokens)
	if err != nil {
		return nil, err
	}
	refreshTokens := make([]storage.RefreshToken, len(authRefreshTokens))
	for idx, t := range authRefreshTokens {
		refreshTokens[idx] = toDexRefreshToken(&t)
	}
	return refreshTokens, nil
}

func (s *GoPgStorage) ListPasswords() ([]storage.Password, error) {
	var authPasswords []UaaAuthPassword
	err := s.db.Select(authPasswords)

	if err != nil {
		return nil, err
	}
	passwords := make([]storage.Password, len(authPasswords))
	for idx, t := range authPasswords {
		passwords[idx] = toDexPassword(t)
	}
	return passwords, nil
}

func (s *GoPgStorage) ListConnectors() ([]storage.Connector, error) {
	var authConnectors []UaaAuthConnector

	err := s.db.Select(authConnectors)

	if err != nil {
		return nil, err
	}
	connectors := make([]storage.Connector, len(authConnectors))
	for idx, t := range authConnectors {
		connectors[idx] = toDexConnector(t)
	}
	return connectors, nil
}
func (s *GoPgStorage) DeleteAuthRequest(id string) error {
	return s.db.Delete(&UaaAuthRequest{ID: id})
}
func (s *GoPgStorage) DeleteAuthCode(code string) error {
	return s.db.Delete(&UaaAuthCode{ID: code})
}
func (s *GoPgStorage) DeleteClient(id string) error {
	return s.db.Delete(&UaaAuthClient{ID: id})
}
func (s *GoPgStorage) DeleteRefresh(id string) error {
	return s.db.Delete(&UaaAuthRefreshToken{ID: id})
}
func (s *GoPgStorage) DeletePassword(email string) error {
	return s.db.Delete(&UaaAuthPassword{Email: email})
}
func (s *GoPgStorage) DeleteOfflineSessions(userID string, connID string) error {
	return s.db.Delete(&UaaAuthOfflineSession{UserID: userID, ConnectorID: connID})
}
func (s *GoPgStorage) DeleteConnector(id string) error {
	return s.db.Delete(&UaaAuthConnector{ID: id})
}
func (s *GoPgStorage) UpdateClient(id string, updater func(old storage.Client) (storage.Client, error)) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		client, err := getClient(tx, id)
		if err != nil {
			return err
		}
		client, err = updater(client)
		if err != nil {
			return err
		}
		err = tx.Update(client)
		return err
	})
}
func (s *GoPgStorage) UpdateKeys(updater func(old storage.Keys) (storage.Keys, error)) (err error) {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		keys, err := getKeys(tx)
		if err != nil {
			return err
		}
		keys, err = updater(keys)
		if err != nil {
			return err
		}
		uaaAuthKeys := toAuthKey(keys)
		_, err = tx.Model(UaaAuthKey{}).Delete()
		if err != nil {
			return err
		}
		err = tx.Insert(uaaAuthKeys)
		return err
	})
}

func toAuthKey(keys storage.Keys) []UaaAuthKey {
	authKeys := make([]UaaAuthKey, len(keys.VerificationKeys)+1)
	authKeys[0] = UaaAuthKey{
		ID:           "0",
		SigningKey:   keys.SigningKey,
		PublicKey:    keys.SigningKeyPub,
		Rotated:      false,
		NextRotation: keys.NextRotation,
	}

	for idx, k := range keys.VerificationKeys {
		authKeys[idx+1] = UaaAuthKey{
			ID:        strconv.Itoa(idx + 1),
			PublicKey: k.PublicKey,
			Rotated:   true,
			ExpiresAt: k.Expiry,
		}
	}
	return authKeys
}
func (s *GoPgStorage) UpdateAuthRequest(id string, updater func(a storage.AuthRequest) (storage.AuthRequest, error)) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		ar, err := getAuthRequest(tx, id)
		if err != nil {
			return err
		}
		ar, err = updater(ar)
		if err != nil {
			return err
		}
		return tx.Update(toAuthRequest(&ar))
	})

}
func (s *GoPgStorage) UpdateRefreshToken(id string, updater func(r storage.RefreshToken) (refresh storage.RefreshToken, err error)) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		refresh, err := getRefresh(tx, id)
		if err != nil {
			return err
		}
		refresh, err = updater(refresh)
		if err != nil {
			return err
		}
		err = tx.Update(toAuthRefreshToken(&refresh))
		return err
	})
}
func (s *GoPgStorage) UpdatePassword(email string, updater func(p storage.Password) (storage.Password, error)) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		password, err := getPassword(tx, email)
		if err != nil {
			return err
		}
		password, err = updater(password)
		if err != nil {
			return err
		}
		err = tx.Update(toAuthPassword(&password))
		return err
	})
}
func (s *GoPgStorage) UpdateOfflineSessions(userID string, connID string, updater func(s storage.OfflineSessions) (storage.OfflineSessions, error)) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		session, err := getOfflineSessions(tx, userID, connID)
		if err != nil {
			return err
		}
		session, err = updater(session)
		if err != nil {
			return err
		}
		err = tx.Update(toAuthOfflineSessions(&session))
		return err
	})
}
func (s *GoPgStorage) UpdateConnector(id string, updater func(c storage.Connector) (storage.Connector, error)) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		connector, err := getConnector(tx, id)
		if err != nil {
			return err
		}
		connector, err = updater(connector)
		if err != nil {
			return err
		}
		conn := toAuthConnector(&connector)
		err = tx.Update(&conn)
		return err
	})
}
func (s *GoPgStorage) GarbageCollect(now time.Time) (result storage.GCResult, err error) {

	err = s.db.RunInTransaction(func(tx *pg.Tx) error {
		r, err := tx.Model(&UaaAuthRequest{}).Where("expiry < ?", now).Delete()
		if err != nil {
			return fmt.Errorf("gc auth_request: %v", err)
		}
		result.AuthRequests = int64(r.RowsAffected())

		r, err = tx.Model(&UaaAuthCode{}).Where("expiry < ?", now).Delete()
		if err != nil {
			return fmt.Errorf("gc auth_code: %v", err)
		}
		result.AuthCodes = int64(r.RowsAffected())
		return nil
	})

	return
}
