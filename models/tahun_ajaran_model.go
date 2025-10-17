package models

import (
    "database/sql"

    "e_raport_digital/config"
)

func GetAllTahunAjaran() ([]TahunAjaran, error) {
    query := `SELECT id_tahun, nama, aktif FROM tahun_ajaran`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tahun []TahunAjaran
    for rows.Next() {
        var t TahunAjaran
        err := rows.Scan(&t.IDTahun, &t.Nama, &t.Aktif)
        if err != nil {
            return nil, err
        }
        tahun = append(tahun, t)
    }
    return tahun, nil
}

func CreateTahunAjaran(nama string, aktif bool) error {
    query := `INSERT INTO tahun_ajaran (nama, aktif) VALUES (?, ?)`
    _, err := config.DB.Exec(query, nama, aktif)
    return err
}

func GetTahunAjaranByID(id int) (*TahunAjaran, error) {
    query := `SELECT id_tahun, nama, aktif FROM tahun_ajaran WHERE id_tahun = ?`
    row := config.DB.QueryRow(query, id)
    var t TahunAjaran
    err := row.Scan(&t.IDTahun, &t.Nama, &t.Aktif)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &t, nil
}

func GetAktifTahunAjaran() (*TahunAjaran, error) {
    query := `SELECT id_tahun, nama, aktif FROM tahun_ajaran WHERE aktif = TRUE LIMIT 1`
    row := config.DB.QueryRow(query)
    var t TahunAjaran
    err := row.Scan(&t.IDTahun, &t.Nama, &t.Aktif)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &t, nil
}

func UpdateTahunAjaran(id int, nama string, aktif bool) error {
    query := `UPDATE tahun_ajaran SET nama = ?, aktif = ? WHERE id_tahun = ?`
    _, err := config.DB.Exec(query, nama, aktif, id)
    return err
}

func DeleteTahunAjaran(id int) error {
    query := `DELETE FROM tahun_ajaran WHERE id_tahun = ?`
    _, err := config.DB.Exec(query, id)
    return err
}
