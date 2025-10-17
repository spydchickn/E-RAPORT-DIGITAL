package models

import (
	"database/sql"
	"e_raport_digital/config"
)

type Siswa struct {
    IDSiswa        int
    NIS            string
    Nama           string
    Alamat         string
    Foto           string
    IDKelas        sql.NullInt64
    IDUser         sql.NullInt64
    OrtuNama       string
    OrtuPekerjaan  string
    OrtuAlamat     string
    OrtuTelepon    string
    OTP            string
}

func GetAllSiswa() ([]Siswa, error) {
    rows, err := config.DB.Query("SELECT id_siswa, nis, nama, alamat, id_kelas, id_user FROM siswa")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var list []Siswa
    for rows.Next() {
        var s Siswa
        rows.Scan(&s.IDSiswa, &s.NIS, &s.Nama, &s.Alamat, &s.IDKelas, &s.IDUser)
        list = append(list, s)
    }
    return list, nil
}

// GetSiswaByUserID returns the siswa record linked to a users.id (id_user)
func GetSiswaByUserID(userID int) (*Siswa, error) {
    row := config.DB.QueryRow("SELECT id_siswa, nis, nama, alamat, id_kelas, id_user FROM siswa WHERE id_user = ?", userID)
    var s Siswa
    err := row.Scan(&s.IDSiswa, &s.NIS, &s.Nama, &s.Alamat, &s.IDKelas, &s.IDUser)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &s, nil
}

// GetSiswaByID returns the siswa record by id_siswa
func GetSiswaByID(id int) (*Siswa, error) {
    row := config.DB.QueryRow("SELECT id_siswa, nis, nama, alamat, id_kelas, id_user FROM siswa WHERE id_siswa = ?", id)
    var s Siswa
    err := row.Scan(&s.IDSiswa, &s.NIS, &s.Nama, &s.Alamat, &s.IDKelas, &s.IDUser)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &s, nil
}

func CreateSiswa(nis, nama, alamat, foto, ortuNama, ortuPekerjaan, ortuAlamat, ortuTelepon, otp string, idUser int64) (int64, error) {
    if idUser > 0 {
        res, err := config.DB.Exec("INSERT INTO siswa (nis, nama, alamat, foto, ortu_nama, ortu_pekerjaan, ortu_alamat, ortu_telepon, otp, id_user) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", nis, nama, alamat, foto, ortuNama, ortuPekerjaan, ortuAlamat, ortuTelepon, otp, idUser)
        if err != nil {
            return 0, err
        }
        return res.LastInsertId()
    }
    res, err := config.DB.Exec("INSERT INTO siswa (nis, nama, alamat, foto, ortu_nama, ortu_pekerjaan, ortu_alamat, ortu_telepon, otp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", nis, nama, alamat, foto, ortuNama, ortuPekerjaan, ortuAlamat, ortuTelepon, otp)
    if err != nil {
        return 0, err
    }
    id, _ := res.LastInsertId()
    return id, nil
}

func UpdateSiswa(id int, nis, nama, alamat string) error {
    _, err := config.DB.Exec("UPDATE siswa SET nis=?, nama=?, alamat=? WHERE id_siswa=?", nis, nama, alamat, id)
    return err
}

func DeleteSiswa(id int) error {
    _, err := config.DB.Exec("DELETE FROM siswa WHERE id_siswa=?", id)
    return err
}
