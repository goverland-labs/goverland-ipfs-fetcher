create table ipfs_data
(
    ipfs_id    text not null
        primary key,
    created_at timestamp with time zone,
    type       text,
    data       jsonb
);


