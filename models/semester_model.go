package models

import (
    "database/sql"

    "e_raport_digital/config"
)

func GetAllSemester() ([]Semester, error) {
    query := `SELECT id_semester, nama, id_tahun_ajaran FROM semester`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var semesters []Semester
    for rows.Next() {
        var s Semester
        err := rows.Scan(&s.IDSemester, &s.Nama, &s.IDTahunAjaran)
        if err != nil {
            return nil, err
        }
        semesters = append(semesters, s)
    }
    return semesters, nil
}

func CreateSemester(nama string, idTahunAjaran int) error {
    query := `INSERT INTO semester (nama, id_tahun_ajaran) VALUES (?, ?)`
    _, err := config.DB.Exec(query, nama, idTahunAjaran)
    return err
}

func GetSemesterByID(id int) (*Semester, error) {
    query := `SELECT id_semester, nama, id_tahun_ajaran FROM semester WHERE id_semester = ?`
    row := config.DB.QueryRow(query, id)
    var s Semester
    err := row.Scan(&s.IDSemester, &s.Nama, &s.IDTahunAjaran)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &s, nil
}

func GetSemesterByTahun(idTahun int) ([]Semester, error) {
    query := `SELECT id_semester, nama, id_tahun_ajaran FROM semester WHERE id_tahun_ajaran = ?`
    rows, err := config.DB.Query(query, idTahun)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var semesters []Semester
    for rows.Next() {
        var s Semester
        err := rows.Scan(&s.IDSemester, &s.Nama, &s.IDTahunAjaran)
        if err != nil {
            return nil, err
        }
        semesters = append(semesters, s)
    }
    return semesters, nil
}

func GetAktifSemester() (*Semester, error) {
    // Assume active semester is tied to active tahun_ajaran
    query := `SELECT s.id_semester, s.nama, s.id_tahun_ajaran 
              FROM semester s 
              JOIN tahun_ajaran t ON s.id_tahun_ajaran = t.id_tahun 
              WHERE t.aktif = TRUE LIMIT 1`
    row := config.DB.QueryRow(query)
    var s Semester
    err := row.Scan(&s.IDSemester, &s.Nama, &s.IDTahunAjaran)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &s, nil
}

func UpdateSemester(id int, nama string, idTahunAjaran int) error {
    query := `UPDATE semester SET nama = ?, id_tahun_ajaran = ? WHERE id_semester = ?`
    _, err := config.DB.Exec(query, nama, idTahunAjaran, id)
    return err
}

func DeleteSemester(id int) error {
    query := `DELETE FROM semester WHERE id_semester = ?`
    _, err := config.DB.Exec(query, id)
    return err
}
