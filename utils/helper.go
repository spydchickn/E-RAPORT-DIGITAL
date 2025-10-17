package utils

import (
    "net/http"

    "e_raport_digital/middlewares"
)

// Helper function placeholder
func KonversiNilai(nilai int) string {
    if nilai >= 90 {
        return "A"
    } else if nilai >= 80 {
        return "B"
    } else if nilai >= 70 {
        return "C"
    } else if nilai >= 60 {
        return "D"
    }
    return "E"
}

// Helper function to get user ID from session
func GetUserIDFromSession(r *http.Request) int {
    session, _ := middlewares.Store.Get(r, "session")
    if session.Values["user_id"] != nil {
        return session.Values["user_id"].(int)
    }
    return 0
}
