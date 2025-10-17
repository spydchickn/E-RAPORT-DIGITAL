package models

import (
    "database/sql"

    "e_raport_digital/config"
)

func GetAllEkstrakurikuler() ([]Ekstrakurikuler, error) {
    query := `SELECT id_ekstra, id_siswa, nama_ekstra, nilai FROM ekstrakurikuler`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ekstras []Ekstrakurikuler
    for rows.Next() {
        var e Ekstrakurikuler
        err := rows.Scan(&e.IDEkstra, &e.IDSiswa, &e.NamaEkstra, &e.Nilai)
        if err != nil {
            return nil, err
        }
        ekstras = append(ekstras, e)
    }
    return ekstras, nil
}

func CreateEkstrakurikuler(idSiswa int, namaEkstra, nilai string) error {
    query := `INSERT INTO ekstrakurikuler (id_siswa, nama_ekstra, nilai) VALUES (?, ?, ?)`
    _, err := config.DB.Exec(query, idSiswa, namaEkstra, nilai)
    return err
}

func GetEkstrakurikulerByID(id int) (*Ekstrakurikuler, error) {
    query := `SELECT id_ekstra, id_siswa, nama_ekstra, nilai FROM ekstrakurikuler WHERE id_ekstra = ?`
    row := config.DB.QueryRow(query, id)
    var e Ekstrakurikuler
    err := row.Scan(&e.IDEkstra, &e.IDSiswa, &e.NamaEkstra, &e.Nilai)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &e, nil
}

func GetEkstrakurikulerBySiswa(idSiswa int) ([]Ekstrakurikuler, error) {
    query := `SELECT id_ekstra, id_siswa, nama_ekstra, nilai FROM ekstrakurikuler WHERE id_siswa = ? ORDER BY nama_ekstra`
    rows, err := config.DB.Query(query, idSiswa)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ekstras []Ekstrakurikuler
    for rows.Next() {
        var e Ekstrakurikuler
        err := rows.Scan(&e.IDEkstra, &e.IDSiswa, &e.NamaEkstra, &e.Nilai)
        if err != nil {
            return nil, err
        }
        ekstras = append(ekstras, e)
    }
    return ekstras, nil
}

func UpdateEkstrakurikuler(id int, idSiswa int, namaEkstra, nilai string) error {
    query := `UPDATE ekstrakurikuler SET id_siswa = ?, nama_ekstra = ?, nilai = ? WHERE id_ekstra = ?`
    _, err := config.DB.Exec(query, idSiswa, namaEkstra, nilai, id)
    return err
}

func DeleteEkstrakurikuler(id int) error {
    query := `DELETE FROM ekstrakurikuler WHERE id_ekstra = ?`
    _, err := config.DB.Exec(query, id)
    return err
}
