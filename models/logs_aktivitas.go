package models

import "database/sql"

type LogsAktivitas struct {
    IDLog      int
    UserID     int
    Action     string
    TableName  string
    RecordID   sql.NullInt64
    Details    sql.NullString
    Tanggal    sql.NullTime
}
