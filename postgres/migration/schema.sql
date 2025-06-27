create table if not exists "user"
(
    login           varchar primary key,
    password        varchar          not null,
    accrual_balance double precision not null default 0
);

create type accrual_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

create table if not exists "order"
(
    id             varchar primary key,
    user_login     varchar        not null,
    uploaded_at    timestamp      not null default current_timestamp,
    accrual_status accrual_status not null default 'NEW'::accrual_status,
    accrual_points double precision,
    foreign key (user_login) references "user" (login)
);

create table if not exists withdrawal
(
    id           serial primary key,
    order_number varchar          not null,
    user_login   varchar          not null,
    sum          double precision not null,
    processed_at timestamp        not null default current_timestamp,
    foreign key (user_login) references "user" (login)
);
