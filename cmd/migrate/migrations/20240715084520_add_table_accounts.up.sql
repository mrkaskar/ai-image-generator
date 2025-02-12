create table if not exists accounts (
    id serial primary key,
    user_id uuid references auth.users,
    username text not null,
    created_at timestamp not null default now()
);
