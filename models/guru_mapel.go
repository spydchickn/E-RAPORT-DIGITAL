package models

import (
    "e_raport_digital/config"
)

// GuruMapel represents the assignment of a subject to a teacher
type GuruMapel struct {
    IDGuruMapel int
    IDGuru      int
    IDMapel     int
}

// GetMapelByGuruID returns all subjects assigned to a specific guru
func GetMapelByGuruID(guruID int) ([]Mapel, error) {
    query := `
        SELECT m.id_mapel, m.nama_mapel 
        FROM guru_mapel gm
        JOIN mapel m ON gm.id_mapel = m.id_mapel
        WHERE gm.id_guru = ?
    `
    rows, err := config.DB.Query(query, guruID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var mapels []Mapel
    for rows.Next() {
        var m Mapel
        err := rows.Scan(&m.IDMapel, &m.NamaMapel)
        if err != nil {
            return nil, err
        }
        mapels = append(mapels, m)
    }
    return mapels, nil
}

// AssignGuruMapel assigns a mapel to a guru (inserts into junction table)
func AssignGuruMapel(guruID, mapelID int) error {
    query := `INSERT INTO guru_mapel (id_guru, id_mapel) VALUES (?, ?)`
    _, err := config.DB.Exec(query, guruID, mapelID)
    if err != nil {
        return err
    }
    return nil
}

// RemoveGuruMapel removes the assignment of a mapel from a guru
func RemoveGuruMapel(guruID, mapelID int) error {
    query := `DELETE FROM guru_mapel WHERE id_guru = ? AND id_mapel = ?`
    _, err := config.DB.Exec(query, guruID, mapelID)
    if err != nil {
        return err
    }
    return nil
}
