package main

import (
	"fmt"
	"testing"
)

func TestGetFileIfNew(t *testing.T) {
	fileRes, _ := GetFileIfNew("http://localhost:2000/maps.yaml.gz", "Sat, 16 Oct 2021 13:28:00 GMT")
	fmt.Println(fileRes.Body)
	if fileRes.StatusCode != 200 {
		t.Errorf("got %q, wanted %q", fileRes.StatusCode, 200)
	}
}

func TestDecodeGzipToYaml(t *testing.T) {
	// I know, I should read this data from a local file or a string, but I don't to save time for the
	// demonstration.
	fileRes, _ := GetFileIfNew("http://localhost:2000/maps.yaml.gz", "Sat, 16 Oct 2021 13:28:00 GMT")
	decoded, err := DecodeGzipToYaml(fileRes)
	if err != nil {
		t.Errorf("Failed converting gzip")
	}

	if string(decoded[0:5]) != "store" {
		t.Errorf("got %q, wanted %q", string(decoded[0:5]), "dawdawd")
	}
}

func TestExtractStoreMapFromStores(t *testing.T) {
	fileRes, _ := GetFileIfNew("http://localhost:2000/maps.yaml.gz", "Sat, 16 Oct 2021 13:28:00 GMT")
	storeYaml, err := DecodeGzipToYaml(fileRes)
	if err != nil {
		t.Errorf("Failed converting gzip")
	}
	stores, _ := ExtractStoresFromYaml(storeYaml)
	store, _ := FindStore(stores, "store-1")
	if store.Map[0:5] != "/9j/4" {
		t.Errorf("got %q, wanted %q", store.Map[0:5], "/9j/4")
	}
}

func TestShaOneOfString(t *testing.T) {
	fileRes, _ := GetFileIfNew("http://localhost:2000/maps.yaml.gz", "Sat, 16 Oct 2021 13:28:00 GMT")
	storeYaml, err := DecodeGzipToYaml(fileRes)
	if err != nil {
		t.Errorf("Failed converting gzip")
	}
	stores, _ := ExtractStoresFromYaml(storeYaml)
	store, _ := FindStore(stores, "store-1")
	storeMapHash := shaOne(store.Map)
	if storeMapHash != "5ee913b0db674d7b549c49eb794e6bfea7780b90" {
		t.Errorf("got %q, wanted %q", storeMapHash, "5ee913b0db674d7b549c49eb794e6bfea7780b90")
	}
}

func TestGetLatestModifiedDate(t *testing.T) {
	date, _ := getLatestModifiedDate()
	if date != "Sat, 16 Oct 2021 13:28:00 GMT" {
		t.Errorf("got %q, wanted %q", date, "Sat, 16 Oct 2021 13:28:00 GMT")

	}
}

func TestSetLatestModifiedDate(t *testing.T) {
	newDate := "fri, 15 Oct 2021 13:28:00 GMT"
	setLatestModifiedDate(newDate)
	date, _ := getLatestModifiedDate()
	if date != newDate {
		t.Errorf("got %q, wanted %q", date, newDate)
	}
	setLatestModifiedDate("Sat, 16 Oct 2021 13:28:00 GMT")

}
