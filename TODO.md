# TODO: Repair verifikasiController.go errors

## Steps:
- [x] Step 1: Add GetSiswaByID function to models/siswa.go
- [x] Step 2: Update VerifikasiStore in controllers/verifikasiController.go to properly use sql.NullInt64 and sql.NullString for CreateNotif and LogAction calls
