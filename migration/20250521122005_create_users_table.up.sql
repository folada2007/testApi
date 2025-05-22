CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username text not null default 'unknow',
    surname text not null default 'unknow',
    age INTEGER not null,
    nationality text not null default 'unknow',
    gender text not null default 'unknow'
);