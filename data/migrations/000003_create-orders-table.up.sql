CREATE TABLE orders (
    order_id VARCHAR(255) PRIMARY KEY,
    user_id INT,
    sauce VARCHAR(255),
    cheese VARCHAR(255),
    main_topping VARCHAR(255),
    extra_topping VARCHAR(255),
    status VARCHAR(50),
    timestamp DATETIME
);
