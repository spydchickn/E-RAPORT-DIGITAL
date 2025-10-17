package models

import (
	"e_raport_digital/config"

	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers() ([]User, error) {
    rows, err := config.DB.Query("SELECT id, username, password, role FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var list []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role)
        list = append(list, u)
    }
    return list, nil
}

// CreateUser hashes the password and returns the inserted user ID
func CreateUser(username, password, role string) (int64, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return 0, err
    }
    res, err := config.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", username, string(hashed), role)
    if err != nil {
        return 0, err
    }
    id, _ := res.LastInsertId()
    return id, nil
}

// UpdateUser updates user fields; password will be re-hashed
func UpdateUser(id int, username, password, role string) error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    _, err = config.DB.Exec("UPDATE users SET username=?, password=?, role=? WHERE id=?", username, string(hashed), role, id)
    return err
}

func DeleteUser(id int) error {
    _, err := config.DB.Exec("DELETE FROM users WHERE id=?", id)
    return err
}

// SetUserPasswordHash updates the stored password hash for a user id
func SetUserPasswordHash(id int, hash string) error {
    _, err := config.DB.Exec("UPDATE users SET password=? WHERE id=?", hash, id)
    return err
}

// CountUsers returns total number of users
func CountUsers() (int, error) {
    var cnt int
    err := config.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&cnt)
    return cnt, err
}
