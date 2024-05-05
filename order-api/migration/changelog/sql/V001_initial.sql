CREATE TABLE "order" (
	order_id serial4 NOT NULL,
	customer_name VARCHAR(100) NOT NULL,
	order_date TIMESTAMP WITH TIME ZONE NOT NULL,
	CONSTRAINT order_pk PRIMARY KEY (order_id)
);

CREATE TABLE "item" (
	item_id serial4 NOT NULL,
	order_id int4 NOT NULL,
	item_name VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
	CONSTRAINT item_pk PRIMARY KEY (item_id)
);

ALTER TABLE item ADD CONSTRAINT item_order_fk FOREIGN KEY (order_id) REFERENCES "order"(order_id);
