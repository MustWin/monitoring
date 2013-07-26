create database monitoring;
create user monitoring with password '' superuser;
create table services (
  id serial primary key,
  name character varying(127) not null,
  url character varying(255) not null,
  status text not null default '{}',
  healthy bool default true,
  acked bool default false,
  updated_at timestamp,
  created_at timestamp);
create unique index on services (name);
