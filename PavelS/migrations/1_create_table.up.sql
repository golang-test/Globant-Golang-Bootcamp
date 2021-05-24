CREATE TABLE books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    price FLOAT,
    genre INT,
    amount INT
);
CREATE TABLE genres (
    id INT AUTO_INCREMENT PRIMARY KEY,
    value INT,
    name VARCHAR(100)
);
INSERT INTO books (name, price, genre, amount) VALUES ('The Three Musketeers', 10.44, 1, 5);
INSERT INTO books (name, price, genre, amount) VALUES ('War and Peace', 7.30, 2, 0);
INSERT INTO books (name, price, genre, amount) VALUES ('Harry Potter', 19.99, 3, 10);
INSERT INTO genres (value, name) VALUES (1, 'Adventure');
INSERT INTO genres (value, name) VALUES (2, 'Classics');
INSERT INTO genres (value, name) VALUES (3, 'Fantasy');

