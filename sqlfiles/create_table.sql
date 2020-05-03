create table user(
	user_id int not null primary key,
    user_type int not null,
    name varchar(64) not null,
    foreign key(user_type)references type(user_type)
);

create table type(
	user_type int not null primary key,
    addbook bool default false,
    removebook bool default false,
    adduser bool default false,
    extend bool default true
);

create table book(
	book_id int not null primary key,
    title varchar(64) not null,
    author varchar(64) not null,
    ISBN varchar(64) not null,
    amount int not null,
    explanation varchar(128) 
);

create table record(
    book_id int not null,
    user_id int not null,
    record_id int not null,
    return_date date,
    borrow_date date not null,
    ddl date not null,
    primary key(record_id, user_id, book_id),
    foreign key(book_id)references book(book_id),
    foreign key(user_id)references user(user_id)
);