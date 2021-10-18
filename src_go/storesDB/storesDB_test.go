package storesDB

import (
	"testing"
)

func TestGetStore(t *testing.T) {
	store, _ := GetStore("store-1")
	if store.LatestMapHash != "4ee212b0db674d7b549c49eb453dcbfea7780b64" {
		t.Errorf("got %q, wanted %q", store.LatestMapHash, "4ee212b0db674d7b549c49eb453dcbfea7780b64")
	}
}

func TestUpdateStore(t *testing.T) {
	err := UpdateStore("mjello", "store-2")
	if err != nil {
		t.Errorf("Could not update, due to: %q", err.Error())
	}
}
