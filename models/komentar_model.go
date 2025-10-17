package models

import (
    "database/sql"

    "e_raport_digital/config"
)

func GetAllKomentar() ([]Komentar, error) {
    query := `SELECT id_komentar, id_nilai, id_siswa, dari_role, id_pengirim, pesan, tanggal FROM komentar`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var komentars []Komentar
    for rows.Next() {
        var k Komentar
        err := rows.Scan(&k.IDKomentar, &k.IDNilai, &k.IDSiswa, &k.DariRole, &k.IDPengirim, &k.Pesan, &k.Tanggal)
        if err != nil {
            return nil, err
        }
        komentars = append(komentars, k)
    }
    return komentars, nil
}

func CreateKomentar(idNilai sql.NullInt64, idSiswa int, dariRole string, idPengirim int, pesan string) error {
    var idNilaiInt *int
    if idNilai.Valid {
        intVal := int(idNilai.Int64)
        idNilaiInt = &intVal
    }
    if idNilaiInt != nil {
        query := `INSERT INTO komentar (id_nilai, id_siswa, dari_role, id_pengirim, pesan) VALUES (?, ?, ?, ?, ?)`
        _, err := config.DB.Exec(query, *idNilaiInt, idSiswa, dariRole, idPengirim, pesan)
        return err
    } else {
        query := `INSERT INTO komentar (id_siswa, dari_role, id_pengirim, pesan) VALUES (?, ?, ?, ?)`
        _, err := config.DB.Exec(query, idSiswa, dariRole, idPengirim, pesan)
        return err
    }
}

func GetKomentarByID(id int) (*Komentar, error) {
    query := `SELECT id_komentar, id_nilai, id_siswa, dari_role, id_pengirim, pesan, tanggal FROM komentar WHERE id_komentar = ?`
    row := config.DB.QueryRow(query, id)
    var k Komentar
    err := row.Scan(&k.IDKomentar, &k.IDNilai, &k.IDSiswa, &k.DariRole, &k.IDPengirim, &k.Pesan, &k.Tanggal)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &k, nil
}

func GetKomentarBySiswa(idSiswa int) ([]Komentar, error) {
    query := `SELECT id_komentar, id_nilai, id_siswa, dari_role, id_pengirim, pesan, tanggal FROM komentar WHERE id_siswa = ? ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query, idSiswa)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var komentars []Komentar
    for rows.Next() {
        var k Komentar
        err := rows.Scan(&k.IDKomentar, &k.IDNilai, &k.IDSiswa, &k.DariRole, &k.IDPengirim, &k.Pesan, &k.Tanggal)
        if err != nil {
            return nil, err
        }
        komentars = append(komentars, k)
    }
    return komentars, nil
}

func GetKomentarByNilai(idNilai int) ([]Komentar, error) {
    query := `SELECT id_komentar, id_nilai, id_siswa, dari_role, id_pengirim, pesan, tanggal FROM komentar WHERE id_nilai = ? ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query, idNilai)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var komentars []Komentar
    for rows.Next() {
        var k Komentar
        err := rows.Scan(&k.IDKomentar, &k.IDNilai, &k.IDSiswa, &k.DariRole, &k.IDPengirim, &k.Pesan, &k.Tanggal)
        if err != nil {
            return nil, err
        }
        komentars = append(komentars, k)
    }
    return komentars, nil
}

func UpdateKomentar(id int, idNilai sql.NullInt64, idSiswa int, dariRole string, idPengirim int, pesan string) error {
    var idNilaiInt *int
    if idNilai.Valid {
        intVal := int(idNilai.Int64)
        idNilaiInt = &intVal
    }
    if idNilaiInt != nil {
        query := `UPDATE komentar SET id_nilai = ?, id_siswa = ?, dari_role = ?, id_pengirim = ?, pesan = ? WHERE id_komentar = ?`
        _, err := config.DB.Exec(query, *idNilaiInt, idSiswa, dariRole, idPengirim, pesan, id)
        return err
    } else {
        query := `UPDATE komentar SET id_siswa = ?, dari_role = ?, id_pengirim = ?, pesan = ? WHERE id_komentar = ?`
        _, err := config.DB.Exec(query, idSiswa, dariRole, idPengirim, pesan, id)
        return err
    }
}

func DeleteKomentar(id int) error {
    query := `DELETE FROM komentar WHERE id_komentar = ?`
    _, err := config.DB.Exec(query, id)
    return err
}
