package ai

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/user/ghost-browser/internal/database"
)

type Personality struct {
	ID             string       `json:"id"`
	ProfileID      string       `json:"profileId"`
	Name           string       `json:"name"`
	Age            int          `json:"age"`
	Gender         string       `json:"gender"`
	Occupation     string       `json:"occupation"`
	Location       string       `json:"location"`
	Bio            string       `json:"bio"`
	Interests      []string     `json:"interests"`
	ExpertiseAreas []string     `json:"expertiseAreas"`
	WritingStyle   WritingStyle `json:"writingStyle"`
	TypingSpeed    TypingProfile `json:"typingSpeed"`
	MouseBehavior  MouseProfile `json:"mouseBehavior"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

type WritingStyle struct {
	Formality     string   `json:"formality"`
	Tone          string   `json:"tone"`
	Verbosity     string   `json:"verbosity"`
	UseEmojis     bool     `json:"useEmojis"`
	UseSlang      bool     `json:"useSlang"`
	CommonPhrases []string `json:"commonPhrases"`
}

type TypingProfile struct {
	AverageWPM     int     `json:"averageWpm"`
	Variance       float64 `json:"variance"`
	ErrorRate      float64 `json:"errorRate"`
	PauseBetween   int     `json:"pauseBetween"`
	ThinkingPauses bool    `json:"thinkingPauses"`
}

type MouseProfile struct {
	MovementSpeed  string  `json:"movementSpeed"`
	Precision      string  `json:"precision"`
	ScrollBehavior string  `json:"scrollBehavior"`
	ClickDelay     int     `json:"clickDelay"`
	Jitter         float64 `json:"jitter"`
}

type Schedule struct {
	ID             string        `json:"id"`
	ProfileID      string        `json:"profileId"`
	Timezone       string        `json:"timezone"`
	ActiveHours    ActiveHours   `json:"activeHours"`
	WeeklySchedule []DaySchedule `json:"weeklySchedule"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
}

type ActiveHours struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type DaySchedule struct {
	DayOfWeek  int                 `json:"dayOfWeek"`
	Activities []ScheduledActivity `json:"activities"`
}

type ScheduledActivity struct {
	TimeStart   string   `json:"timeStart"`
	TimeEnd     string   `json:"timeEnd"`
	Activity    string   `json:"activity"`
	Sites       []string `json:"sites"`
	Description string   `json:"description"`
}

type Conversation struct {
	ID        string    `json:"id"`
	ProfileID string    `json:"profileId"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type Engine struct {
	db        *database.Database
	ollamaURL string
	model     string
}

func NewEngine(db *database.Database) *Engine {
	return &Engine{
		db:        db,
		ollamaURL: "http://localhost:11434",
		model:     "llama3.2",
	}
}

func (e *Engine) SetModel(model string)      { e.model = model }
func (e *Engine) SetOllamaURL(url string)    { e.ollamaURL = url }

func (e *Engine) GetPersonality(profileID string) (*Personality, error) {
	row := e.db.DB().QueryRow(`
		SELECT id, profile_id, name, age, gender, occupation, location, bio,
		       interests, expertise_areas, writing_style, typing_speed, mouse_behavior,
		       created_at, updated_at
		FROM personalities WHERE profile_id = ?`, profileID)

	var p Personality
	var interests, expertise, writing, typing, mouse string

	err := row.Scan(&p.ID, &p.ProfileID, &p.Name, &p.Age, &p.Gender, &p.Occupation,
		&p.Location, &p.Bio, &interests, &expertise, &writing, &typing, &mouse,
		&p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	json.Unmarshal([]byte(interests), &p.Interests)
	json.Unmarshal([]byte(expertise), &p.ExpertiseAreas)
	json.Unmarshal([]byte(writing), &p.WritingStyle)
	json.Unmarshal([]byte(typing), &p.TypingSpeed)
	json.Unmarshal([]byte(mouse), &p.MouseBehavior)

	return &p, nil
}

func (e *Engine) UpdatePersonality(profileID string, p *Personality) error {
	existing, _ := e.GetPersonality(profileID)

	interests, _ := json.Marshal(p.Interests)
	expertise, _ := json.Marshal(p.ExpertiseAreas)
	writing, _ := json.Marshal(p.WritingStyle)
	typing, _ := json.Marshal(p.TypingSpeed)
	mouse, _ := json.Marshal(p.MouseBehavior)
	now := time.Now()

	if existing == nil {
		p.ID = uuid.New().String()
		p.ProfileID = profileID
		p.CreatedAt = now
		p.UpdatedAt = now

		_, err := e.db.DB().Exec(`
			INSERT INTO personalities (id, profile_id, name, age, gender, occupation, location, bio,
			                          interests, expertise_areas, writing_style, typing_speed, mouse_behavior,
			                          created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.ID, p.ProfileID, p.Name, p.Age, p.Gender, p.Occupation, p.Location, p.Bio,
			string(interests), string(expertise), string(writing), string(typing), string(mouse),
			p.CreatedAt, p.UpdatedAt)
		return err
	}

	p.UpdatedAt = now
	_, err := e.db.DB().Exec(`
		UPDATE personalities SET name=?, age=?, gender=?, occupation=?, location=?, bio=?,
		       interests=?, expertise_areas=?, writing_style=?, typing_speed=?, mouse_behavior=?, updated_at=?
		WHERE profile_id=?`,
		p.Name, p.Age, p.Gender, p.Occupation, p.Location, p.Bio,
		string(interests), string(expertise), string(writing), string(typing), string(mouse),
		p.UpdatedAt, profileID)
	return err
}

