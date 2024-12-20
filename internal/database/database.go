package database

import (
	"errors"

	"github.com/kamuridesu/ip-syncer/internal/shared"
)

type Database interface {
	Connect() error
	Disconnect() error
	CreateIfNotexists() error
	Insert(info *shared.IPInfo) error
	Update(info *shared.IPInfo) error
	Delete(id string) error
	DeleteByName(name string) error
	GetByIP(ip string) (*shared.IPInfo, error)
	GetByName(name string) (*shared.IPInfo, error)
	GetAll() (*[]shared.IPInfo, error)
}

func New(dbType string, info string) (Database, error) {
	switch dbType {
	case "sqlite":
		return NewSQLite(info)
	default:
		return nil, errors.New("unsupported database type")
	}
}
