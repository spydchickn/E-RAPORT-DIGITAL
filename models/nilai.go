package models

import (
    "database/sql"
    "e_raport_digital/config"
)

type Nilai struct {
    IDNilai  int
    IDSiswa  int
    IDMapel  int
    Nilai    int
    Status   string
    Komentar string
}

func CreateNilai(idSiswa, idMapel, nilai int, komentar string) error {
    _, err := config.DB.Exec("INSERT INTO nilai (id_siswa, id_mapel, nilai, komentar, status) VALUES (?, ?, ?, ?, 'draft')", idSiswa, idMapel, nilai, komentar)
    return err
}

func GetNilaiBySiswa(idSiswa int) ([]Nilai, error) {
    rows, err := config.DB.Query("SELECT id_nilai, id_siswa, id_mapel, nilai, status, komentar FROM nilai WHERE id_siswa = ?", idSiswa)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var list []Nilai
    for rows.Next() {
        var n Nilai
        rows.Scan(&n.IDNilai, &n.IDSiswa, &n.IDMapel, &n.Nilai, &n.Status, &n.Komentar)
        list = append(list, n)
    }
    return list, nil
}

func GetNilaiByID(id int) (*Nilai, error) {
    row := config.DB.QueryRow("SELECT id_nilai, id_siswa, id_mapel, nilai, status, komentar FROM nilai WHERE id_nilai = ?", id)
    var n Nilai
    err := row.Scan(&n.IDNilai, &n.IDSiswa, &n.IDMapel, &n.Nilai, &n.Status, &n.Komentar)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &n, nil
}

func UpdateNilaiStatus(id int, status string) error {
    _, err := config.DB.Exec("UPDATE nilai SET status = ? WHERE id_nilai = ?", status, id)
    return err
}
