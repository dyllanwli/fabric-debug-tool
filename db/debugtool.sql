use debugtool;
drop table if exists debugtool;

create table debugtool
(
    userid int primary key auto_increment,
    username varchar(50) not null,
    userpassword varchar(50) not null
    org varchar(100) not null
);

drop table if exists logsinfo;
create table logsinfo
(
	logid int primary key auto_increment,
	logname varchar(200) not null,
	logpath varchar(500) not null,
	saveflag int not null
);
insert into debugtool (userid,username,userpassword,org) value(0,'diya',1234,'org1');