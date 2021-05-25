create table BookStore (id int, name char(100), price int, genre  int, amount int);
create table Genres (id int, genre_name char(100));

insert into BookStore values (1, 'book', 100, 1, 3) ;

insert into Genres values (1, 'Adventure') ;
insert into Genres values (2, 'Classics') ;
insert into Genres values (3, 'Fantasy') ;