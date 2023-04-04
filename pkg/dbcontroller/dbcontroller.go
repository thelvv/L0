package DBcontroller

import (
	models "L0/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSLMode  string
}

func (cfg *Config) ConnectDB() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBname, cfg.SSLMode)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

type DBcontroller struct {
	db *sqlx.DB
}

func NewDBcontroller(db *sqlx.DB) *DBcontroller {
	return &DBcontroller{db}
}

func (DC *DBcontroller) GetDelivery(order_uid string, delivery *models.Delivery) (err error) {
	row := DC.db.QueryRowx(`SELECT name, phone, zip, city, address, region, email
	FROM delivery WHERE order_uid = $1;`, order_uid)

	err = row.StructScan(delivery)
	if err != nil {
		return err
	}
	return nil
}

func (DC *DBcontroller) GetPayment(order_uid string, payment *models.Payment) (err error) {
	row := DC.db.QueryRowx(`SELECT transaction, request_id, currency, provider, amount,
 	payment_dt, bank, delivery_cost, goods_total, custom_fee
	FROM payment WHERE order_uid = $1;`, order_uid)

	err = row.StructScan(payment)
	if err != nil {
		return err
	}
	return nil
}

func (DC *DBcontroller) GetItems(order_uid string, items *models.Items) (err error) {
	rows, err := DC.db.Queryx(`SELECT chrt_id, track_number, price, rid, name,
    sale, size, total_price, nm_id, brand, status
	FROM items WHERE order_uid = $1;`, order_uid)

	if err != nil {
		return err
	}

	var item models.Item
	for rows.Next() {
		err = rows.StructScan(&item)
		if err != nil {
			return err
		}
		*items = append(*items, item)
	}

	return nil
}

func (dbCon *DBcontroller) GetOrder(order_uid string) (models.Order, error) {
	row := dbCon.db.QueryRowx(`SELECT order_uid, track_number, entry, locale,
    internal_signature, customer_id, delivery_service, shardkey, sm_id,
    date_created, oof_shard
    FROM orders WHERE order_uid = $1;`, order_uid)

	var order models.Order
	err := row.StructScan(&order)
	if err != nil {
		return models.Order{}, err
	}

	err = dbCon.GetDelivery(order_uid, &order.Delivery)
	if err != nil {
		return models.Order{}, err
	}

	err = dbCon.GetPayment(order_uid, &order.Payment)
	if err != nil {
		return models.Order{}, err
	}

	err = dbCon.GetItems(order_uid, &order.Items)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (DC *DBcontroller) GetAllOrders() (map[string]models.Order, error) {
	rows, err := DC.db.Queryx(`SELECT order_uid FROM orders;`)
	if err != nil {
		return nil, err
	}

	orders := make(map[string]models.Order)
	var order_uid string
	for rows.Next() {
		err := rows.Scan(&order_uid)
		if err != nil {
			return nil, err
		}
		ord, err := DC.GetOrder(order_uid)
		if err != nil {
			return nil, err
		}

		orders[order_uid] = ord
	}

	return orders, nil
}

func (DC *DBcontroller) InsertDelivery(order_uid string, delivery *models.Delivery) (err error) {
	query := fmt.Sprintf(`INSERT INTO delivery(
			order_uid, name, phone, zip, city, address, region, email)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8);`)

	_, err = DC.db.Exec(
		query,
		order_uid,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)

	if err != nil {
		return err
	}
	return nil
}

func (DC *DBcontroller) InsertPayment(order_uid string, pay *models.Payment) (err error) {
	query := fmt.Sprintf(`INSERT INTO payment(
			order_uid, transaction, request_id, currency, provider, amount, 
            payment_dt, bank, delivery_cost, goods_total, custom_fee)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`)

	_, err = DC.db.Exec(
		query, order_uid,
		pay.Transaction,
		pay.Request_id,
		pay.Currency,
		pay.Provider,
		pay.Amount,
		pay.Payment_dt,
		pay.Bank,
		pay.Delivery_cost,
		pay.Goods_total,
		pay.Custom_fee,
	)

	if err != nil {
		return err
	}
	return nil
}

func (DC *DBcontroller) InsertItem(order_uid string, itm *models.Item) (err error) {
	query := fmt.Sprintf(`INSERT INTO items(
			order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`)

	_, err = DC.db.Exec(
		query, order_uid,
		itm.Chrt_id,
		itm.Track_number,
		itm.Price,
		itm.Rid,
		itm.Name,
		itm.Sale,
		itm.Size,
		itm.Total_price,
		itm.Nm_id,
		itm.Brand,
		itm.Status,
	)
	if err != nil {
		return err
	}
	return nil
}

func (DC *DBcontroller) InsertOrder(ord *models.Order) (err error) {
	query := fmt.Sprintf(`INSERT INTO orders(
			order_uid, track_number, entry, locale, internal_signature, 
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
			VALUES(:order_uid, :track_number, :entry, :locale, :internal_signature, :customer_id, 
			:delivery_service, :shardkey, :sm_id, :date_created, :oof_shard);`)
	_, err = DC.db.NamedExec(query, ord)
	if err != nil {
		return err
	}

	err = DC.InsertDelivery(ord.Order_uid, &ord.Delivery)
	if err != nil {
		return err
	}

	err = DC.InsertPayment(ord.Order_uid, &ord.Payment)
	if err != nil {
		return err
	}

	for _, itm := range ord.Items {
		err = DC.InsertItem(ord.Order_uid, &itm)
		if err != nil {
			return err
		}
	}

	return nil
}
