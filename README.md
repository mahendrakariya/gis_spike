# PostGIS Spike


* Create database named `cabspike` in Postgres.
* Enable PostGIS by running the command
```
CREATE EXTENSION postgis;
```
* Create a table named drivers and insert a few records. [See this](https://github.com/mahendrakariya/postgis_spike/tree/master/query_gen) for details.
* Start `pgbouncer`.
```
cd pgbouncer
pgbouncer -d pgbouncer.ini
cd ..
```
* Run the code.
```
go run src/read_write.go
```
