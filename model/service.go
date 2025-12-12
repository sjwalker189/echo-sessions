package model

import (
	"github.com/jmoiron/sqlx"
)

type Service struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Desc string `db:"description"`
}

func ListServices(db *sqlx.DB) ([]Service, error) {
	var services []Service
	err := db.Select(&services, "SELECT * from services")
	return services, err
}
