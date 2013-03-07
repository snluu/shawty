alter table `shawties` add column `CreatorIP` varchar(45) not null default '127.0.0.1' after `Url`;

create index idx_creator_ip on shawties(CreatorIP);
create index idx_created_on on shawties(CreatedOn);