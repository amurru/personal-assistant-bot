package db

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	_ "github.com/dotenv-org/godotenvvault/autoload"
	"github.com/supabase-community/supabase-go"
)

var (
	instance Persistence
	apiUrl   = os.Getenv("SUPABASE_URL")
	apiKey   = os.Getenv("SUPABASE_KEY")
)

type Supabase struct {
	client *supabase.Client
}

func InstanceOrNew() Persistence {
	if instance != nil {
		return instance
	}

	if apiUrl == "" || apiKey == "" {
		log.Fatalf("SUPABASE_URL or SUPABASE_KEY is not set")
	}
	client, err := supabase.NewClient(apiUrl, apiKey, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Error init SupaBase: %v", err)
	}
	instance = &Supabase{
		client: client,
	}

	return instance
}

// methods

func (s *Supabase) Init() error {
	return nil
}
func (s *Supabase) Close() error {
	return nil
}

func (s *Supabase) GetUsers() ([]User, error) {
	return []User{}, nil
}

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

func (s *Supabase) UpdateUser(user User) error {
	return nil
}

func (s *Supabase) DeleteUser(user User) error {
	return nil
}

func (s *Supabase) GetUserNotes(user User) ([]Note, error) {
	return []Note{}, nil
}

func (s *Supabase) AddNote(user User, note Note) error {
	return nil
}

func (s *Supabase) UpdateNote(note Note) error {
	return nil
}
func (s *Supabase) DeleteNote(note Note) error {
	return nil
}
