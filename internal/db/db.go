package db

type Persistence interface {
	// State

	// Init persistence connection
	Init() error
	// Close terminates persistence connection
	Close() error

	// Users

	// GetUsers returns a slice of all users
	GetUsers() ([]User, error)
	// GetUser fetches a user by their telegram id
	GetUser(telegramId int64) (User, error)
	// IsKnownUser returns bool inndicating if user is known to system
	IsKnownUser(telegramId int64) bool
	// AddUser persists a user in system
	AddUser(user User) error
	// UpdateUser updates user information in system
	UpdateUser(user User) error
	// DeleteUser deletes user information from system
	DeleteUser(user User) error

	// Notes

	// GetUserNotes returns a slice of all user's notes
	GetUserNotes(user User) ([]Note, error)
	// AddNote adds a new note to user's notes
	AddNote(user User, note Note) error
	// UpdateNote updates note content
	UpdateNote(note Note) error
	// DeleteNote deletes a note
	DeleteNote(note Note) error
}