func (e *Engine) GenerateRandom() (*Personality, error) {
	return &Personality{
		Name:       "Random User",
		Age:        25,
		Gender:     "other",
		Occupation: "Software Developer",
		Location:   "New York, USA",
		Bio:        "A tech enthusiast who loves exploring new technologies.",
		Interests:  []string{"technology", "gaming", "music"},
		ExpertiseAreas: []string{"programming", "web development"},
		WritingStyle: WritingStyle{
			Formality: "casual",
			Tone:      "friendly",
			Verbosity: "concise",
			UseEmojis: true,
			UseSlang:  false,
		},
		TypingSpeed: TypingProfile{
			AverageWPM:     60,
			Variance:       0.2,
			ErrorRate:      0.02,
			PauseBetween:   100,
			ThinkingPauses: true,
		},
		MouseBehavior: MouseProfile{
			MovementSpeed:  "medium",
			Precision:      "medium",
			ScrollBehavior: "smooth",
			ClickDelay:     50,
			Jitter:         0.1,
		},
	}, nil
}

func (e *Engine) Chat(profileID, message string) (string, error) {
	personality, err := e.GetPersonality(profileID)
	if err != nil {
		return "", err
	}

	systemPrompt := "You are a helpful assistant."
	if personality != nil {
		systemPrompt = fmt.Sprintf(`You are %s, a %d year old %s from %s.
Your occupation is %s. Bio: %s
Your interests include: %v
Your writing style is %s, %s, and %s.
%s`,
			personality.Name, personality.Age, personality.Gender, personality.Location,
			personality.Occupation, personality.Bio, personality.Interests,
			personality.WritingStyle.Formality, personality.WritingStyle.Tone, personality.WritingStyle.Verbosity,
			func() string {
				if personality.WritingStyle.UseEmojis {
					return "You often use emojis."
				}
				return ""
			}())
	}

	// Get conversation history
	history := e.getConversationHistory(profileID, 10)

	// Build messages
	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
	}
	for _, h := range history {
		messages = append(messages, map[string]string{"role": h.Role, "content": h.Content})
	}
	messages = append(messages, map[string]string{"role": "user", "content": message})

	// Call Ollama
	response, err := e.callOllama(messages)
	if err != nil {
		return "", err
	}

	// Save conversation
	e.saveConversation(profileID, "user", message)
	e.saveConversation(profileID, "assistant", response)

	return response, nil
}

func (e *Engine) callOllama(messages []map[string]string) (string, error) {
	payload := map[string]interface{}{
		"model":    e.model,
		"messages": messages,
		"stream":   false,
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(e.ollamaURL+"/api/chat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}

func (e *Engine) getConversationHistory(profileID string, limit int) []Conversation {
	rows, err := e.db.DB().Query(`
		SELECT id, profile_id, role, content, created_at
		FROM conversations WHERE profile_id = ?
		ORDER BY created_at DESC LIMIT ?`, profileID, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var history []Conversation
	for rows.Next() {
		var c Conversation
		rows.Scan(&c.ID, &c.ProfileID, &c.Role, &c.Content, &c.CreatedAt)
		history = append([]Conversation{c}, history...)
	}
	return history
}

func (e *Engine) saveConversation(profileID, role, content string) {
	e.db.DB().Exec(`
		INSERT INTO conversations (id, profile_id, role, content, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		uuid.New().String(), profileID, role, content, time.Now())
}

func (e *Engine) GetSchedule(profileID string) (*Schedule, error) {
	row := e.db.DB().QueryRow(`
		SELECT id, profile_id, timezone, active_hours, weekly_schedule, created_at, updated_at
		FROM schedules WHERE profile_id = ?`, profileID)

	var s Schedule
	var activeHours, weekly string

	err := row.Scan(&s.ID, &s.ProfileID, &s.Timezone, &activeHours, &weekly, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	json.Unmarshal([]byte(activeHours), &s.ActiveHours)
	json.Unmarshal([]byte(weekly), &s.WeeklySchedule)
	return &s, nil
}

func (e *Engine) UpdateSchedule(profileID string, s *Schedule) error {
	existing, _ := e.GetSchedule(profileID)

	activeHours, _ := json.Marshal(s.ActiveHours)
	weekly, _ := json.Marshal(s.WeeklySchedule)
	now := time.Now()

	if existing == nil {
		s.ID = uuid.New().String()
		s.ProfileID = profileID
		s.CreatedAt = now
		s.UpdatedAt = now

		_, err := e.db.DB().Exec(`
			INSERT INTO schedules (id, profile_id, timezone, active_hours, weekly_schedule, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			s.ID, s.ProfileID, s.Timezone, string(activeHours), string(weekly), s.CreatedAt, s.UpdatedAt)
		return err
	}

	s.UpdatedAt = now
	_, err := e.db.DB().Exec(`
		UPDATE schedules SET timezone=?, active_hours=?, weekly_schedule=?, updated_at=? WHERE profile_id=?`,
		s.Timezone, string(activeHours), string(weekly), s.UpdatedAt, profileID)
	return err
}
