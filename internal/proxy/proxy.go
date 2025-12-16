package proxy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/user/ghost-browser/internal/database"
)

type ProxyType string

const (
	HTTP   ProxyType = "http"
	HTTPS  ProxyType = "https"
	SOCKS4 ProxyType = "socks4"
	SOCKS5 ProxyType = "socks5"
)

type Proxy struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	Type             ProxyType  `json:"type"`
	Host             string     `json:"host"`
	Port             int        `json:"port"`
	Username         string     `json:"username,omitempty"`
	Password         string     `json:"password,omitempty"`
	Country          string     `json:"country,omitempty"`
	LastCheckAt      *time.Time `json:"lastCheckAt,omitempty"`
	LastCheckStatus  string     `json:"lastCheckStatus"`
	LastCheckLatency int        `json:"lastCheckLatency"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

func (p *Proxy) ToURL() string {
	var auth string
	if p.Username != "" {
		auth = p.Username
		if p.Password != "" {
			auth += ":" + p.Password
		}
		auth += "@"
	}
	return fmt.Sprintf("%s://%s%s:%d", p.Type, auth, p.Host, p.Port)
}

type CheckResult struct {
	ProxyID string `json:"proxyId"`
	Success bool   `json:"success"`
	Latency int    `json:"latency"`
	IP      string `json:"ip,omitempty"`
	Country string `json:"country,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Manager struct {
	db *database.Database
}

func NewManager(db *database.Database) *Manager {
	return &Manager{db: db}
}

func (m *Manager) GetAll() ([]*Proxy, error) {
	rows, err := m.db.DB().Query(`
		SELECT id, name, type, host, port, username, password, country,
		       last_check_at, last_check_status, last_check_latency, created_at, updated_at
		FROM proxies ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proxies []*Proxy
	for rows.Next() {
		p, err := scanProxy(rows)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, p)
	}
	return proxies, nil
}

func (m *Manager) GetByID(id string) (*Proxy, error) {
	row := m.db.DB().QueryRow(`
		SELECT id, name, type, host, port, username, password, country,
		       last_check_at, last_check_status, last_check_latency, created_at, updated_at
		FROM proxies WHERE id = ?`, id)
	return scanProxyRow(row)
}

func (m *Manager) Add(p *Proxy) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.Name == "" {
		p.Name = fmt.Sprintf("%s:%d", p.Host, p.Port)
	}
	now := time.Now()
	p.CreatedAt, p.UpdatedAt = now, now
	p.LastCheckStatus = "unknown"

	_, err := m.db.DB().Exec(`
		INSERT INTO proxies (id, name, type, host, port, username, password, country, last_check_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Type, p.Host, p.Port, p.Username, p.Password, p.Country, p.LastCheckStatus, p.CreatedAt, p.UpdatedAt)
	return err
}

func (m *Manager) Update(p *Proxy) error {
	p.UpdatedAt = time.Now()
	_, err := m.db.DB().Exec(`
		UPDATE proxies SET name=?, type=?, host=?, port=?, username=?, password=?, country=?, updated_at=? WHERE id=?`,
		p.Name, p.Type, p.Host, p.Port, p.Username, p.Password, p.Country, p.UpdatedAt, p.ID)
	return err
}

func (m *Manager) Delete(id string) error {
	_, err := m.db.DB().Exec("DELETE FROM proxies WHERE id = ?", id)
	return err
}

func (m *Manager) Check(id string) (*CheckResult, error) {
	p, err := m.GetByID(id)
	if err != nil {
		return nil, err
	}

	result := m.checkProxy(p)
	now := time.Now()
	status := "working"
	if !result.Success {
		status = "failed"
	}

	m.db.DB().Exec(`
		UPDATE proxies SET last_check_at=?, last_check_status=?, last_check_latency=?, country=COALESCE(NULLIF(?, ''), country), updated_at=? WHERE id=?`,
		now, status, result.Latency, result.Country, now, id)

	return result, nil
}

func (m *Manager) CheckAll() ([]*CheckResult, error) {
	proxies, err := m.GetAll()
	if err != nil {
		return nil, err
	}
	results := make([]*CheckResult, 0, len(proxies))
	for _, p := range proxies {
		r, _ := m.Check(p.ID)
		results = append(results, r)
	}
	return results, nil
}

func (m *Manager) checkProxy(p *Proxy) *CheckResult {
	result := &CheckResult{ProxyID: p.ID}

	proxyURL, err := url.Parse(p.ToURL())
	if err != nil {
		result.Error = err.Error()
		return result
	}

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		Timeout:   10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get("https://api.ipify.org?format=json")
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.Latency = int(time.Since(start).Milliseconds())
	result.Success = true

	var ipResp struct{ IP string `json:"ip"` }
	json.NewDecoder(resp.Body).Decode(&ipResp)
	result.IP = ipResp.IP

	return result
}

func (m *Manager) ImportFromText(text, defaultType string) (int, error) {
	if defaultType == "" {
		defaultType = "http"
	}

	count := 0
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		proxy, err := parseProxyLine(line, defaultType)
		if err != nil {
			continue
		}
		if m.Add(proxy) == nil {
			count++
		}
	}
	return count, nil
}

func parseProxyLine(line, defaultType string) (*Proxy, error) {
	proxy := &Proxy{Type: ProxyType(defaultType)}

	if strings.Contains(line, "://") {
		u, err := url.Parse(line)
		if err != nil {
			return nil, err
		}
		proxy.Type = ProxyType(u.Scheme)
		proxy.Host = u.Hostname()
		fmt.Sscanf(u.Port(), "%d", &proxy.Port)
		if u.User != nil {
			proxy.Username = u.User.Username()
			proxy.Password, _ = u.User.Password()
		}
		return proxy, nil
	}

	parts := strings.Split(line, ":")
	if len(parts) >= 2 {
		proxy.Host = parts[0]
		fmt.Sscanf(parts[1], "%d", &proxy.Port)
		if len(parts) >= 4 {
			proxy.Username = parts[2]
			proxy.Password = parts[3]
		}
	}

	if proxy.Host == "" || proxy.Port == 0 {
		return nil, fmt.Errorf("invalid format")
	}
	return proxy, nil
}

func scanProxy(rows *sql.Rows) (*Proxy, error) {
	var p Proxy
	var lastCheck sql.NullTime
	err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Host, &p.Port, &p.Username, &p.Password,
		&p.Country, &lastCheck, &p.LastCheckStatus, &p.LastCheckLatency, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if lastCheck.Valid {
		p.LastCheckAt = &lastCheck.Time
	}
	return &p, nil
}

func scanProxyRow(row *sql.Row) (*Proxy, error) {
	var p Proxy
	var lastCheck sql.NullTime
	err := row.Scan(&p.ID, &p.Name, &p.Type, &p.Host, &p.Port, &p.Username, &p.Password,
		&p.Country, &lastCheck, &p.LastCheckStatus, &p.LastCheckLatency, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("proxy not found")
		}
		return nil, err
	}
	if lastCheck.Valid {
		p.LastCheckAt = &lastCheck.Time
	}
	return &p, nil
}
