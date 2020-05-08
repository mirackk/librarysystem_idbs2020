create table if not exists Type(
	user_type int not null primary key,
    updatebook bool default false,
    adduser bool default false,
    borrowbook bool default false
);

create table if not exists User(
	user_id int not null auto_increment primary key,
    user_type int not null,
    name varchar(64) not null,
    foreign key(user_type)references Type(user_type)
);

create table if not exists Book(
	book_id int not null auto_increment primary key,
    title varchar(128) not null,
    author varchar(64) not null,
    ISBN varchar(64) not null,
    amount int not null,
    explanation varchar(128) 
);

create table if not exists Record(
    book_id int not null,
    user_id int not null,
    record_id int not null auto_increment,
    return_date date,
    borrow_date date not null,
    ddl date not null,
    etdtimes int not null,
    primary key(record_id, user_id, book_id),
    foreign key(book_id)references Book(book_id),
    foreign key(user_id)references User(user_id)
);
