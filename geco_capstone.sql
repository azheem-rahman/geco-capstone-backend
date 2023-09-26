CREATE DATABASE capstonedb

CREATE TABLE accounts (
    account_id INT NOT NULL AUTO_INCREMENT,
    email varchar(255) NOT NULL,
    password text NOT NULL,
    account_type ENUM ('admin','partner_malaysia','partner_indonesia'),
    PRIMARY KEY (account_id),
);

CREATE TABLE accounts_details (
    detail_id INT NOT NULL AUTO_INCREMENT,
    account_id INT,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    PRIMARY KEY (detail_id),
    CONSTRAINT fk_user
        FOREIGN KEY (account_id)
        REFERENCES accounts(account_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE orders (
    order_id INT NOT NULL AUTO_INCREMENT,
    account_id INT,
    order_length INT NOT NULL,
    order_width INT NOT NULL,
    order_height INT NOT NULL,
    order_weight INT NOT NULL,
    consignee_name varchar(255) NOT NULL,
    consignee_number varchar(20) NOT NULL,
    consignee_country varchar(255) NOT NULL,
    consignee_address text NOT NULL,
    consignee_postal varchar(10) NOT NULL,
    consignee_state text NOT NULL,
    consignee_city text NOT NULL,
    consignee_province text NOT NULL,
    consignee_email varchar(255) NOT NULL,
    pickup_contact_name varchar(255) NOT NULL,
    pickup_contact_number varchar(20) NOT NULL,
    pickup_country varchar(255) NOT NULL,
    pickup_address text NOT NULL,
    pickup_postal varchar(10) NOT NULL,
    pickup_state text NOT NULL,
    pickup_city text NOT NULL,
    pickup_province text NOT NULL,
    due_date TIMESTAMP NOT NULL,
    completed INT NOT NULL,
    PRIMARY KEY (order_id),
    CONSTRAINT fk_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(account_id)
);

CREATE TABLE items (
	item_id INT NOT NULL AUTO_INCREMENT,
    order_id INT,
    item_description text NOT NULL,
    item_category varchar(255) NOT NULL,
    item_product_id varchar(255) NOT NULL,
    item_sku varchar(255) NOT NULL,
    item_quantity INT NOT NULL,
    item_price_value decimal(10,2) NOT NULL,
    item_price_currency varchar(5) NOT NULL,
    PRIMARY KEY(item_id)
);

CREATE TABLE order_items (
	order_id int,
    item_id int,
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    FOREIGN KEY (item_id) REFERENCES items(item_id)
);