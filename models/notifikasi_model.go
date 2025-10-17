package models

import (
    "database/sql"
    "fmt"

    "e_raport_digital/config"
)

func GetAllNotifikasi() ([]Notifikasi, error) {
    query := `SELECT id_notif, to_user_id, to_role, pesan, dibaca, tanggal FROM notifikasi ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifs []Notifikasi
    for rows.Next() {
        var n Notifikasi
        err := rows.Scan(&n.IDNotif, &n.ToUserID, &n.ToRole, &n.Pesan, &n.Dibaca, &n.Tanggal)
        if err != nil {
            return nil, err
        }
        notifs = append(notifs, n)
    }
    return notifs, nil
}

func CreateNotif(toUserID sql.NullInt64, toRole sql.NullString, pesan string) error {
    var toUserIDInt *int
    if toUserID.Valid {
        intVal := int(toUserID.Int64)
        toUserIDInt = &intVal
    }
    var toRoleStr *string
    if toRole.Valid {
        toRoleStr = &toRole.String
    }
    if toUserIDInt != nil && toRoleStr != nil {
        query := `INSERT INTO notifikasi (to_user_id, to_role, pesan) VALUES (?, ?, ?)`
        _, err := config.DB.Exec(query, *toUserIDInt, *toRoleStr, pesan)
        return err
    } else if toUserIDInt != nil {
        query := `INSERT INTO notifikasi (to_user_id, pesan) VALUES (?, ?)`
        _, err := config.DB.Exec(query, *toUserIDInt, pesan)
        return err
    } else if toRoleStr != nil {
        query := `INSERT INTO notifikasi (to_role, pesan) VALUES (?, ?)`
        _, err := config.DB.Exec(query, *toRoleStr, pesan)
        return err
    } else {
        return fmt.Errorf("either to_user_id or to_role must be provided")
    }
}

func GetNotifByID(id int) (*Notifikasi, error) {
    query := `SELECT id_notif, to_user_id, to_role, pesan, dibaca, tanggal FROM notifikasi WHERE id_notif = ?`
    row := config.DB.QueryRow(query, id)
    var n Notifikasi
    err := row.Scan(&n.IDNotif, &n.ToUserID, &n.ToRole, &n.Pesan, &n.Dibaca, &n.Tanggal)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &n, nil
}

func GetUnreadByUser(userID int) ([]Notifikasi, error) {
    query := `SELECT id_notif, to_user_id, to_role, pesan, dibaca, tanggal FROM notifikasi WHERE to_user_id = ? AND dibaca = FALSE ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifs []Notifikasi
    for rows.Next() {
        var n Notifikasi
        err := rows.Scan(&n.IDNotif, &n.ToUserID, &n.ToRole, &n.Pesan, &n.Dibaca, &n.Tanggal)
        if err != nil {
            return nil, err
        }
        notifs = append(notifs, n)
    }
    return notifs, nil
}

func GetNotifByRole(role string) ([]Notifikasi, error) {
    query := `SELECT id_notif, to_user_id, to_role, pesan, dibaca, tanggal FROM notifikasi WHERE to_role = ? ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query, role)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifs []Notifikasi
    for rows.Next() {
        var n Notifikasi
        err := rows.Scan(&n.IDNotif, &n.ToUserID, &n.ToRole, &n.Pesan, &n.Dibaca, &n.Tanggal)
        if err != nil {
            return nil, err
        }
        notifs = append(notifs, n)
    }
    return notifs, nil
}

func MarkRead(id int) error {
    query := `UPDATE notifikasi SET dibaca = TRUE WHERE id_notif = ?`
    _, err := config.DB.Exec(query, id)
    return err
}

func MarkAllReadByUser(userID int) error {
    query := `UPDATE notifikasi SET dibaca = TRUE WHERE to_user_id = ? AND dibaca = FALSE`
    _, err := config.DB.Exec(query, userID)
    return err
}

func DeleteNotif(id int) error {
    query := `DELETE FROM notifikasi WHERE id_notif = ?`
    _, err := config.DB.Exec(query, id)
    return err
}
