# Бизнес логика работы с заказами.

## Требования

**POST /api/user/orders** — загрузка пользователем номера заказа для расчёта;
**GET /api/user/orders** — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
**GET /api/user/balance** — получение текущего баланса счёта баллов лояльности пользователя;
**POST /api/user/balance/withdraw** — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
**GET /api/user/withdrawals** — получение информации о выводе средств с накопительного счёта пользователем.

## Реализация

**Предварительный вид таблиц для работы с логикой**

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  login VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
  id BIGSERIAL PRIMARY KEY,
  number VARCHAR(255) UNIQUE NOT NULL,
  status INTEGER,
  accrual INTEGER,
  user_id BIGINT REFERENCES users(id),
  uploaded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_balance (
  id BIGSERIAL PRIMARY KEY,
  current INTEGER,
  withdraw INTEGER
  user_id BIGINT REFERENCES users(id),
);
```

**Статусы заказа**

_NEW_ — заказ загружен в систему, но не попал в обработку;
_PROCESSING_ — вознаграждение за заказ рассчитывается;
_INVALID_ — система расчёта вознаграждений отказала в расчёте;
_PROCESSED_ — данные по заказу проверены и информация о расчёте успешно получена.

**Этапы логики**

1. Создание заказа с первоначальным статусом "NEW" с асинхронным запуском тайм тикера(паттерн worker pool) для получение информации о текущем состояния статуса заказа. После получения обновляем таблицу orders(accrual),orders(status), user_balance(current)
