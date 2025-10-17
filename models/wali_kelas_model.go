package models

import (
    "database/sql"

    "e_raport_digital/config"
)

func GetAllWaliKelas() ([]WaliKelas, error) {
    query := `SELECT id_wali, id_guru, id_kelas FROM wali_kelas`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var walis []WaliKelas
    for rows.Next() {
        var w WaliKelas
        err := rows.Scan(&w.IDWali, &w.IDGuru, &w.IDKelas)
        if err != nil {
            return nil, err
        }
        walis = append(walis, w)
    }
    return walis, nil
}

func CreateWaliKelas(idGuru, idKelas int) error {
    query := `INSERT INTO wali_kelas (id_guru, id_kelas) VALUES (?, ?)`
    _, err := config.DB.Exec(query, idGuru, idKelas)
    return err
}

func GetWaliByKelas(idKelas int) (*WaliKelas, error) {
    query := `SELECT id_wali, id_guru, id_kelas FROM wali_kelas WHERE id_kelas = ?`
    row := config.DB.QueryRow(query, idKelas)
    var w WaliKelas
    err := row.Scan(&w.IDWali, &w.IDGuru, &w.IDKelas)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &w, nil
}

func GetWaliByGuru(idGuru int) (*WaliKelas, error) {
    query := `SELECT id_wali, id_guru, id_kelas FROM wali_kelas WHERE id_guru = ?`
    row := config.DB.QueryRow(query, idGuru)
    var w WaliKelas
    err := row.Scan(&w.IDWali, &w.IDGuru, &w.IDKelas)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &w, nil
}

func DeleteWaliKelas(idWali int) error {
    query := `DELETE FROM wali_kelas WHERE id_wali = ?`
    _, err := config.DB.Exec(query, idWali)
    return err
}
