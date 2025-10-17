package models

import "e_raport_digital/config"

func GetAllMapel() ([]Mapel, error) {
    rows, err := config.DB.Query("SELECT id_mapel, nama_mapel, kkm FROM mapel")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var list []Mapel
    for rows.Next() {
        var m Mapel
        rows.Scan(&m.IDMapel, &m.NamaMapel, &m.KKM)
        list = append(list, m)
    }
    return list, nil
}

func CreateMapel(nama string, kkm string) error {
    _, err := config.DB.Exec("INSERT INTO mapel (nama_mapel, kkm) VALUES (?, ?)", nama, kkm)
    return err
}

func UpdateMapel(id int, nama string) error {
    _, err := config.DB.Exec("UPDATE mapel SET nama_mapel=? WHERE id_mapel=?", nama, id)
    return err
}

func DeleteMapel(id int) error {
    _, err := config.DB.Exec("DELETE FROM mapel WHERE id_mapel=?", id)
    return err
}
