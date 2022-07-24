CREATE DATABASE wb_db;    

\c wb_db; 

CREATE TABLE wb_data (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "order_uid" VARCHAR(50),
    "track_number" VARCHAR(50),
    "entry" VARCHAR(50),
    "delivery" JSON,
    "payment" JSON,
    "items" JSON,
    "locale" VARCHAR(10),
    "internal_signature" VARCHAR(50),
    "customer_id" VARCHAR(50),
    "delivery_service" VARCHAR(50),
    "shardkey" VARCHAR(50),
    "sm_id" INT,
    "date_created" TIMESTAMP,
    "oof_shard" VARCHAR(50)
);
