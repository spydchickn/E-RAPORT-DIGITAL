package models

import "database/sql"

type Notifikasi struct {
    IDNotif    int
    ToUserID   sql.NullInt64
    ToRole     sql.NullString
    Pesan      string
    Dibaca     bool
    Tanggal    sql.NullTime
}
