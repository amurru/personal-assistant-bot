package db

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

var (
	instance Persistence
	apiURL   = os.Getenv("SUPABASE_URL")
	apiKey   = os.Getenv("SUPABASE_KEY")
)

// Supabase is a persistence layer for the bot
type Supabase struct {
	client *supabase.Client
}

// InstanceOrNew returns a singleton instance of the Supabase persistence layer
func InstanceOrNew() Persistence {
	if instance != nil {
		return instance
	}

	if apiURL == "" || apiKey == "" {
		log.Fatalf("SUPABASE_URL or SUPABASE_KEY is not set")
	}
	client, err := supabase.NewClient(apiURL, apiKey, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Error init SupaBase: %v", err)
	}
	instance = &Supabase{
		client: client,
	}

	return instance
}

// methods

// Init persistence connection
func (s *Supabase) Init() error {
	return nil
}

// Close terminates persistence connection
func (s *Supabase) Close() error {
	return nil
}

// GetUsers returns a slice of all users
func (s *Supabase) GetUsers() ([]User, error) {
	return []User{}, nil
}

// GetUser fetches a user by their telegram id
func (s *Supabase) GetUser(telegramID int64) (User, error) {
	id := strconv.FormatInt(telegramID, 10)
	result, _, err := s.client.From("users").
		Select("*", "exact", false).
		Single().
		Eq("id", id).
		Execute()
	if err != nil {
		log.Printf("GetUser error: %v", err)
		return User{}, err
	}

	var user User
	err = json.Unmarshal(result, &user)
	if err != nil {
		log.Printf("Error unmarshalling user: %v", err)
		return User{}, err
	}
	return user, nil
}

// IsKnownUser returns bool inndicating if user is known to system
func (s *Supabase) IsKnownUser(telegramID int64) bool {
	id := strconv.FormatInt(telegramID, 10)
	_, count, err := s.client.From("users").
		Select("*", "exact", false).
		Eq("id", id).
		Execute()
	if err != nil {
		// log error and return false
		log.Printf("IsKnownUser error: %v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true
}

// AddUser persists a user in system
func (s *Supabase) AddUser(user User) error {
	_, count, err := s.client.From("users").
		Insert(user, false, "", "", "exact").
		Execute()
	if err != nil {
		log.Printf("AddUser error: %v", err)
		return err
	}
	log.Printf("Added %d user to system", count)
	return nil
}

// UpdateUser updates user information in system
func (s *Supabase) UpdateUser(user User) error {
	return nil
}

// DeleteUser deletes user information from system
func (s *Supabase) DeleteUser(user User) error {
	return nil
}

// GetUserNotes returns a slice of all user's notes
func (s *Supabase) GetUserNotes(userID int64) ([]Note, error) {
	result, _, err := s.client.From("user_notes").
		Select("*", "exact", false).
		Eq("user_id", strconv.FormatInt(userID, 10)).
		Execute()
	if err != nil {
		log.Printf("GetUserNotes error: %v", err)
		return []Note{}, err
	}
	var notes []Note
	err = json.Unmarshal(result, &notes)
	if err != nil {
		log.Printf("Error unmarshalling notes: %v", err)
		return []Note{}, err
	}
	return notes, nil
}

// AddNote adds a new note to user's notes
func (s *Supabase) AddNote(note Note) error {
	_, count, err := s.client.From("user_notes").
		Insert(note, false, "", "", "exact").
		Execute()
	if err != nil {
		log.Printf("AddNote error: %v", err)
		return err
	}
	log.Printf("Added %d note to system", count)
	return nil
}

// UpdateNote updates note content
func (s *Supabase) UpdateNote(note Note) error {
	return nil
}

// DeleteNote deletes a note
func (s *Supabase) DeleteNote(note Note) error {
	return nil
}
