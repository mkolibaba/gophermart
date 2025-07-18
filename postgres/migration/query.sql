-- name: UserGet :one
select *
from "user"
where login = $1;

-- name: UserSave :exec
insert into "user" (login, password)
values ($1, $2);

-- name: UserAddAccrualBalance :exec
update "user"
set accrual_balance = accrual_balance + $1
where login = $2;

-- name: OrderGet :one
select *
from "order"
where id = $1;

-- name: OrderGetAll :many
select *
from "order"
where user_login = $1;

-- name: OrderSave :exec
insert into "order" (id, user_login)
values ($1, $2);

-- name: OrderGetWithNonFinalAccrualStatus :many
select *
from "order"
where accrual_status not in ('INVALID', 'PROCESSED')
limit $1;

-- name: OrderUpdateAccrual :exec
update "order"
set accrual_status = $1,
    accrual_points = $2
where id = $3;

-- name: WithdrawalGetAll :many
select w.*
from withdrawal w
where w.user_login = $1;

-- name: WithdrawalSave :exec
insert into withdrawal (order_number, user_login, sum)
values ($1, $2, $3);

-- name: BalanceGet :one
select cast(coalesce((select sum(w.sum) from withdrawal w where w.user_login = u.login),
                     0) as double precision) as withdrawn,
       u.accrual_balance                     as "current"
from "user" u
where u.login = $1;
