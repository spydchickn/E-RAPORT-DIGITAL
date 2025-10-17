package models

import (
    "database/sql"
    "time"
)

type Presensi struct {
    IDPresensi int
    IDSiswa    int
    IDMapel    int
    IDGuru     int
    Tanggal    time.Time
    Status     string
    Catatan    sql.NullString
}
