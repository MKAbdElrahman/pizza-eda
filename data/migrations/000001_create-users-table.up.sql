CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(20),
    created DATETIME NOT NULL,
    CONSTRAINT users_uc_email UNIQUE (email)
);


