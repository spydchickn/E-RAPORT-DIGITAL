// ...existing code...
package models

import "database/sql"

type Guru struct {
    IDGuru int
    Nama   string
    NIP    string
    Alamat string
    Foto   string
    IDUser sql.NullInt64
}
