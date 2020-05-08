package main
import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	// mysql connector
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const (
	User     = "root"
	Password = "135421"
	DBName   = "library"
)

type Authority struct{
	update bool
	adduser bool
	extend bool
	borrow bool
}

type Bookinfo struct{
	title string
	author string
	ISBN string
}

type Book struct{
	id int
	title string
	author string
	ISBN string
	amount int
	explanation string
}

type Record struct{
	bookid int
	userid int
	recordid int
	ifreturn bool
	returndate time.Time
	borrowdate time.Time
	ddl time.Time
	etdtimes int
}

func ConnectDB() error{
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
	if err != nil {
		panic(err)
	}
	return nil
}

// creat database and tables
func Init() error {
	//res,err:=ioutil.ReadFile("create_table.sql")
	/*query1, err := ioutil.ReadFile("/home/mirack/ass3/librarysystem_idbs2020/sqlfiles/create_database.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query1)); err != nil {
		panic(err)
	}*/
	_,err:=db.Exec("create database if not exists library")
	_,err=db.Exec("use library")
	if err!=nil{
		fmt.Println("can not create database")
		return err
	}
	query, err := ioutil.ReadFile("../sqlfiles/create_table.sql")
	if err != nil {
		panic(err)
	}
	sqls:=string(query)
	for _, s := range strings.Split(sqls, ";") {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}
	return nil
}

var errnouser=errors.New("no this user")
//return if one can do the thing
//we can check one's authority to do one thing
//in interaction I will command for user's id
func Usertypejudge(user int)(Authority,error){
	query:=`select updatebook,adduser,borrowbook from User,Type where User.user_type=Type.user_type and user_id=?`
	//var row *sql.Row
	row := db.QueryRow(query,user)
	var ath Authority
	err:=row.Scan(&ath.update,&ath.adduser,&ath.borrow)
	if err==sql.ErrNoRows{
		//err1:=errors.New("no this user")
		return ath,errnouser
	}
	if err!=nil{
		//err1:=errors.New("no this user")
		return ath,errnouser
	}
	//fmt.Printf("update:%v,adduser:%v,borrow:%v\n",ath.update,ath.adduser,ath.borrow)
	return ath,nil
}

//check if one book is in the library
func Checkbook(isbn string) bool{
	var tmp int
	err:=db.QueryRow(`SELECT book_id FROM Book WHERE ISBN =?`,isbn).Scan(&tmp)
	if err==sql.ErrNoRows{
		fmt.Println("here1")
		return false
	}
	if err!=nil{
		fmt.Println("here2")
		return false
	}
	return true
}

//add a book into the library;check first
func AddBook(info Bookinfo) error{
	flg:=Checkbook(info.ISBN)
	if flg==false{
		fmt.Println("it is new book")
		_,err:=db.Exec("insert into Book(title,author,ISBN,amount) values(?,?,?,1)",info.title,info.author,info.ISBN)
		if err!=nil{
			return err
		}
		return nil
	}
	if flg==true{
		fmt.Println("it is old book")
		_,err:=db.Exec("update Book set amount=amount+1 where ISBN=?",info.ISBN)
		if err!=nil{
			return err
		}
		return nil
	}
	return nil
}

//remove a book
var errnosuchbook= errors.New("no such book;can not delete")
var bookout=errors.New("book has been borrowed out")
func Remove(isbn string)error{
	flg:=Checkbook(isbn)
	if flg==false{
		//var err = errors.New("no such book;can not delete")
		return errnosuchbook
	}
	_, err :=db.Exec(`update Book set amount=amount-1,explanation="one of the book is lost" where amount>0 and ISBN=?`,isbn)
	if err==sql.ErrNoRows{
		//err1:=errors.New("no this book,can not remove")
		return bookout
	}
	if err != nil {
		return err
	}
	return nil
}

