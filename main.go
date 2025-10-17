package main

import (
    "fmt"
    "log"
    "net/http"
    "runtime/debug"

    "e_raport_digital/config"
    "e_raport_digital/routes"
)

// recoverMiddleware wraps an http.Handler and recovers from panics,
// logs the stack trace and returns HTTP 500 with a clear message.
func recoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                log.Printf("PANIC recovered: %v\n%s", rec, debug.Stack())
                // In development, show stack trace in response for easier debugging.
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprintf(w, "Internal Server Error: %v\n\n%s", rec, debug.Stack())
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func main() {
    config.ConnectDB()
    defer func() {
        if config.DB != nil {
            config.DB.Close()
        }
    }()

    r := routes.SetupRoutes()
    // wrap router with recover middleware so panics get logged and return 500
    handler := recoverMiddleware(r)

    log.Println("Server running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Fatal(err)
    }
}
