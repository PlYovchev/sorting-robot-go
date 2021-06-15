package main

import (
	"testing"

	"github.com/plyovchev/go-at-ocado/sort/gen"
)

func TestLoadItems(t *testing.T) {
	sortingService := newSortingService()
	request := &gen.LoadItemsRequest{
		Items: nil,
	}

	_, err := sortingService.LoadItems(nil, request)
	if err == nil {
		t.Fatal("expected error")
	}
}
