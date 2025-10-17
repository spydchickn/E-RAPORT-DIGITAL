package models

import (
    "database/sql"

    "e_raport_digital/config"
)

func GetAllCatatanSikap() ([]CatatanSikap, error) {
    query := `SELECT id_catatan, id_siswa, id_semester, deskripsi, nilai_sikap FROM catatan_sikap`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var catatans []CatatanSikap
    for rows.Next() {
        var c CatatanSikap
        err := rows.Scan(&c.IDCatatan, &c.IDSiswa, &c.IDSemester, &c.Deskripsi, &c.NilaiSikap)
        if err != nil {
            return nil, err
        }
        catatans = append(catatans, c)
    }
    return catatans, nil
}

func CreateCatatanSikap(idSiswa, idSemester int, deskripsi, nilaiSikap string) error {
    query := `INSERT INTO catatan_sikap (id_siswa, id_semester, deskripsi, nilai_sikap) VALUES (?, ?, ?, ?)`
    _, err := config.DB.Exec(query, idSiswa, idSemester, deskripsi, nilaiSikap)
    return err
}

func GetCatatanSikapByID(id int) (*CatatanSikap, error) {
    query := `SELECT id_catatan, id_siswa, id_semester, deskripsi, nilai_sikap FROM catatan_sikap WHERE id_catatan = ?`
    row := config.DB.QueryRow(query, id)
    var c CatatanSikap
    err := row.Scan(&c.IDCatatan, &c.IDSiswa, &c.IDSemester, &c.Deskripsi, &c.NilaiSikap)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &c, nil
}

func GetCatatanSikapBySiswa(idSiswa int) ([]CatatanSikap, error) {
    query := `SELECT id_catatan, id_siswa, id_semester, deskripsi, nilai_sikap FROM catatan_sikap WHERE id_siswa = ? ORDER BY id_catatan DESC`
    rows, err := config.DB.Query(query, idSiswa)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var catatans []CatatanSikap
    for rows.Next() {
        var c CatatanSikap
        err := rows.Scan(&c.IDCatatan, &c.IDSiswa, &c.IDSemester, &c.Deskripsi, &c.NilaiSikap)
        if err != nil {
            return nil, err
        }
        catatans = append(catatans, c)
    }
    return catatans, nil
}

func GetCatatanSikapBySemester(idSemester int) ([]CatatanSikap, error) {
    query := `SELECT id_catatan, id_siswa, id_semester, deskripsi, nilai_sikap FROM catatan_sikap WHERE id_semester = ? ORDER BY id_siswa`
    rows, err := config.DB.Query(query, idSemester)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var catatans []CatatanSikap
    for rows.Next() {
        var c CatatanSikap
        err := rows.Scan(&c.IDCatatan, &c.IDSiswa, &c.IDSemester, &c.Deskripsi, &c.NilaiSikap)
        if err != nil {
            return nil, err
        }
        catatans = append(catatans, c)
    }
    return catatans, nil
}

func UpdateCatatanSikap(id int, idSiswa, idSemester int, deskripsi, nilaiSikap string) error {
    query := `UPDATE catatan_sikap SET id_siswa = ?, id_semester = ?, deskripsi = ?, nilai_sikap = ? WHERE id_catatan = ?`
    _, err := config.DB.Exec(query, idSiswa, idSemester, deskripsi, nilaiSikap, id)
    return err
}

func DeleteCatatanSikap(id int) error {
    query := `DELETE FROM catatan_sikap WHERE id_catatan = ?`
    _, err := config.DB.Exec(query, id)
    return err
}
