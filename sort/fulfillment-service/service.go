package main

import (
	"context"
	"errors"
	"log"

	"github.com/plyovchev/sorting-robot-go/sort/gen"
	"github.com/preslavmihaylov/ordertocubby"
)

func newFulfillmentService(sortRobotClient gen.SortingRobotClient) gen.FulfillmentServer {
	return &fulfillmentService{
		sortRobotClient: sortRobotClient,
		cubbyToOrder:    make(map[string]string),
		itemToCubbies:   make(map[string][]*gen.Cubby),
	}
}

type fulfillmentService struct {
	// orderToCubby map[string]string
	sortRobotClient gen.SortingRobotClient
	cubbyToOrder    map[string]string
	itemToCubbies   map[string][]*gen.Cubby
	preparedOrders  []*gen.PreparedOrder
}

func (fs *fulfillmentService) LoadOrders(ctx context.Context, in *gen.LoadOrdersRequest) (*gen.CompleteResponse, error) {
	fs.mapOrdersToCubbies(in.Orders)

	err := fs.executeOrders()
	if err != nil {
		return nil, err
	}

	return &gen.CompleteResponse{Orders: fs.preparedOrders}, nil
}

func (fs *fulfillmentService) mapOrdersToCubbies(ordersToMap []*gen.Order) {
	fs.preparedOrders = make([]*gen.PreparedOrder, 0, len(ordersToMap))
	for _, order := range ordersToMap {
		cubbyId := ""
		counter := uint32(1)
		for cubbyId == "" || fs.cubbyToOrder[cubbyId] != "" {
			cubbyId = ordertocubby.Map(order.Id, counter, 10)
			counter++
		}

		fs.cubbyToOrder[cubbyId] = order.Id
		log.Printf("Order %s is assigned to cubby %s", order.Id, cubbyId)

		cubby := &gen.Cubby{Id: cubbyId}
		fs.populateItemToCubbiesForOrder(order, cubby)

		fs.preparedOrders = append(fs.preparedOrders, &gen.PreparedOrder{Order: order, Cubby: cubby})
	}
}

func (fs *fulfillmentService) populateItemToCubbiesForOrder(order *gen.Order, cubby *gen.Cubby) {
	for _, item := range order.Items {
		if fs.itemToCubbies[item.Code] == nil {
			fs.itemToCubbies[item.Code] = []*gen.Cubby{cubby}
		} else {
			fs.itemToCubbies[item.Code] = append(fs.itemToCubbies[item.Code], cubby)
		}
	}
}

func (fs *fulfillmentService) executeOrders() error {
	for len(fs.itemToCubbies) > 0 {
		pickItemResponse, err := fs.sortRobotClient.PickItem(context.Background(), &gen.Empty{})
		if err != nil {
			log.Fatalf("failed to pick item: %v", err)
			return errors.New("error occured when tried to pick item")
		}
		log.Printf("Picked item is %s", pickItemResponse.Item.Label)

		cubbies := fs.itemToCubbies[pickItemResponse.Item.Code]
		cubbyForSelectedItem := cubbies[len(cubbies)-1]
		_, err = fs.sortRobotClient.PlaceInCubby(context.Background(), &gen.PlaceInCubbyRequest{Cubby: cubbyForSelectedItem})
		if err != nil {
			return errors.New("error occured when tried to place item in cubby")
		}

		cubbies = cubbies[:len(cubbies)-1]
		if len(cubbies) == 0 {
			delete(fs.itemToCubbies, pickItemResponse.Item.Code)
		} else {
			fs.itemToCubbies[pickItemResponse.Item.Code] = cubbies
		}
	}

	return nil
}
