package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "ec2-54-228-139-34.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "oelbjbkvvncwcm"
	password = "623be7f9c2232ad797503ff14d82aa87ef433962a2d17b60c0f15d3ff8817187"
	dbname   = "d6cgt8qmojrjpn"
)

func Init() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		return
	}

	// close database
	defer db.Close()

	// check db
	err = db.Ping()

	if err != nil {
		return
	}

	fmt.Println("Connected!")
}
