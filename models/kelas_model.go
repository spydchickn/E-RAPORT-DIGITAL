package models

import "e_raport_digital/config"

func GetAllKelas() ([]Kelas, error) {
    rows, err := config.DB.Query("SELECT id_kelas, nama_kelas FROM kelas")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var list []Kelas
    for rows.Next() {
        var k Kelas
        rows.Scan(&k.IDKelas, &k.NamaKelas)
        list = append(list, k)
    }
    return list, nil
}

func CreateKelas(nama string) error {
    _, err := config.DB.Exec("INSERT INTO kelas (nama_kelas) VALUES (?)", nama)
    return err
}

func UpdateKelas(id int, nama string) error {
    _, err := config.DB.Exec("UPDATE kelas SET nama_kelas=? WHERE id_kelas=?", nama, id)
    return err
}

func DeleteKelas(id int) error {
    _, err := config.DB.Exec("DELETE FROM kelas WHERE id_kelas=?", id)
    return err
}
