package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const READ_MAX = 4999

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=cabspike port=6543 sslmode=disable password=123456")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	finished := make(chan bool)

	for i := 0; i < READ_MAX; i++ {
		go readDatabase(db, finished)
	}

	var total_read_time int64
	for i := 0; i < READ_MAX; i++ {
		start := time.Now()
		<-finished
		total_read_time += time.Since(start).Nanoseconds()
		start = time.Now()
	}
	avg_time_ns := total_read_time / READ_MAX
	fmt.Printf("Avg read time: %v ms\n", avg_time_ns/int64(time.Millisecond))
	fmt.Printf("Total read time for %v connections: %v ms\n", READ_MAX, total_read_time/int64(time.Millisecond))
}

func readDatabase(db *sql.DB, finished chan bool) {
	rows, err := db.Query("select driver_id from drivers d where ST_DWithin(d.geog, ST_GeomFromText('POINT(12.99612 77.57553)'), 1500) limit 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			fmt.Println(err)
		}
	}

	finished <- true
}

func writeDatabase(db *sql.DB) {
	result, err := db.Exec("UPDATE drivers set geog='POINT(77.6130574652 12.9018253405)' where driver_id=1")

	if err != nil {
		log.Fatal(err)
	}
	res, err := result.RowsAffected()
	fmt.Printf("Rows affected %v", res)
}
