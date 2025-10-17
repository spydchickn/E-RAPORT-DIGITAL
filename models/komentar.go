package models

import "database/sql"

type Komentar struct {
    IDKomentar int
    IDNilai    sql.NullInt64
    IDSiswa    int
    DariRole   string
    IDPengirim int
    Pesan      string
    Tanggal    sql.NullTime
}
