package models

import (
    "database/sql"
    "e_raport_digital/config"
)

type User struct {
    ID       int
    Username string
    Password string
    Role     string
}

func GetUserByUsername(username string) (*User, error) {
    var u User
    row := config.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", username)
    err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Role)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &u, nil
}
