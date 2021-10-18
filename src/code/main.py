import sys
import requests
import zlib
from requests.models import Response
import yaml
import pymysql
import hashlib
import multiprocessing as mp


def main() -> int:
    # Missing some error handling
    file_res = get_file_if_new(
        "http://localhost:2000/maps.yaml.gz", get_latest_modified_date())
    if (file_res.status_code == 304):
        print("File not downloaded as no new changed has occured since last read. Exiting.")
        return -1
    # I declare a variable here to maintain readability.
    stores_as_yaml = zlib.decompress(file_res.content, 16+zlib.MAX_WBITS)
    stores = yaml.safe_load(stores_as_yaml)
    mp.set_start_method('spawn')
    for store_key in stores:
        p = mp.Process(target=update_map_to_store_if_modified,
                       args=(store_key, stores[store_key]['map']))
        p.start()
        p.join()
    # The multiprocessing should be handled in such a way that only if everything
    # is "all_good", would we update the latest modified date.
    all_good = True
    if (all_good):
        set_latest_modified_date(file_res.headers["last-modified"])
    return 0


def update_map_to_store_if_modified(store_key: str, store_map: str):
    latest_map_hash = get_store(store_key)[1]
    store_map_hash = hashlib.sha1(str.encode(store_map)).hexdigest()
    if store_map_hash != str(latest_map_hash):
        send_new_map_to_store(store_key, store_map)
        update_store_in_db(store_map_hash, store_key)


def send_new_map_to_store(store_id: str, store_map: str):
    # Missing mqtt implementation.
    # To see how I would handle this, please see the file:
    # ../src_ggo/main.go function sendNewMapToStore
    print(store_id)


def update_store_in_db(new_map_hash: str, store_id: str):
    conn = get_sql_connection()
    cur = conn.cursor()
    cur.execute(
        f'update stores set LatestMapHash = "{new_map_hash}" where StoreId = "{store_id}"')
    cur.close()


def get_sql_connection():
    return pymysql.connect(host='localhost', port=23312,
                           user='dev', passwd='dev', db='stores')


def get_store(store_id: str) -> tuple:
    conn = get_sql_connection()
    cur = conn.cursor()
    cur.execute(f'SELECT * FROM stores WHERE StoreId = "{store_id}"')
    res = cur.fetchone()
    cur.close()
    conn.close()
    return res


def get_file_if_new(url: str, lastly_modified: str) -> Response:
    return requests.get(url, headers={"If-Modified-Since": lastly_modified})


def get_latest_modified_date() -> str:
    f = open("./data/modified-date.txt", "r")
    date = f.read()
    f.close()
    return date


def set_latest_modified_date(date: str):
    f = open("./data/modified-date.txt", "w")
    f.write(date)
    f.close()


if __name__ == '__main__':
    sys.exit(main())
