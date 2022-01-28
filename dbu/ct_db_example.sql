create database webuser;
create database dndcharacter;
\c webuser;
create user usersvc with password 'usersvcpass123';
alter database webuser owner to usersvc;
\c dndcharacter;
create user dndcharacter with password 'dndcharpass123';
alter database dndcharacter owner to dndcharacter;

