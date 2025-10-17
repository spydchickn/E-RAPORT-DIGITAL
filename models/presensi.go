package models

import (
    "database/sql"
)

type Presensi struct {
    IDPresensi int
    IDSiswa    int
    IDMapel    int
    IDGuru     int
    Tanggal    sql.NullTime  // For DATE field
    Status     string
    Catatan    sql.NullString
}
