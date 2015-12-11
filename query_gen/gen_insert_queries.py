import random
import sys


def gen_name():
    alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
    name = []
    for i in range(15):
        name.append(alphabet[random.randint(0, 51)])
    name = "".join(name)
    return "'" + name + "'"

def get_random_lat():
    lat1 = 12.813196
    lat2 = 13.055798
    return random.uniform(lat1, lat2)

def get_random_long():
    long1 = 77.474313
    long2 = 77.767158
    return random.uniform(long1, long2)

def generate_insert_queries(n):
    queries = []
    values = []
    for i in range(n):
        latitude = get_random_lat()
        longitude = get_random_long()
        q = "(%d, 'POINT(%s %s)')" %(i+1, str(latitude), str(longitude))
        values.append(q)
    queries.append("INSERT INTO drivers VALUES\n" + ",\n".join(values) + ";")
    return queries


if len(sys.argv) == 1:
    print "Usage: python gen_insert_queries.py 1000000 > queries.sql"
    exit(1)

n = int(sys.argv[1])
print "\n".join(generate_insert_queries(n))
