create table shawties(
	ID bigint primary key not null auto_increment,
	Rand char not null,
	Hits bigint not null default 0,
	Url varchar(2048) not null,
	CreatedOn bigint not null
);

create index idx_rand on shawties(ID, Rand);
create unique index idx_url on shawties(Url(78));