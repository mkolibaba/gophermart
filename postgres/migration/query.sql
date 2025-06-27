-- name: UserGet :one
select *
from "user"
where login = $1;

-- name: UserExists :one
select exists (select 1 from "user" where login = $1);

-- name: UserSave :exec
insert into "user" (login, password)
values ($1, $2);

-- name: UserUpdateAccrualBalance :exec
update "user"
set accrual_balance = $1
where login = $2;

-- name: UserGetForLoginAndPassword :one
select *
from "user"
where login = $1
  and password = $2;

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

-- name: OrderUpdateAccrualStatus :exec
update "order"
set accrual_status = $1
where id = $2;

-- name: WithdrawalGetAll :many
select w.*
from withdrawal w
         join "order" o on o.id = w.order_number
where o.user_login = $1;

-- name: WithdrawalSave :exec
insert into withdrawal (order_number, user_login, sum)
values ($1, $2, $3);

-- name: BalanceGet :one
select cast(sum(w.sum) as double precision) as "current", u.accrual_balance as withdrawn
from withdrawal w
         join "user" u on u.login = w.user_login
where u.login = $1
group by u.accrual_balance;
