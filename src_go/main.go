package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	storesDB "github.com/TobiasFP/cm-q3-watcher/storesDB"

	yaml "gopkg.in/yaml.v3"
)

func main() {
	dateOfLastFileRead, err := getLatestModifiedDate()
	if err != nil {
		log.Fatalf("Problem reading the last date the file was modified. " + err.Error())
	}

	fileRes, err := GetFileIfNew("http://localhost:2000/maps.yaml.gz", dateOfLastFileRead)
	if err != nil {
		log.Fatalf("Problem downloading file: " + err.Error())
	}

	if fileRes.StatusCode == 304 {
		log.Fatalf("File not downloaded as no new changed has occured since last read. Exiting.")
	}

	StoresAsBytes, err := DecodeGzipToYaml(fileRes)
	if err != nil {
		log.Fatalf("Problem decoding gzipped file: " + err.Error() + " pro")
	}

	stores, err := ExtractStoresFromYaml(StoresAsBytes)
	for storeID, storeMap := range stores {
		go updateMapToStoreIfModified(storeID, storeMap)
	}
	// Waiting for the goroutines to finish.
	// Should implement channels to
	// wait for all goroutines to be done,
	// in order to verify all has gone well,
	// so we could "setLatestModifiedDate" but this
	// is quick and dirty
	time.Sleep(5 * time.Second)
	err = setLatestModifiedDate(fileRes.Header.Get("Last-Modified"))
	if err != nil {
		log.Fatalf("Problem setting the file modified date: " + err.Error())
	}
	fmt.Println("Done!")
}

func updateMapToStoreIfModified(storeId string, storeMap Store) {
	store, err := storesDB.GetStore(storeId)
	if err != nil {
		log.Fatalf("Problem finding store in db " + err.Error())
	}

	mapHashFromYaml := shaOne(storeMap.Map)
	if mapHashFromYaml != store.LatestMapHash {
		sendNewMapToStore(storeId, storeMap)
		storesDB.UpdateStore(mapHashFromYaml, storeId)
	}
}

func sendNewMapToStore(storeId string, store Store) {
	broker := "localhost"
	port := 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("tobias")
	opts.SetUsername("tobias")
	opts.SetPassword("notasecurepassword")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := client.Publish("topic/new-map", 0, false, store.Map); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Updated map in " + storeId)
}

func GetFileIfNew(url string, lastlyModified string) (ress *http.Response, err error) {
	res := &http.Response{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return res, err
	}

	req.Header.Set("If-Modified-Since", lastlyModified)
	res, err = client.Do(req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func DecodeGzipToYaml(gzFileRes *http.Response) ([]byte, error) {
	var body []byte
	gzFileAsBytes, err := ioutil.ReadAll(gzFileRes.Body)
	if err != nil {
		return body, err
	}
	gzReader, err := gzip.NewReader(bytes.NewBuffer(gzFileAsBytes))
	if err != nil {
		return body, err
	}
	return ioutil.ReadAll(gzReader)
}

type stores map[string]Store

type Store struct {
	Map string
}

func ExtractStoresFromYaml(yamlRes []byte) (stores, error) {
	stores := make(stores)
	err := yaml.Unmarshal(yamlRes, &stores)
	return stores, err
}

func FindStore(stores stores, storeIdNeedle string) (Store, error) {
	for storeId, storeMap := range stores {
		if storeId == storeIdNeedle {
			return storeMap, nil
		}
	}
	return Store{}, errors.New("No store with that identifier found.")
}

func shaOne(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func getLatestModifiedDate() (string, error) {
	date, err := os.ReadFile("./data/modified-date.txt")
	return string(date), err
}

func setLatestModifiedDate(date string) error {
	dateAsBytes := []byte(date)
	return os.WriteFile("./data/modified-date.txt", dateAsBytes, 0644)
}
