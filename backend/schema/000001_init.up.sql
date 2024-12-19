
CREATE TABLE users (
    id serial not null unique,
    username varchar(200) not null unique,
    email varchar(200) not null unique,
    password_hash varchar(256) not null,
    wallet int not null,
    role varchar(200) not null,
    subscription boolean,
    end_of_sub varchar(200)
);


CREATE TABLE sessions (
    id serial not null unique,
    user_id serial references users (id) on delete cascade not null,
    refresh_token text not null unique,
    fingerprint varchar(256) not null,
    Ip varchar(256)
);

CREATE TABLE planned_manga (
    id serial not null unique,
    user_id serial references users (id) on delete cascade not null,
    site varchar(256) not null,
    manga_id int
);

CREATE TABLE archived_files (
    id serial not null unique,
    user_id serial references users (id) on delete cascade not null,
    name varchar(256) not null,
    url varchar(256) not null
);

CREATE TABLE ads (
    id serial not null unique,
    ad_type int,
    banner_url varchar(256) not null,
    position int
);