//searchbook by 3 methods
var errnothisbook=errors.New("no this book")
func Searchbytitle(tl string)([]Book,error){
	rows,err:=db.Query("select book_id,title,author,ISBN,amount from Book where title = ?",tl)
	if err!=nil{
		return nil ,err
	}
	defer rows.Close()
	books:=[]Book{}
	flg:=0
	for rows.Next(){
		var tmp Book
		err:=rows.Scan(&tmp.id,&tmp.title,&tmp.author,&tmp.ISBN,&tmp.amount)
		if err==sql.ErrNoRows{
			//err1:=errors.New("no this book")
			return nil ,errnothisbook
		}
		if err!=nil{
			return nil,err
		}
		books=append(books,tmp)
		flg=1
	}
	if flg==0{
		fmt.Println("no this book")
		return nil ,errnothisbook
	}
	return books,nil
}

func Searchbyauthor(ar string)([]Book,error){
	rows,err:=db.Query("select book_id,title,author,ISBN,amount from Book where author = ?",ar)
	if err!=nil{
		return nil ,err
	}
	defer rows.Close()
	books:=[]Book{}
	flg:=0
	for rows.Next(){
		var tmp Book
		err:=rows.Scan(&tmp.id,&tmp.title,&tmp.author,&tmp.ISBN,&tmp.amount)
		if err==sql.ErrNoRows{
			//err1:=errors.New("no this book")
			return nil ,errnothisbook
		}
		if err!=nil{
			return nil,err
		}
		books=append(books,tmp)
		flg=1
	}
	if flg==0{
		fmt.Println("no this book")
		return nil ,errnothisbook
	}
	return books,nil
}

func SearchbyISBN(isbn string)([]Book,error){
	rows,err:=db.Query("select book_id,title,author,ISBN,amount from Book where ISBN = ?",isbn)
	if err!=nil{
		return nil ,err
	}
	defer rows.Close()
	books:=[]Book{}
	flg:=0
	for rows.Next(){
		var tmp Book
		err:=rows.Scan(&tmp.id,&tmp.title,&tmp.author,&tmp.ISBN,&tmp.amount)
		if err==sql.ErrNoRows{
			//err1:=errors.New("no this book")
			return nil ,errnothisbook
		}
		if err!=nil{
			return nil,err
		}
		books=append(books,tmp)
		flg=1
	}
	if flg==0{
		fmt.Println("no this book")
		return nil ,errnothisbook
	}
	return books,nil
}

//add user into the database(2 type,addministrator, student)
func Adduser(usertype int,name string)(int,error){
	res,err:=db.Exec("insert into User(user_type,name) values(?,?)",usertype,name)
	if err!=nil{
		return -1 ,err
	}
	user_id,err:=res.LastInsertId()
	if err!=nil{
		return -1 ,err
	}
	return int(user_id),err
}

