CREATE TABLE orders (
    order_id INT NOT NULL AUTO_INCREMENT,
    due_date TIMESTAMP NOT NULL,
    completed INT NOT NULL,
    PRIMARY KEY (order_id)
);

CREATE TABLE order_details (
    order_detail_id INT NOT NULL AUTO_INCREMENT,
    order_id INT,
    order_length INT NOT NULL,
    order_width INT NOT NULL,
    order_height INT NOT NULL,
    order_weight INT NOT NULL,
    PRIMARY KEY(order_detail_id),
    FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE consignee_details (
    consignee_detail_id INT NOT NULL AUTO_INCREMENT,
    order_id INT,
    consignee_name varchar(255) NOT NULL,
    consignee_number varchar(20) NOT NULL,
    consignee_country varchar(255) NOT NULL,
    consignee_address text NOT NULL,
    consignee_postal varchar(10) NOT NULL,
    consignee_state text NOT NULL,
    consignee_city text NOT NULL,
    consignee_province text NOT NULL,
    consignee_email varchar(255) NOT NULL,
    PRIMARY KEY(consignee_detail_id),
    FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE pickup_details (
    pickup_detail_id INT NOT NULL AUTO_INCREMENT,
    order_id INT,
    pickup_contact_name varchar(255) NOT NULL,
    pickup_contact_number varchar(20) NOT NULL,
    pickup_country varchar(255) NOT NULL,
    pickup_address text NOT NULL,
    pickup_postal varchar(10) NOT NULL,
    pickup_state text NOT NULL,
    pickup_city text NOT NULL,
    pickup_province text NOT NULL,
    PRIMARY KEY(pickup_detail_id),
    FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);