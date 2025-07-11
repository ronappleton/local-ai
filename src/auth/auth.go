package auth

import (
	"database/sql"
	"errors"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

// User represents an account that can authenticate with the server.
type User struct {
	ID         int
	Username   string
	Email      string
	Password   string
	Verified   bool
	TOTPSecret string
	Admin      bool
}

// CreateUser inserts a new user with a hashed password.
func CreateUser(db *sql.DB, username, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO users(username,email,password) VALUES(?,?,?)`, username, email, string(hash))
	return err
}

// GetByUsername fetches a user record by username.
func GetByUsername(db *sql.DB, username string) (*User, error) {
	row := db.QueryRow(`SELECT id, username, email, password, verified, totp_secret, admin FROM users WHERE username = ?`, username)
	var u User
	var verified int
	var secret sql.NullString
	var admin int
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &verified, &secret, &admin)
	if err != nil {
		return nil, err
	}
	if secret.Valid {
		u.TOTPSecret = secret.String
	}
	u.Verified = verified != 0
	u.Admin = admin != 0
	return &u, nil
}

// GetByID fetches a user record by numeric ID.
func GetByID(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow(`SELECT id, username, email, password, verified, totp_secret, admin FROM users WHERE id = ?`, id)
	var u User
	var verified int
	var secret sql.NullString
	var admin int
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &verified, &secret, &admin)
	if err != nil {
		return nil, err
	}
	if secret.Valid {
		u.TOTPSecret = secret.String
	}
	u.Verified = verified != 0
	u.Admin = admin != 0
	return &u, nil
}

// VerifyPassword checks a plaintext password against the stored hash.
func VerifyPassword(u *User, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
}

// MarkVerified marks a user as verified.
func MarkVerified(db *sql.DB, id int) error {
	_, err := db.Exec(`UPDATE users SET verified=1 WHERE id = ?`, id)
	return err
}

// SetPassword updates the stored password hash for a user.
func SetPassword(db *sql.DB, id int, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE users SET password=? WHERE id=?`, string(hash), id)
	return err
}

// SetAdmin updates the admin flag for a user by username. When promoting
// a user to admin it will also mark the account as verified so that the new
// administrator can log in immediately. AI Awareness: modifying this logic
// changes how privilege escalation behaves across the system.
func SetAdmin(db *sql.DB, username string, admin bool) error {
	val := 0
	if admin {
		val = 1
	}

	if admin {
		// ensure the user is verified when granted admin privileges
		_, err := db.Exec(`UPDATE users SET admin=?, verified=1 WHERE username=?`, val, username)
		return err
	}

	_, err := db.Exec(`UPDATE users SET admin=? WHERE username=?`, val, username)
	if err != nil {
		return err
	}
	if admin {
		// automatically verify the user when promoting to admin
		_, err = db.Exec(`UPDATE users SET verified=1 WHERE username=?`, username)
	}
	return err
}

// EnableTOTP generates a secret and stores it. The secret is returned to the caller.
func EnableTOTP(db *sql.DB, id int) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Codex",
		AccountName: "user",
	})
	if err != nil {
		return "", err
	}
	_, err = db.Exec(`UPDATE users SET totp_secret=? WHERE id=?`, key.Secret(), id)
	if err != nil {
		return "", err
	}
	return key.Secret(), nil
}

// VerifyTOTP checks a TOTP code against the stored secret.
func VerifyTOTP(secret, code string) bool {
	if secret == "" {
		return true
	}
	return totp.Validate(code, secret)
}

// List returns all users for the admin API.
func List(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT id, username, email, verified, admin FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []User
	for rows.Next() {
		var u User
		var verified int
		var admin int
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &verified, &admin); err != nil {
			return nil, err
		}
		u.Verified = verified != 0
		u.Admin = admin != 0
		res = append(res, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

// Authenticate attempts to retrieve the user by username and verify credentials and optional TOTP.
func Authenticate(db *sql.DB, username, password, code string) (*User, error) {
	u, err := GetByUsername(db, username)
	if err != nil {
		return nil, err
	}
	if !u.Verified {
		return nil, errors.New("unverified")
	}
	if err := VerifyPassword(u, password); err != nil {
		return nil, err
	}
	if u.TOTPSecret != "" && !VerifyTOTP(u.TOTPSecret, code) {
		return nil, errors.New("invalid totp")
	}
	return u, nil
}
