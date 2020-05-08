insert into Type(user_type,updatebook,adduser,borrowbook)
values
(1,true,true,true),
(2,false,false,true),
(3,false,false,false);

insert into User(user_type,name)
values
(1,"adminone"),
(1,"admintwo"),
(1,"adminthree"),
(2,"mirack"),
(2,"meow"),
(2,"rabbit"),
(3,"badguy");

insert into Book(title,author,ISBN,amount,explanation)
values
("The art of translation","Liu, Xiaodong","9787544657624",2,NULL),
("Speaking for ourselves","Nikolajeva, Maria","9787544656511",3,NULL),
("Chinese literature","Epiphanius Wilson","9787508539270",4,NULL),
("Toward a network theory of acculturation","Chi, Ruobing","9787313186430",1,NULL),
("Anti-mimesis from Plato to Hitchcock","Cohen, Tom","9787521301083",0,NULL),
("On deconstruction : theory and criticism after structuralism","Culler, Jonathan","9787521301045",5,NULL),
("American dream, American nightmare : fiction since 1960","Hume, Kathryn","9787521301106",8,NULL),
("The antinomies of realism","Jameson, Fredric","9787521301137",3,NULL),
("Fictions of authority : women writers and narrative voice","Lanser, Susan Sniader","9787521301090",1,NULL),
("Madame Bovary","FLaubert, Gustave","9787506297455",2,NULL),
("Key concepts in contemporary literature","Padley, Steve","9787544646086",4,NULL),
("Key concepts in postcolonial literature","Wisker, Gina","9787544646062",5,NULL);

##not return 
insert into Record(book_id,user_id,borrow_date,ddl,etdtimes)
values
(1,4,"2020-05-06","2020-06-13",1),
(2,4,"2020-04-06","2020-05-13",1),
(3,4,"2020-03-23","2020-05-14",3),
(6,6,"2020-03-30","2020-05-06",1),
(10,5,"2020-05-06","2020-06-06",0);

insert into Record(book_id,user_id,return_date,borrow_date,ddl,etdtimes)
values
(5,4,"2020-03-01","2020-02-06","2020-03-13",1),
(12,4,"2019-09-01","2020-08-06","2020-09-06",0),
(7,6,"2020-03-15","2020-02-06","2020-03-20",2);

