use debugtool;
drop table if exists debugtool;

create table users
(
    userid int primary key auto_increment,
    username varchar(50) not null,
    userpassword varchar(50) not null
    org varchar(100) not null
);

drop table if exists channel;
create table channel
(
	channelid int primary key auto_increment,
	log varchar(500) not null,
);
insert into users (userid,username,userpassword,org) value(0,'jack','1234','org2');
insert into channel (channelid,log) value(0,'exammple_channel||../artifacts/channel/mychannel.tx');
