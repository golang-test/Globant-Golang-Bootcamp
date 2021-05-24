DROP TABLE IF EXISTS book;
DROP TABLE IF EXISTS genre;

CREATE TABLE IF NOT EXISTS genre (
	id SERIAL PRIMARY KEY,
	name VARCHAR(45) NOT NULL);

CREATE TABLE IF NOT EXISTS book (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) UNIQUE NOT NULL,
	price double precision NOT NULL,
	genre_id INT NOT NULL REFERENCES genre (Id),
	amount INT NOT NULL);

INSERT INTO genre (name) VALUES ('Adventure');
INSERT INTO genre (name) VALUES ('Classics');
INSERT INTO genre (name) VALUES ('Fantasy');

INSERT INTO book (name, price, genre_id, amount) VALUES ('book1', '10', '1', '50');
INSERT INTO book (name, price, genre_id, amount) VALUES ('book2', '11', '2', '1');
INSERT INTO book (name, price, genre_id, amount) VALUES ('book3', '20.6', '3', '3');
INSERT INTO book (name, price, genre_id, amount) VALUES ('book4', '25', '1', '4');
INSERT INTO book (name, price, genre_id, amount) VALUES ('book5', '30.5', '2', '2');