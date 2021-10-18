package storesDB

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

func getMySQLDB() (db *sql.DB, err error) {
	return sql.Open("mysql", "dev:dev@tcp(127.0.0.1:23312)/stores")
}

type store struct {
	StoreID       string
	LatestMapHash string
}

func GetStore(storeId string) (store, error) {
	var emptyStore store
	db, err := getMySQLDB()
	defer db.Close()
	if err != nil {
		return emptyStore, err
	}

	res, err := db.Query("SELECT * FROM stores WHERE StoreId = ?", storeId)
	defer res.Close()

	if err != nil {
		return emptyStore, err
	}

	for res.Next() {

		var store store
		err := res.Scan(&store.StoreID, &store.LatestMapHash)

		if err != nil {
			return store, err
		}

		return store, nil
	}
	return emptyStore, errors.New("No such store found in database.")
}

func UpdateStore(newMapHash string, storeId string) error {
	db, err := getMySQLDB()
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Exec("update stores set LatestMapHash = ? where StoreId = ?", newMapHash, storeId)
	return err
}

// Ideally, you would create an UpdateOrCreateStore function, that handles when the store is not already in the database
