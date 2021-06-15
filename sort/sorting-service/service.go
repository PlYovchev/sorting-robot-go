package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/plyovchev/go-at-ocado/sort/gen"
)

func newSortingService() gen.SortingRobotServer {
	return &sortingService{}
}

type sortingService struct {
	items        []*gen.Item
	selectedItem *gen.Item
	cubbiesItems map[string]*gen.Item
}

func (s *sortingService) LoadItems(context context.Context, request *gen.LoadItemsRequest) (*gen.LoadItemsResponse, error) {
	itemsToLoad := request.GetItems()
	if itemsToLoad == nil {
		return nil, errors.New("no items to load")
	}
	s.items = itemsToLoad
	return &gen.LoadItemsResponse{}, nil
}

func (s *sortingService) MoveItem(context context.Context, req *gen.MoveItemRequest) (*gen.MoveItemResponse, error) {
	if s.selectedItem == nil {
		return nil, errors.New("no item is selected")
	}

	if req.Cubby == nil {
		return nil, errors.New("no cubby specified")
	}

	cubbyId := req.Cubby.Id

	if s.cubbiesItems == nil {
		s.cubbiesItems = make(map[string]*gen.Item)
	}

	s.cubbiesItems[cubbyId] = s.selectedItem
	s.selectedItem = nil
	return &gen.MoveItemResponse{}, nil
}

func (s *sortingService) SelectItem(context.Context, *gen.SelectItemRequest) (*gen.SelectItemResponse, error) {
	if s.selectedItem != nil {
		return nil, errors.New("item already selected")
	}

	if len(s.items) == 0 {
		return nil, errors.New("no items in the main cubby")
	}

	fmt.Println(s.items)
	selectedItemAtIndex := rand.Intn(len(s.items))
	s.selectedItem = s.items[selectedItemAtIndex]
	s.items = deleteItemAtIndex(s.items, selectedItemAtIndex)

	return &gen.SelectItemResponse{Item: s.selectedItem}, nil
}

func deleteItemAtIndex(items []*gen.Item, index int) []*gen.Item {
	copy(items[index:], items[index+1:])
	items[len(items)-1] = nil
	return items[:len(items)-1]
}
