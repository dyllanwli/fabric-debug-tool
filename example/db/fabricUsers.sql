use debugtool;
drop table if exists fabricusers;
create table fabricusers
(
    userid int primary key auto_increment,
    username varchar(50) not null,
    userpassword varchar(50) not null,
    token varchar(255) not null,
    org varchar(100) not null,
    balance double 
);
drop table if exists logsinfo;
create table logsinfo
(
	logid int primary key auto_increment,
	logname varchar(200) not null,
	logpath varchar(500) not null,
	saveflag int not null
);
insert into fabricusers (userid,username,userpassword,token,org,balance) value(0,'diya',1234,1234567890,'org1',1000);