//borrow book
//check if there is the book
//check if one can borrow
//add record
var errcantfind=errors.New("can find the book")
var errover3=errors.New("you have more than 3 overdue")
var errbooknoamount=errors.New("the book has no copy available")
func Borrow(userid,bookid int)(int ,error){
	var tmp int
	err:=db.QueryRow("select book_id from Book where book_id=?",bookid).Scan(&tmp)
	if err == sql.ErrNoRows {
		//err1:=errors.New("can find the book")
		return -1, errcantfind
	}
	if err != nil {
		return -1, err
	}
	//now check if one can borrow
	now:=time.Now()
	var overdue_cnt int
	err = db.QueryRow(`select count(*) from Record where user_id=? and return_date is null and ddl <?`,userid,now).Scan(&overdue_cnt)
	if err != nil {
		return -1, err
	}
	if overdue_cnt > 3 {
		//var err1=errors.New("you have more than 3 overdue")
		return -1, errover3
	}

	due :=now.Add(time.Hour*24*30)

	res,err:=db.Exec(`update Book set amount=amount-1 where book_id=? and amount>0`,bookid)
	if err!=nil{
		//err2:=errors.New("the book has no copy available")
		return -1,errbooknoamount
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if cnt == 0 {
		//var err1=errors.New("no book available to borrow")
		return -1, errbooknoamount
	}

	res,err=db.Exec(`insert into Record(book_id,user_id,borrow_date,ddl,etdtimes) values (?,?,?,?,0)`,userid, bookid, now, due)
	if err != nil {
		return -1, err
	}

	recordid, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(recordid),nil
}

//return one book
var errnothisrecord=errors.New("no this record")
var erralready=errors.New("already returned")
func Returnbook(recordid int)error{
	var flg int
	//var rt_date time.Time
	err:=db.QueryRow(`select ISNULL(return_date) from Record where record_id=? `,recordid).Scan(&flg)
	if err==sql.ErrNoRows{
		//err1:=errors.New("no this record")
		return errnothisrecord
	}
	if err!=nil{
		return err
	}
	//fmt.Println(rt_date)
	//var zerodate time.Time
	//rt_date, _= time.Parse("2006-01-02 15:04:05", string(getdate)+" 00:00:00")
	if flg==0{
		//err2:=errors.New("already returned")
		return erralready
	}

	var book_id int
	err = db.QueryRow("select book_id from Record where record_id=?", recordid).
		Scan(&book_id)
	if err != nil { // must be valid record ID
		return err
	}
	_,err=db.Exec("update Book set amount=amount+1 where book_id=?",book_id)
	if err!=nil{
		return err
	}

	now:=time.Now()
	_,err=db.Exec("update Record set return_date=? where record_id=?",now,recordid)
	if err!=nil{
		return err
	}
	return nil
}

//extend return time
//one question not solved:if one book overdue but one can still extend,should system allow him?
//how about the overdue rules?
var errextd3=errors.New("already extend equal or over 3 times")
var erroverdue=errors.New("already overdue")
func Extendddl(recordid int)error{
	var getdue string
	var due time.Time
	var flg,etdtimes int
	err:=db.QueryRow(`select ISNULL(return_date),ddl,etdtimes from Record where record_id=? `,recordid).Scan(&flg,&getdue,&etdtimes)
	if err==sql.ErrNoRows{
		//err1:=errors.New("no this record")
		return errnothisrecord
	}
	if err!=nil{
		return err
	}
	if flg==0{
		//err2:=errors.New("already returned")
		return erralready
	}
	if etdtimes>=3{
		//err3:=errors.New("already extend over 3 times")
		return errextd3
	}

	//get the string "getdue" into time.Time to "due"
	due,_=time.Parse("2006-01-02 15:04:05", getdue+" 00:00:00")
	now:=time.Now()
	if due.Before(now) {
		//var err4=errors.New("already overdue")
		return erroverdue
	}
	new_ddl:=due.Add(time.Hour*24*7)
	_,err=db.Exec("update Record set ddl=? , etdtimes=etdtimes+1 where record_id=?",new_ddl,recordid)
	if err!=nil{
		return err
	}
	return nil
}

//read record by user_id,we can deal with a lot about one student through this api
//I am not sure about the argc type yet
//I am not sure about the empty value for time.Time
//with this api we can 1)query the books a student has borrowed and not returned yet
//and 2)check the deadline of returning a borrowed book
type Scantype interface {
	Scan(a ...interface{}) error
}

func Getrecord1(row Scantype)(Record,error){
	var r Record
	var returndate string
	var borrowdate,ddl string
	err:=row.Scan(&r.bookid,&r.userid,&r.recordid,&returndate,&borrowdate,&ddl,&r.etdtimes)
	if err!=nil{
		return Record{},err
	}
	/*if returndate.Valid{
		r.returndate=returndate.Time
	}*/
	r.returndate,_=time.Parse("2006-01-02 15:04:05", returndate+" 00:00:00")
	r.borrowdate,_=time.Parse("2006-01-02 15:04:05", borrowdate+" 00:00:00")
	r.ddl,_=time.Parse("2006-01-02 15:04:05", ddl+" 00:00:00")
	return r,nil
}

func Getrecord2(row Scantype)(Record,error){
	var r Record
	var borrowdate,ddl string
	err:=row.Scan(&r.bookid,&r.userid,&r.recordid,&borrowdate,&ddl,&r.etdtimes)
	if err!=nil{
		return Record{},err
	}
	/*if returndate.Valid{
		r.returndate=returndate.Time
	}*/
	r.borrowdate,_=time.Parse("2006-01-02 15:04:05", borrowdate+" 00:00:00")
	r.ddl,_=time.Parse("2006-01-02 15:04:05", ddl+" 00:00:00")
	return r,nil
}

//show one's borrow recodes
//guess sql.Row might work
//we don't have to check user_id here because we check when login
func Showhistory(user_id int)([]Record,error){
	histroy:=[]Record{}
	query1:= `SELECT book_id, user_id,record_id,return_date, borrow_date,ddl,etdtimes from Record join User using(user_id) where return_date is not NULL and user_id=?`
	rows,err:=db.Query(query1,user_id)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next(){
		res,err:=Getrecord1(rows)
		if err != nil {
			return nil, err
		}
		histroy=append(histroy,res)
	}
	query2:= `SELECT book_id, user_id,record_id,borrow_date,ddl,etdtimes from Record join User using(user_id) where return_date is NULL and user_id=?`
	rows,err=db.Query(query2,user_id)
	for rows.Next(){
		res,err:=Getrecord2(rows)
		if err != nil {
			return nil, err
		}
		histroy=append(histroy,res)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return histroy,nil
}

//check if one has registered
var errnosuchuser=errors.New("no such user")
func Useridcheck(user_id int)(string,error){
	var name string
	err:=db.QueryRow("select name from User where user_id=?",user_id).Scan(&name)
	if err==sql.ErrNoRows{
		//err1:=errors.New("no such user")
		return "",errnosuchuser
	}
	if err!=nil{
		return "",err
	}
	return name,nil
}

func main() {
	//db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	err:=ConnectDB()
	if err!=nil{
		fmt.Println("connect error")
	}
	fmt.Println("succeed in connect")
	err=Init()
	if err!=nil{
		fmt.Println("init error")
	}
	fmt.Printf("Welcome to the Library Management System!")
	for true {//login layer
		fmt.Printf("\n\n")
		fmt.Println("if you want to quit please cin \"quit\"")
		fmt.Println("if you don't;cin anything to continue")
		qt:="go on"
		fmt.Scanln(&qt)
		if qt=="quit"{
			break
		}
		fmt.Println("please cin your user_id to login")
		var user_id int
		fmt.Scanln(&user_id)
		name,err:=Useridcheck(user_id)
		if err!=nil{
			fmt.Println(err)
			continue
		}
		for true {
			fmt.Printf("\n\n")
			ath, _ := Usertypejudge(user_id)
			var num int
			fmt.Printf("hello %v. welcome back\n",name)
			fmt.Printf("update:%v,adduser:%v,borrow:%v\n",ath.update,ath.adduser,ath.borrow)
			fmt.Println("please cin a number to choose what you want to do")
			fmt.Println("1.add book\n2.remove book\n3.add account\n4.query book\n5.borrow book")
			fmt.Println("6.show history\n7.check return deadline\n8.extend ddl\n9.return book")
			fmt.Println("if you want to relogin please cin \"10\"")
			fmt.Scanln(&num)
			if num == 10 {
				break
			}
			switch num {
			case 1:
				if ath.update == false {
					fmt.Println("you don't have authority to do")
					continue
				}
				fmt.Println("plz cin title,author,ISBN in order")
				var book Bookinfo
				//fmt.Scanln(&book.title, &book.author, &book.ISBN)
				scanner1 := bufio.NewScanner(os.Stdin)
				scanner1.Scan()
				book.title=scanner1.Text()
				scanner2 := bufio.NewScanner(os.Stdin)
				scanner2.Scan()
				book.author=scanner2.Text()
				scanner3 := bufio.NewScanner(os.Stdin)
				scanner3.Scan()
				book.ISBN=scanner3.Text()

				err := AddBook(book)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("success in add")
			case 2:
				if ath.update == false {
					fmt.Println("you don't have authority to do")
					continue
				}
				fmt.Println("plz cin ISBN")
				var ISBN string
				fmt.Scanln(&ISBN)
				err := Remove(ISBN)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("success in add")
			case 3:
				if ath.adduser == false {
					fmt.Println("you don't have authority to do")
					continue
				}
				fmt.Println("plz cin usertype,name in order")
				var usertype int
				var name string
				fmt.Scanln(&usertype)
				fmt.Scanln(&name)
				id, err := Adduser(usertype, name)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("success in add,the new user_id=?", id)
			case 4:
				fmt.Println("choose a query method:1.title 2.author 3.ISBN")
				var mth int
				fmt.Scanln(&mth)
				if mth == 1 {
					var name string
					fmt.Println("cin title:")
					fmt.Scanln(&name)
					book, err := Searchbytitle(name)
					if err != nil {
						fmt.Println("can not find!")
						continue
					}
					for _, i := range book {
						fmt.Println("book info:")
						fmt.Printf("title:%v,author:%v,ISBN:%v,amount:%v,explanation:%v\n", i.title, i.author, i.ISBN, i.amount, i.explanation)
					}
				}
				if mth == 2 {
					var name string
					fmt.Println("cin author:")
					fmt.Scanln(&name)
					book, err := Searchbyauthor(name)
					if err != nil {
						fmt.Println(err)
						continue
					}
					for _, i := range book {
						fmt.Println("book info:")
						fmt.Printf("title:%v,author:%v,ISBN:%v,amount:%v,explanation:%v\n", i.title, i.author, i.ISBN, i.amount, i.explanation)
					}
				}
				if mth == 3 {
					var name string
					fmt.Println("cin ISBN:")
					fmt.Scanln(&name)
					book, err := SearchbyISBN(name)
					if err != nil {
						fmt.Println(err)
						continue
					}
					for _, i := range book {
						fmt.Println("book info:")
						fmt.Printf("title:%v,author:%v,ISBN:%v,amount:%v,explanation:%v\n", i.title, i.author, i.ISBN, i.amount, i.explanation)
					}
				}
			case 5:
				fmt.Println("cin the book_id you want to borrow:")
				var book_id int
				fmt.Scanln(&book_id)
				record_id, err := Borrow(user_id, book_id)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("success in borrow.your record id is:%v\n", record_id)
			case 6:
				rcd, err := Showhistory(user_id)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, i := range rcd {
					fmt.Println("your borrow history:")
					fmt.Printf("bookid:%v\nuserid:%v\nrecordid:%v\nreturndate:%v\nborrowdate:%v\nddl:%v\netdtimes:%v\n", i.bookid, i.userid, i.recordid, i.returndate, i.borrowdate, i.ddl, i.etdtimes)
				}
			case 7:
				rcd, err := Showhistory(user_id)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("plz cin the book id you borrow the book:")
				var bid int
				fmt.Scanln(&bid)
				flg := 0
				for _, i := range rcd{
					if i.bookid == bid {
						fmt.Printf("the ddl of the book is:%v\n", i.ddl)
						flg = 1
					}
				}
				if flg == 0 {
					fmt.Println("you haven't borrow this book")
				}
			case 8:
				fmt.Println("cin your record id;if you don't know please check your borrow history")
				var rcdid int
				fmt.Scanln(&rcdid)
				err := Extendddl(rcdid)
				if err != nil {
					fmt.Println(err)
					continue
				}
			case 9:
				fmt.Println("cin your record id;if you don't know please check your borrow history")
				var rcdid int
				fmt.Scanln(&rcdid)
				err := Returnbook(rcdid)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
			fmt.Println("tap anything int to coninue")
			var a int
			fmt.Scanln(&a)

		}
	}
}
