package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const readMax = 15000
const writeMax = 8000
const database = "postgres"

func main() {
	db, err := getDatabase(database)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	finished := make(chan bool)

	for i := 0; i < readMax; i++ {
		go readDatabase(db, finished)
	}

	for i := 0; i < writeMax; i++ {
		go writeDatabase(db, finished)
	}

	var totalTime int64
	for i := 0; i < readMax+writeMax; i++ {
		start := time.Now()
		<-finished
		totalTime += time.Since(start).Nanoseconds()
		start = time.Now()
	}
	avgTimeInNS := totalTime / readMax
	fmt.Printf("Avg read time: %v ms\n", avgTimeInNS/int64(time.Millisecond))
	fmt.Printf("Total read time for %v connections: %v ms\n", readMax+writeMax, totalTime/int64(time.Millisecond))
}

func readDatabase(db *sql.DB, finished chan bool) {
	query := getReadQuery(database)
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

func writeDatabase(db *sql.DB, finished chan bool) {
	randomLat := randLat()
	randomLong := randLong()
	driverID := randDriverID()

	query := fmt.Sprintf("UPDATE drivers set geog='POINT(%v %v)' where driver_id=%v", randomLat, randomLong, driverID)
	result, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
	result.RowsAffected()
	//fmt.Printf("Rows affected %v", res)
	finished <- true
}

func randomBetween(min, max float64) float64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Float64()*(max-min) + min
}

func randLat() float64 {
	latMin := 12.813196
	latMax := 13.055798
	return randomBetween(latMin, latMax)
}

func randLong() float64 {
	longMin := 77.474313
	longMax := 77.767158
	return randomBetween(longMin, longMax)
}

func randDriverID() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(100000) + 1
}

func getDatabase(dbname string) (*sql.DB, error) {
	if dbname == "postgres" {
		return sql.Open("postgres", "user=postgres dbname=cabspike port=6543 sslmode=disable password=123456")
	}
	if dbname == "memsql" {
		return sql.Open("mysql", "root:@tcp(localhost:3037)/cabspike")
	}
	return nil, nil
}

func getReadQuery(dbname string) string {
	randomLat := randLat()
	randomLong := randLong()

	if dbname == "postgres" {
		return fmt.Sprintf("select driver_id from drivers d where ST_DWithin(d.geog, ST_GeomFromText('POINT(%v %v)'), 1500) limit 10", randomLat, randomLong)
	}
	if dbname == "memsql" {
		return fmt.Sprintf("select driver_id from drivers d where GEOGRAPHY_DISTANCE_WITHIN(d.geog, ST_GeomFromText('POINT(%v %v)'), 1500) limit 10", randomLat, randomLong)
	}
	return ""
}
