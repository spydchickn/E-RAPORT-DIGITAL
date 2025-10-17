package models

import (
    "database/sql"
    "time"

    "e_raport_digital/config"
)

func GetAllLogs() ([]LogsAktivitas, error) {
    query := `SELECT id_log, user_id, action, table_name, record_id, details, tanggal FROM logs_aktivitas ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var logs []LogsAktivitas
    for rows.Next() {
        var l LogsAktivitas
        err := rows.Scan(&l.IDLog, &l.UserID, &l.Action, &l.TableName, &l.RecordID, &l.Details, &l.Tanggal)
        if err != nil {
            return nil, err
        }
        logs = append(logs, l)
    }
    return logs, nil
}

func LogAction(userID int, action, tableName string, recordID sql.NullInt64, details sql.NullString) error {
    var recordIDInt *int
    if recordID.Valid {
        intVal := int(recordID.Int64)
        recordIDInt = &intVal
    }
    var detailsStr *string
    if details.Valid {
        detailsStr = &details.String
    }
    if recordIDInt != nil && detailsStr != nil {
        query := `INSERT INTO logs_aktivitas (user_id, action, table_name, record_id, details) VALUES (?, ?, ?, ?, ?)`
        _, err := config.DB.Exec(query, userID, action, tableName, *recordIDInt, *detailsStr)
        return err
    } else if recordIDInt != nil {
        query := `INSERT INTO logs_aktivitas (user_id, action, table_name, record_id) VALUES (?, ?, ?, ?)`
        _, err := config.DB.Exec(query, userID, action, tableName, *recordIDInt)
        return err
    } else if detailsStr != nil {
        query := `INSERT INTO logs_aktivitas (user_id, action, table_name, details) VALUES (?, ?, ?, ?)`
        _, err := config.DB.Exec(query, userID, action, tableName, *detailsStr)
        return err
    } else {
        query := `INSERT INTO logs_aktivitas (user_id, action, table_name) VALUES (?, ?, ?)`
        _, err := config.DB.Exec(query, userID, action, tableName)
        return err
    }
}

func GetLogByID(id int) (*LogsAktivitas, error) {
    query := `SELECT id_log, user_id, action, table_name, record_id, details, tanggal FROM logs_aktivitas WHERE id_log = ?`
    row := config.DB.QueryRow(query, id)
    var l LogsAktivitas
    err := row.Scan(&l.IDLog, &l.UserID, &l.Action, &l.TableName, &l.RecordID, &l.Details, &l.Tanggal)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &l, nil
}

func GetLogsByUser(userID int) ([]LogsAktivitas, error) {
    query := `SELECT id_log, user_id, action, table_name, record_id, details, tanggal FROM logs_aktivitas WHERE user_id = ? ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var logs []LogsAktivitas
    for rows.Next() {
        var l LogsAktivitas
        err := rows.Scan(&l.IDLog, &l.UserID, &l.Action, &l.TableName, &l.RecordID, &l.Details, &l.Tanggal)
        if err != nil {
            return nil, err
        }
        logs = append(logs, l)
    }
    return logs, nil
}

func GetLogsByDate(tanggal time.Time) ([]LogsAktivitas, error) {
    query := `SELECT id_log, user_id, action, table_name, record_id, details, tanggal FROM logs_aktivitas WHERE DATE(tanggal) = ? ORDER BY tanggal DESC`
    rows, err := config.DB.Query(query, tanggal.Format("2006-01-02"))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var logs []LogsAktivitas
    for rows.Next() {
        var l LogsAktivitas
        err := rows.Scan(&l.IDLog, &l.UserID, &l.Action, &l.TableName, &l.RecordID, &l.Details, &l.Tanggal)
        if err != nil {
            return nil, err
        }
        logs = append(logs, l)
    }
    return logs, nil
}

func GetRecentLogs(limit int) ([]LogsAktivitas, error) {
    query := `SELECT id_log, user_id, action, table_name, record_id, details, tanggal FROM logs_aktivitas ORDER BY tanggal DESC LIMIT ?`
    rows, err := config.DB.Query(query, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var logs []LogsAktivitas
    for rows.Next() {
        var l LogsAktivitas
        err := rows.Scan(&l.IDLog, &l.UserID, &l.Action, &l.TableName, &l.RecordID, &l.Details, &l.Tanggal)
        if err != nil {
            return nil, err
        }
        logs = append(logs, l)
    }
    return logs, nil
}

func DeleteLog(id int) error {
    query := `DELETE FROM logs_aktivitas WHERE id_log = ?`
    _, err := config.DB.Exec(query, id)
    return err
}
