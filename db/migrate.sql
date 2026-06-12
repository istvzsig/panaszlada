create table if not exists public.reports (
    id text primary key,
    tracking_code text unique not null,

    title text not null,
    description text not null,
    category text not null,

    latitude double precision,
    longitude double precision,

    status text not null,
    created_at timestamptz not null
);