package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/kamuridesu/ip-syncer/internal/shared"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	DB *sql.DB
}

func NewSQLite(filename string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to open database: %s", err.Error()))
		return nil, err
	}
	return &SQLite{DB: db}, nil
}

func (s *SQLite) Connect() error {
	return s.DB.Ping()
}

func (s *SQLite) Disconnect() error {
	return s.DB.Close()
}

func (s *SQLite) CreateIfNotexists() error {
	_, err := s.DB.Exec(`CREATE TABLE IF NOT EXISTS "ips" (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT, 
		"ip" TEXT, 
		"name" TEXT, 
		"authkey" TEXT, changed BOOLEAN)
	`)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create table: %s", err.Error()))
		return err
	}
	return nil
}

func (s *SQLite) Insert(info *shared.IPInfo) error {
	_, err := s.DB.Exec(`INSERT INTO ips (
		ip,
		name,
		authkey,
		changed
	) VALUES (?, ?, ?, ?)`, info.IP, info.Name, info.Authkey, info.Changed)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to insert record: %s", err.Error()))
		return err
	}
	return nil
}

func (s *SQLite) Update(info *shared.IPInfo) error {
	_, err := s.DB.Exec(`UPDATE ips SET
		ip = ?,
		name = ?,
		authkey = ?,
		changed = ?
	WHERE id = ?`, info.IP, info.Name, info.Authkey, info.Changed)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to update record: %s", err.Error()))
		return err
	}
	return nil
}

func (s *SQLite) Delete(id string) error {
	_, err := s.DB.Exec(`DELETE FROM ips WHERE id = ?`, id)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to delete record: %s", err.Error()))
		return err
	}
	return nil
}

func (s *SQLite) DeleteByName(name string) error {
	_, err := s.DB.Exec(`DELETE FROM ips WHERE name = ?`, name)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to delete record: %s", err.Error()))
		return err
	}
	return nil
}

func (s *SQLite) GetById(id string) (*shared.IPInfo, error) {
	var info shared.IPInfo
	err := s.DB.QueryRow(`SELECT * FROM ips WHERE id = ?`, id).Scan(&info.IP, &info.Name, &info.Authkey, &info.Changed)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get record: %s", err.Error()))
		return nil, err
	}
	return &info, nil
}

func (s *SQLite) GetByName(name string) (*shared.IPInfo, error) {
	var info shared.IPInfo
	err := s.DB.QueryRow(`SELECT * FROM ips WHERE name = ?`, name).Scan(&info.IP, &info.Name, &info.Authkey, &info.Changed)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get record: %s", err.Error()))
		return nil, err
	}
	return &info, nil
}

func (s *SQLite) GetAll() (*[]shared.IPInfo, error) {
	rows, err := s.DB.Query(`SELECT * FROM ips`)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get records: %s", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var infos []shared.IPInfo
	for rows.Next() {
		var info shared.IPInfo
		err := rows.Scan(&info.IP, &info.Name, &info.Authkey, &info.Changed)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to scan record: %s", err.Error()))
			return nil, err
		}
		infos = append(infos, info)
	}
	return &infos, nil
}
