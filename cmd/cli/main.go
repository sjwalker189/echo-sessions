package main

import (
	model "app/db"
	hash "app/lib"
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		panic(err)
	}

	q := model.New(db)

	pw, err := hash.Make("password")
	if err != nil {
		log.Fatal(err)
	}

	err = q.NewUser(context.Background(), model.NewUserParams{
		FirstName: "Sam",
		LastName:  "Walker",
		Email:     "sjwalker189@gmail.com",
		Password:  pw,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Complete")
}
