package models

import (
    "database/sql"
    "time"

    "e_raport_digital/config"
)

func GetAllPresensi() ([]Presensi, error) {
    query := `SELECT id_presensi, id_siswa, id_mapel, id_guru, tanggal, status, catatan FROM presensi`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var presensis []Presensi
    for rows.Next() {
        var p Presensi
        err := rows.Scan(&p.IDPresensi, &p.IDSiswa, &p.IDMapel, &p.IDGuru, &p.Tanggal, &p.Status, &p.Catatan)
        if err != nil {
            return nil, err
        }
        presensis = append(presensis, p)
    }
    return presensis, nil
}

func CreatePresensi(idSiswa, idMapel, idGuru int, tanggal string, status, catatan string) error {
    query := `INSERT INTO presensi (id_siswa, id_mapel, id_guru, tanggal, status, catatan) VALUES (?, ?, ?, ?, ?, ?)`
    _, err := config.DB.Exec(query, idSiswa, idMapel, idGuru, tanggal, status, catatan)
    return err
}

func GetPresensiByID(id int) (*Presensi, error) {
    query := `SELECT id_presensi, id_siswa, id_mapel, id_guru, tanggal, status, catatan FROM presensi WHERE id_presensi = ?`
    row := config.DB.QueryRow(query, id)
    var p Presensi
    err := row.Scan(&p.IDPresensi, &p.IDSiswa, &p.IDMapel, &p.IDGuru, &p.Tanggal, &p.Status, &p.Catatan)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &p, nil
}

func GetPresensiBySiswa(idSiswa int) ([]Presensi, error) {
    query := `SELECT id_presensi, id_siswa, id_mapel, id_guru, tanggal, status, catatan FROM presensi WHERE id_siswa = ? ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query, idSiswa)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var presensis []Presensi
    for rows.Next() {
        var p Presensi
        err := rows.Scan(&p.IDPresensi, &p.IDSiswa, &p.IDMapel, &p.IDGuru, &p.Tanggal, &p.Status, &p.Catatan)
        if err != nil {
            return nil, err
        }
        presensis = append(presensis, p)
    }
    return presensis, nil
}

func GetPresensiByMapel(idMapel int) ([]Presensi, error) {
    query := `SELECT id_presensi, id_siswa, id_mapel, id_guru, tanggal, status, catatan FROM presensi WHERE id_mapel = ? ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query, idMapel)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var presensis []Presensi
    for rows.Next() {
        var p Presensi
        err := rows.Scan(&p.IDPresensi, &p.IDSiswa, &p.IDMapel, &p.IDGuru, &p.Tanggal, &p.Status, &p.Catatan)
        if err != nil {
            return nil, err
        }
        presensis = append(presensis, p)
    }
    return presensis, nil
}

func GetPresensiByTanggal(tanggal time.Time) ([]Presensi, error) {
    query := `SELECT id_presensi, id_siswa, id_mapel, id_guru, tanggal, status, catatan FROM presensi WHERE tanggal = ? ORDER BY id_siswa`
    rows, err := config.DB.Query(query, tanggal.Format("2006-01-02"))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var presensis []Presensi
    for rows.Next() {
        var p Presensi
        err := rows.Scan(&p.IDPresensi, &p.IDSiswa, &p.IDMapel, &p.IDGuru, &p.Tanggal, &p.Status, &p.Catatan)
        if err != nil {
            return nil, err
        }
        presensis = append(presensis, p)
    }
    return presensis, nil
}

func GetRekapPresensi(idKelas int, tanggal time.Time) (map[string]int, error) {
    query := `SELECT p.status, COUNT(*) 
              FROM presensi p 
              JOIN siswa s ON p.id_siswa = s.id_siswa 
              WHERE s.id_kelas = ? AND p.tanggal = ? 
              GROUP BY p.status`
    rows, err := config.DB.Query(query, idKelas, tanggal.Format("2006-01-02"))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    rekap := make(map[string]int)
    for rows.Next() {
        var status string
        var count int
        err := rows.Scan(&status, &count)
        if err != nil {
            return nil, err
        }
        rekap[status] = count
    }
    return rekap, nil
}

func UpdatePresensi(id int, idSiswa, idMapel, idGuru int, tanggal string, status, catatan string) error {
    query := `UPDATE presensi SET id_siswa = ?, id_mapel = ?, id_guru = ?, tanggal = ?, status = ?, catatan = ? WHERE id_presensi = ?`
    _, err := config.DB.Exec(query, idSiswa, idMapel, idGuru, tanggal, status, catatan, id)
    return err
}

func DeletePresensi(id int) error {
    query := `DELETE FROM presensi WHERE id_presensi = ?`
    _, err := config.DB.Exec(query, id)
    return err
}
