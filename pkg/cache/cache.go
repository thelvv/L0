package cache

import (
	models "L0/models"
	"errors"
	"sync"
)

type Cache struct {
	cache map[string]models.Order
	mutex *sync.Mutex
}

func New() *Cache {
	cacheNew := Cache{
		cache: make(map[string]models.Order),
		mutex: new(sync.Mutex),
	}
	return &cacheNew
}

func (c *Cache) GetOrder(order_uid string) (models.Order, error) {
	_, exists := c.cache[order_uid]

	if exists {
		return c.cache[order_uid], nil
	}

	return models.Order{}, errors.New("error! Order does not exist")
}

func (c *Cache) GetAllOrders() (map[string]models.Order, error) {
	if len(c.cache) == 0 {
		return nil, errors.New("error! Cache is empty")
	}

	return c.cache, nil

}

func (c *Cache) InsertOrder(order *models.Order) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cache == nil {
		c.cache = make(map[string]models.Order)
	}

	c.cache[order.Order_uid] = *order

	if c.cache[order.Order_uid].Order_uid != order.Order_uid {
		return errors.New("error! Can't insert to cache")
	}

	return nil
}

func (c *Cache) UpdateCache(orders map[string]models.Order) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = orders

	if len(c.cache) == 0 && len(orders) != 0 {
		return errors.New("error! Can't update cache")
	}

	return nil
}
