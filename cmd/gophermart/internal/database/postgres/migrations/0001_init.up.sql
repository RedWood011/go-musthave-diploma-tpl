CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create table if not exists schema_migrations
(
    version bigint  not null,
    dirty   boolean not null,
    primary key (version)
    );


CREATE TABLE IF NOT EXISTS public.users (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  login VARCHAR(50) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  balance  DECIMAL (16,2),
  spend DECIMAL (16,2)
);

CREATE TABLE IF NOT EXISTS public.withdrawals (
 user_id uuid,
 order_number TEXT NOT NULL UNIQUE ,
 status TEXT DEFAULT 'NEW',
 processed_at TIMESTAMP,
 sum DECIMAL (16,2)
);

CREATE TABLE IF NOT EXISTS  public.orders (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY ,
  user_id uuid REFERENCES users(id) ON DELETE CASCADE ,
  number TEXT NOT NULL UNIQUE,
  status TEXT,
  uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  accrual DECIMAL (16,2)
);