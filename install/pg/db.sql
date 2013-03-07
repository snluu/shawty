create sequence shawties_id_seq;

create table shawties(
	ID bigint primary key not null default nextval('shawties_id_seq'),
	Rand char(1) not null,
	Hits bigint not null default 0,
	Url varchar(2048) not null,
	CreatedOn bigint not null,
	CreatorIP varchar(45) not null default '127.0.0.1'
);

create index idx_rand on shawties(ID, Rand);
create unique index idx_url on shawties(Url);

create index idx_creator_ip on shawties(CreatorIP);
create index idx_created_on on shawties(CreatedOn);

grant all on shawties to dev;
grant all on shawties_id_seq to dev;