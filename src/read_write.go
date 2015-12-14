package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const READ_MAX = 900
const WRITE_MAX = 1000

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=geospat sslmode=disable")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	finished := make(chan int)

	for i := 0; i < READ_MAX; i++ {
		go readDatabase(db, i, finished)
	}

	for i := 0; i < WRITE_MAX; i++ {
		//go writeDatabase(db)
	}

	for i := 0; i < READ_MAX; i++ {
		res := <-finished
		fmt.Printf("Result %v\n", res)
	}
}

func readDatabase(db *sql.DB, i int, finished chan int) {
	rows, err := db.Query("SELECT count(*) from drivers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
	}
	finished <- count
	//fmt.Printf("Rows %v", count)
}

func writeDatabase(db *sql.DB) {
	result, err := db.Exec("UPDATE drivers set geog='POINT(77.6130574652 12.9018253405)' where driver_id=1")
	if err != nil {
		log.Fatal(err)
	}
	res, err := result.RowsAffected()
	fmt.Printf("RESULT %v", res)
}
