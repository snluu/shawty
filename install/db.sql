create table shawties(
	ID bigint primary key not null auto_increment,
	Rand char not null,
	Hits bigint not null default 0,
	Url varchar(2048) not null,
	CreatedOn bigint not null
);

create index idx_rand on shawties(ID, Rand);
create unique index idx_url on shawties(Url(78));

alter table `shawties` add column `CreatorIP` varchar(45) not null default '127.0.0.1' after `Url`;

create index idx_creator_ip on shawties(CreatorIP);
create index idx_created_on on shawties(CreatedOn);