package repository

import (
	models "L0/models"
	"L0/pkg/cache"
	DBcontroller "L0/pkg/dbcontroller"
	"fmt"
)

type Repo interface {
	GetOrder(order_uid string) (models.Order, error)
	GetAllOrders() (map[string]models.Order, error)
	InsertOrder(order *models.Order) error
}

type Repository struct {
	DBcontroller    DBcontroller.DBcontroller
	cacheController cache.Cache
}

func NewRepo() *Repository {
	cfg := DBcontroller.Config{
		Host:     "127.0.0.1",
		Port:     "5432",
		Username: "al",
		Password: "123456",
		// убрать пароль
		DBname:  "l0",
		SSLMode: "disable",
	}

	db, err := cfg.ConnectDB()

	if err != nil {
		fmt.Println("[LOG]: error in connecting DB:", err.Error())
		return nil
	}

	dbController := DBcontroller.NewDBcontroller(db)

	Repo := &Repository{
		DBcontroller:    *dbController,
		cacheController: *cache.New(),
	}
	fmt.Println("Repo =", Repo)
	Repo.updateCache(nil)

	return Repo
}

func (Repo *Repository) GetOrder(order_uid string) (models.Order, error) {
	order, err := Repo.cacheController.GetOrder(order_uid)

	if err != nil {
		order, err = Repo.DBcontroller.GetOrder(order_uid)

		if err != nil {
			return models.Order{}, err
		}

		go Repo.cacheController.InsertOrder(&order)
		return order, nil
	}

	return order, nil
}

func (Repo *Repository) updateCache(orders map[string]models.Order) (err error) {
	if len(orders) == 0 {
		allOrders, err := Repo.DBcontroller.GetAllOrders()

		if err != nil {
			return err
		}

		Repo.cacheController.UpdateCache(allOrders)
	}

	Repo.cacheController.UpdateCache(orders)

	return nil
}

func (Repo *Repository) GetAllOrders() (map[string]models.Order, error) {
	orders, err := Repo.cacheController.GetAllOrders()

	if err != nil {
		orders, err = Repo.DBcontroller.GetAllOrders()

		if err != nil {
			return nil, err
		}

		go Repo.updateCache(orders)
		return orders, nil
	}

	return orders, nil
}

func (Repo *Repository) InsertOrder(order *models.Order) error {
	err := Repo.cacheController.InsertOrder(order)
	if err != nil {
		return err
	}

	err = Repo.DBcontroller.InsertOrder(order)
	if err != nil {
		return err
	}

	return nil
}
