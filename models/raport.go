package models

import (
    "database/sql"
    "e_raport_digital/config"
)


type Raport struct {
    IDSiswa   int
    NamaSiswa string
    NamaKelas string
    NilaiList []RaportNilai
    RataRata  float64
}

type RaportNilai struct {
    IDNilai   int
    NamaMapel string
    Nilai     int
    Grade     string // Letter grade
}

// KonversiNilai converts numeric grade to letter grade
func KonversiNilai(nilai int) string {
    if nilai >= 90 {
        return "A"
    } else if nilai >= 80 {
        return "B"
    } else if nilai >= 70 {
        return "C"
    } else if nilai >= 60 {
        return "D"
    }
    return "E"
}

func GetRaportByUserID(userID int) (*Raport, error) {
    // First, get siswa details by user_id
    var siswaID int
    var namaSiswa, namaKelas string
    err := config.DB.QueryRow(`
        SELECT s.id_siswa, s.nama, k.nama_kelas 
        FROM siswa s 
        LEFT JOIN kelas k ON s.id_kelas = k.id_kelas 
        WHERE s.id_user = ?`, userID).Scan(&siswaID, &namaSiswa, &namaKelas)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // No siswa for this user
        }
        return nil, err
    }

    // Get nilai list
    query := `SELECT n.id_nilai, m.nama_mapel, n.nilai
              FROM nilai n
              JOIN mapel m ON n.id_mapel = m.id_mapel
              WHERE n.id_siswa = ?`
    rows, err := config.DB.Query(query, siswaID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var nilaiList []RaportNilai
    var totalNilai, count int
    for rows.Next() {
        var rn RaportNilai
        err := rows.Scan(&rn.IDNilai, &rn.NamaMapel, &rn.Nilai)
        if err != nil {
            return nil, err
        }
        rn.Grade = KonversiNilai(rn.Nilai)
        nilaiList = append(nilaiList, rn)
        totalNilai += rn.Nilai
        count++
    }

    rataRata := 0.0
    if count > 0 {
        rataRata = float64(totalNilai) / float64(count)
    }

    raport := &Raport{
        IDSiswa:   siswaID,
        NamaSiswa: namaSiswa,
        NamaKelas: namaKelas,
        NilaiList: nilaiList,
        RataRata:  rataRata,
    }
    return raport, nil
}
