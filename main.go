package main

import (
	"first_go_project/app/helpers"
	"first_go_project/app/migrations"
	"first_go_project/app/orm/models"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	// Определяем флаг для запуска миграций
	migrate := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	// Если флаг --migrate установлен, запускаем миграции и выходим
	if *migrate {
		err := migrations.RunMigrations()
		if err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully.")
		return
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	}).Methods("GET")

	r.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, models.Product{}.All().ToJson())
	}).Methods("GET")

	r.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, models.Product{}.Find(helpers.StringOf(mux.Vars(r)["id"]).ToInt()).ToJson())
	}).Methods("GET")

	// TODO категории не выводятся, нужно дописать with и many to many в builder .
	r.HandleFunc("/products/{id}/categories", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, models.Product{}.Find(
			helpers.StringOf(mux.Vars(r)["id"]).ToInt(),
		).Categories().ToJson())
	}).Methods("GET")

	r.HandleFunc("/products/first", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, models.Product{}.First().ToJson())
	}).Methods("GET")

	r.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, models.Category{}.All().ToJson())
	}).Methods("GET")

	r.HandleFunc("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, models.Category{}.Find(helpers.StringOf(mux.Vars(r)["id"]).ToInt()).ToJson())
	}).Methods("GET")

	http.ListenAndServe(":8080", r)
}
