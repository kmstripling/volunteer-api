-- Active: 1668178322497@@127.0.0.1@3306@volunteer
create database volunteer;

create table event
(
    id            int auto_increment
        primary key,
    name          varchar(64) not null,
    startdatetime datetime    null,
    enddatetime   datetime    null
);

create table person
(
    id        int auto_increment
        primary key,
    firstname varchar(20) null,
    lastname  varchar(30) null
);



create table registration
(
    id           int auto_increment
        primary key,
    event_id     int                                  not null,
    volunteer_id int                                  not null,
    regdatetime  datetime default current_timestamp() not null
);


create table timeentry
(
    id              int auto_increment
        primary key,
    registration_id int                                    not null,
    time_in         datetime default current_timestamp()   not null,
    time_out        datetime default '0000-00-00 00:00:00' not null
);
