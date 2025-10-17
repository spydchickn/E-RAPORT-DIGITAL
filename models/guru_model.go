package models

import (
	"database/sql"
	"e_raport_digital/config"
)

// GetAllGuru returns all guru records including optional linked user id
func GetAllGuru() ([]Guru, error) {
    rows, err := config.DB.Query("SELECT id_guru, nama, nip, alamat, id_user FROM guru")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var list []Guru
    for rows.Next() {
        var g Guru
        var idUser sql.NullInt64
        if err := rows.Scan(&g.IDGuru, &g.Nama, &g.NIP, &g.Alamat, &idUser); err != nil {
            return nil, err
        }
        g.IDUser = idUser
        list = append(list, g)
    }
    return list, nil
}

// CreateGuruWithUser creates a guru and links it to a user via id_user column if provided
func CreateGuruWithUser(nama, nip, alamat, foto string, idUser int64) (int64, error) {
    if idUser > 0 {
        res, err := config.DB.Exec("INSERT INTO guru (nama, nip, alamat, foto, id_user) VALUES (?, ?, ?, ?, ?)", nama, nip, alamat, foto, idUser)
        if err != nil {
            return 0, err
        }
        return res.LastInsertId()
    }
    res, err := config.DB.Exec("INSERT INTO guru (nama, nip, alamat, foto) VALUES (?, ?, ?, ?)", nama, nip, alamat, foto)
    if err != nil {
        return 0, err
    }
    id, _ := res.LastInsertId()
    return id, nil
}

// GetGuruByUserID returns a guru linked to a users.id via id_user column (if schema uses it)
func GetGuruByUserID(userID int64) (*Guru, error) {
    row := config.DB.QueryRow("SELECT id_guru, nama, nip, alamat, id_user FROM guru WHERE id_user = ?", userID)
    var g Guru
    var idUser sql.NullInt64
    err := row.Scan(&g.IDGuru, &g.Nama, &g.NIP, &g.Alamat, &idUser)
    g.IDUser = idUser
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &g, nil
}

// GetGuruByID returns a guru by its ID
func GetGuruByID(id int) (*Guru, error) {
    row := config.DB.QueryRow("SELECT id_guru, nama, nip, alamat, id_user FROM guru WHERE id_guru = ?", id)
    var g Guru
    var idUser sql.NullInt64
    err := row.Scan(&g.IDGuru, &g.Nama, &g.NIP, &g.Alamat, &idUser)
    g.IDUser = idUser
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &g, nil
}

// CreateGuru creates a guru without linking to a user
func CreateGuru(nama, nip, alamat, foto string) error {
    _, err := config.DB.Exec("INSERT INTO guru (nama, nip, alamat, foto) VALUES (?, ?, ?, ?)", nama, nip, alamat, foto)
    return err
}

// DeleteGuru deletes a guru by id
func DeleteGuru(id int) error {
    _, err := config.DB.Exec("DELETE FROM guru WHERE id_guru=?", id)
    return err
}