package main

import (
    "log"
    "net/http"
            
    "e_raport_digital/config"
    "e_raport_digital/routes"
)

func main() {
    config.ConnectDB()
    defer config.DB.Close()

    r := routes.SetupRoutes()
    log.Println("Server running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
