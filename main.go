package main

import (
	"net/http"
	"github.com/heshambahamish/CRUD/db"
        "github.com/heshambahamish/CRUD/handlers"
)

func main() {
	db.Init()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/create", handlers.Create)
	http.HandleFunc("/store", handlers.Store)
	http.HandleFunc("/edit", handlers.Edit)
	http.HandleFunc("/update", handlers.Update)
	http.HandleFunc("/delete", handlers.Delete)

	http.ListenAndServe(":8080", nil)
}
