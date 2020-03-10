CREATE TABLE payments (
    id serial primary key,
    amount float not null,
    description varchar(512) not null,
    sender varchar(512) not null,
    datetime timestamp not null
);