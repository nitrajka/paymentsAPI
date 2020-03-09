CREATE TABLE payments (
    id serial primary key,
    amount numeric not null,
    description varchar(512),
    sender varchar(512) not null
);