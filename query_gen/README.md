`gen_insert_queries` is a script to generate insertion queries with random longitudes and latitudes within Bangalore.

**Usage:**
```
python gen_insert_queries.py 1000000 > queries.sql
```

Table can be created using the following query.
**For postgres:**
```
CREATE TABLE drivers(
     driver_id  INTEGER,
     geog       GEOGRAPHY(Point)
)
```

**For memsql:**
```
CREATE TABLE drivers(
     driver_id  INT PRIMARY KEY,
     geog       GEOGRAPHYPOINT
)
```
