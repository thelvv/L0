CREATE TABLE orders (
    order_uid               varchar(255) PRIMARY KEY,
    track_number            varchar(255),
    entry                   varchar(255),
    locale                  varchar(255),
    internal_signature      varchar(255),
    customer_id             varchar(255),
    delivery_service        varchar(255),
    shardkey                varchar(255),
    sm_id                   int,
    oof_shard               varchar(255),
    date_created            timestamp,
    recording_date          timestamp
);

CREATE TABLE delivery (
    order_uid               varchar(255) references orders (order_uid) on delete cascade not null unique,
    name                    varchar(255),
    phone                   varchar(255),
    zip                     varchar(255),
    city                    varchar(255),
    address                 varchar(255),
    region                  varchar(255),
    email                   varchar(255),
    recording_date          timestamp
);

CREATE TABLE payment (
     order_uid               varchar(255) references orders (order_uid) on delete cascade not null unique,
     transaction             varchar(255),
     request_id              varchar(255),
     currency                varchar(255),
     provider                varchar(255),
     amount                  int,
     payment_dt              int,
     bank                    varchar(255),
     delivery_cost           int,
     goods_total             int,
     custom_fee              int,
     recording_date          timestamp
);

CREATE TABLE items (
    order_uid               varchar(255) references orders (order_uid) on delete cascade not null,
    chrt_id                 int,
    track_number            varchar(255),
    price                   int,
    rid                     varchar(255),
    name                    varchar(255),
    sale                    int,
    size                    varchar(255),
    total_price             int,
    nm_id                   int,
    brand                   varchar(255),
    status                  int,
    recording_date          timestamp
);

INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', '', 'test', 'meest', '9', 99, '2022-07-22 17:42:39.5555+03:0', '1');

INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
VALUES ('b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809', 'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com');

INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('b563feb7b2b84b6test', 9934930, 'WBILMTESTTRACK', 453, 'ab4219087a764ae0btest', 'Mascaras', 30, '0', 317, 2389212, 'Vivienne Sabo', 202);

INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('b563feb7b2b84b6test', '', 'USD', 'wbpay', '123', 1817, 1637907727, 'alpha', 1500, 317, 0);
