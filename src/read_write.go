package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

const READ_MAX = 14999

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
	latMin := 12.813196
	latMax := 13.055798

	longMin := 77.474313
	longMax := 77.767158

	randomLat := randomBetween(latMin, latMax)
	randomLong := randomBetween(longMin, longMax)

	query := fmt.Sprintf("select driver_id from drivers d where ST_DWithin(d.geog, ST_GeomFromText('POINT(%v %v)'), 1500) limit 10", randomLat, randomLong)
	rows, err := db.Query(query)

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

func randomBetween(min, max float64) float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Float64()*(max-min) + min
}
