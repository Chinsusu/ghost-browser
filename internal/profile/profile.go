package profile

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/user/ghost-browser/internal/database"
	"github.com/user/ghost-browser/internal/fingerprint"
)

type Profile struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Fingerprint *fingerprint.Fingerprint `json:"fingerprint"`
	ProxyID     *string                  `json:"proxyId,omitempty"`
	DataDir     string                   `json:"dataDir"`
	Notes       string                   `json:"notes"`
	Tags        []string                 `json:"tags"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
	LastUsedAt  *time.Time               `json:"lastUsedAt,omitempty"`
}

type CreateOptions struct {
	Fingerprint *fingerprint.Fingerprint `json:"fingerprint,omitempty"`
	ProxyID     *string                  `json:"proxyId,omitempty"`
	Notes       string                   `json:"notes,omitempty"`
	Tags        []string                 `json:"tags,omitempty"`
	OS          string                   `json:"os,omitempty"`
	Browser     string                   `json:"browser,omitempty"`
}

type Manager struct {
	db *database.Database
}

func NewManager(db *database.Database) *Manager {
	return &Manager{db: db}
}

func (m *Manager) GetAll() ([]*Profile, error) {
	rows, err := m.db.DB().Query(`
		SELECT id, name, fingerprint, proxy_id, data_dir, notes, tags,
		       created_at, updated_at, last_used_at FROM profiles ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*Profile
	for rows.Next() {
		p, err := scanProfile(rows)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (m *Manager) GetByID(id string) (*Profile, error) {
	row := m.db.DB().QueryRow(`
		SELECT id, name, fingerprint, proxy_id, data_dir, notes, tags,
		       created_at, updated_at, last_used_at FROM profiles WHERE id = ?`, id)
	return scanProfileRow(row)
}

func (m *Manager) Create(name string, opts *CreateOptions) (*Profile, error) {
	if opts == nil {
		opts = &CreateOptions{}
	}

	id := uuid.New().String()
	fp := opts.Fingerprint
	if fp == nil {
		gen := fingerprint.NewGenerator()
		fp = gen.Generate(&fingerprint.GenerateOptions{OS: opts.OS, Browser: opts.Browser})
	}

	dataDir, err := m.createDataDir(id)
	if err != nil {
		return nil, err
	}

	fpJSON, _ := json.Marshal(fp)
	tags := opts.Tags
	if tags == nil {
		tags = []string{}
	}
	tagsJSON, _ := json.Marshal(tags)
	now := time.Now()

	_, err = m.db.DB().Exec(`
		INSERT INTO profiles (id, name, fingerprint, proxy_id, data_dir, notes, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, name, string(fpJSON), opts.ProxyID, dataDir, opts.Notes, string(tagsJSON), now, now)
	if err != nil {
		return nil, err
	}

	return &Profile{
		ID: id, Name: name, Fingerprint: fp, ProxyID: opts.ProxyID,
		DataDir: dataDir, Notes: opts.Notes, Tags: tags, CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (m *Manager) Update(p *Profile) error {
	fpJSON, _ := json.Marshal(p.Fingerprint)
	tagsJSON, _ := json.Marshal(p.Tags)
	p.UpdatedAt = time.Now()

	_, err := m.db.DB().Exec(`
		UPDATE profiles SET name=?, fingerprint=?, proxy_id=?, notes=?, tags=?, updated_at=? WHERE id=?`,
		p.Name, string(fpJSON), p.ProxyID, p.Notes, string(tagsJSON), p.UpdatedAt, p.ID)
	return err
}

func (m *Manager) Delete(id string) error {
	p, err := m.GetByID(id)
	if err != nil {
		return err
	}
	_, err = m.db.DB().Exec("DELETE FROM profiles WHERE id = ?", id)
	if err != nil {
		return err
	}
	if p.DataDir != "" {
		os.RemoveAll(p.DataDir)
	}
	return nil
}

func (m *Manager) GenerateRandom() (*Profile, error) {
	gen := fingerprint.NewGenerator()
	return m.Create(gen.GenerateRandomName(), &CreateOptions{OS: "windows", Browser: "edge"})
}

func (m *Manager) Duplicate(id string) (*Profile, error) {
	orig, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}
	return m.Create(orig.Name+" (Copy)", &CreateOptions{
		Fingerprint: orig.Fingerprint, ProxyID: orig.ProxyID, Notes: orig.Notes, Tags: orig.Tags,
	})
}

func (m *Manager) Export(id, path string) error {
	p, err := m.GetByID(id)
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(p, "", "  ")
	return os.WriteFile(path, data, 0644)
}

func (m *Manager) Import(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return m.Create(p.Name+" (Imported)", &CreateOptions{
		Fingerprint: p.Fingerprint, Notes: p.Notes, Tags: p.Tags,
	})
}

func (m *Manager) UpdateLastUsed(id string) error {
	now := time.Now()
	_, err := m.db.DB().Exec("UPDATE profiles SET last_used_at=?, updated_at=? WHERE id=?", now, now, id)
	return err
}

func (m *Manager) createDataDir(id string) (string, error) {
	dataDir, err := database.GetDataDir()
	if err != nil {
		return "", err
	}
	profileDir := filepath.Join(dataDir, "profiles", id)
	return profileDir, os.MkdirAll(profileDir, 0755)
}

func scanProfile(rows *sql.Rows) (*Profile, error) {
	var p Profile
	var fpJSON, tagsJSON string
	var proxyID sql.NullString
	var lastUsed sql.NullTime

	err := rows.Scan(&p.ID, &p.Name, &fpJSON, &proxyID, &p.DataDir,
		&p.Notes, &tagsJSON, &p.CreatedAt, &p.UpdatedAt, &lastUsed)
	if err != nil {
		return nil, err
	}

	if proxyID.Valid {
		p.ProxyID = &proxyID.String
	}
	if lastUsed.Valid {
		p.LastUsedAt = &lastUsed.Time
	}
	json.Unmarshal([]byte(fpJSON), &p.Fingerprint)
	json.Unmarshal([]byte(tagsJSON), &p.Tags)
	return &p, nil
}

func scanProfileRow(row *sql.Row) (*Profile, error) {
	var p Profile
	var fpJSON, tagsJSON string
	var proxyID sql.NullString
	var lastUsed sql.NullTime

	err := row.Scan(&p.ID, &p.Name, &fpJSON, &proxyID, &p.DataDir,
		&p.Notes, &tagsJSON, &p.CreatedAt, &p.UpdatedAt, &lastUsed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, err
	}

	if proxyID.Valid {
		p.ProxyID = &proxyID.String
	}
	if lastUsed.Valid {
		p.LastUsedAt = &lastUsed.Time
	}
	json.Unmarshal([]byte(fpJSON), &p.Fingerprint)
	json.Unmarshal([]byte(tagsJSON), &p.Tags)
	return &p, nil
}
