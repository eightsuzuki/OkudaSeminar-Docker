package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type TestUser struct {
	UserID   int
	Password string
}

func main() {
	http.HandleFunc("/a", handler)
	http.ListenAndServe(":8181", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	var Db *sql.DB
	Db, err := sql.Open("postgres", "host=postgres user=app_user password=password dbname=app_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	sql := "SELECT user_id, user_password FROM TEST_USER WHERE user_id=$1;"

	pstatement, err := Db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}

	queryID := 1
	var testUser TestUser

	err = pstatement.QueryRow(queryID).Scan(&testUser.UserID, &testUser.Password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, testUser.Password)
}
