create table if not exists genre
(
    id int not null,
    `name` varchar(20) unique not null,
    primary key (id)
);

create table if not exists book
(
    id int auto_increment not null,
    `name` varchar(100) unique not null,
    price float default 0 not null,
    genre int not null,
    amount int default 0 not null,
    primary key (id),
    foreign key (genre) references genre(id) on delete cascade
);