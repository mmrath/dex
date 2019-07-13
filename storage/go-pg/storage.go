package go_pg

import (
	"github.com/dexidp/dex/storage"
	"time"
	"github.com/go-pg/pg"
)

type GoPgStorage struct {
	db pg.DB
}

func (s *GoPgStorage) Close() error {
	return s.db.Close()
}
func (s *GoPgStorage) CreateAuthRequest(a storage.AuthRequest) error {
	return nil
}
func (s *GoPgStorage) CreateClient(c storage.Client) error {
	return nil
}
func (s *GoPgStorage) CreateAuthCode(c storage.AuthCode) error {
	return nil
}
func (s *GoPgStorage) CreateRefresh(r storage.RefreshToken) error {
	return nil
}
func (s *GoPgStorage) CreatePassword(p storage.Password) error {
	return nil
}
func (s *GoPgStorage) CreateOfflineSessions(s storage.OfflineSessions) error {
	return nil
}
func (s *GoPgStorage) CreateConnector(c storage.Connector) error {
	return nil
}
func (s *GoPgStorage) GetAuthRequest(id string) (storage.AuthRequest, error) {
	return storage.AuthRequest{}, nil
}
func (s *GoPgStorage) GetAuthCode(id string) (storage.AuthCode, error) {
	return storage.AuthCode{}, nil
}
func (s *GoPgStorage) GetClient(id string) (storage.Client, error) {
	return storage.Client{}, nil
}
func (s *GoPgStorage) GetKeys() (storage.Keys, error) {
	return storage.Keys{}, nil
}
func (s *GoPgStorage) GetRefresh(id string) (storage.RefreshToken, error) {
	return storage.RefreshToken{}, nil
}
func (s *GoPgStorage) GetPassword(email string) (storage.Password, error) {
	return storage.Password{}, nil
}
func (s *GoPgStorage) GetOfflineSessions(userID string, connID string) (storage.OfflineSessions, error) {
	return storage.OfflineSessions{}, nil
}
func (s *GoPgStorage) GetConnector(id string) (storage.Connector, error) {
	return storage.Connector{}, nil
}
func (s *GoPgStorage) ListClients() ([]storage.Client, error) {
	return nil, nil
}
func (s *GoPgStorage) ListRefreshTokens() ([]storage.RefreshToken, error) {
	return nil, nil
}
func (s *GoPgStorage) ListPasswords() ([]storage.Password, error) {
	return nil, nil
}
func (s *GoPgStorage) ListConnectors() ([]storage.Connector, error) {
	return nil, nil
}
func (s *GoPgStorage) DeleteAuthRequest(id string) error {
	return nil
}
func (s *GoPgStorage) DeleteAuthCode(code string) error {
	return nil
}
func (s *GoPgStorage) DeleteClient(id string) error {
	return nil
}
func (s *GoPgStorage) DeleteRefresh(id string) error {
	return nil
}
func (s *GoPgStorage) DeletePassword(email string) error {
	return nil
}
func (s *GoPgStorage) DeleteOfflineSessions(userID string, connID string) error {
	return nil
}
func (s *GoPgStorage) DeleteConnector(id string) error {
	return nil
}
func (s *GoPgStorage) UpdateClient(id string, updater func(old storage.Client) (storage.Client, error)) error {
	return nil
}
func (s *GoPgStorage) UpdateKeys(updater func(old storage.Keys) (storage.Keys, error)) error {
	return nil
}
func (s *GoPgStorage) UpdateAuthRequest(id string, updater func(a storage.AuthRequest) (storage.AuthRequest, error)) error {
	return nil
}
func (s *GoPgStorage) UpdateRefreshToken(id string, updater func(r storage.RefreshToken) (storage.RefreshToken, error)) error {
	return nil
}
func (s *GoPgStorage) UpdatePassword(email string, updater func(p storage.Password) (storage.Password, error)) error {
	return nil
}
func (s *GoPgStorage) UpdateOfflineSessions(userID string, connID string, updater func(s storage.OfflineSessions) (storage.OfflineSessions, error)) error {
	return nil
}
func (s *GoPgStorage) UpdateConnector(id string, updater func(c storage.Connector) (storage.Connector, error)) error {
	return nil
}
func (s *GoPgStorage) GarbageCollect(now time.Time) (storage.GCResult, error) {
	return storage.GCResult{}, nil
}
