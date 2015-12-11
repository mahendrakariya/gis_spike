`gen_insert_queries` is a script to generate insertion queries with random longitudes and latitudes within Bangalore.

**Usage:**
```
python gen_insert_queries.py 1000000 > queries.sql
```

Table can be created using the following query.
```
CREATE TABLE drivers(
     driver_id  INTEGER,
     geog       GEOGRAPHY(Point)
)
```
