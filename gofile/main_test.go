package main

import (
	"fmt"
	"testing"
)

func TestUsertypejudge(t *testing.T) {
	//err1:=errors.New("no this user")
	ConnectDB()
	Init()
	userid:=[]struct{
		id int
		err error
	}{
		{1,nil},
		{2,nil},
		{3,nil},
		{100,errnouser},
	}
	for _,i:=range userid{
		_,err:=Usertypejudge(i.id)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
	}
}

func TestCheckbook(t *testing.T) {
	ConnectDB()
	Init()
	isbns:=[]struct{
		s string
		flg bool
	}{
		{"9787544657624",true},
		{"9787508539270",true},
		{"123",false},
		{"1234567891011",false},
	}
	for _,i:=range isbns{
		tmp:=Checkbook(i.s)
		if tmp!=i.flg{
			t.Errorf("should be %v,get %v",i.flg,tmp)
		}
	}
}

func TestAddBook(t *testing.T) {
	ConnectDB()
	Init()
	info:=[]struct{
		book Bookinfo
		err error
	}{
		{Bookinfo{"Number theory. Volume I","Cohen, Henri","9787519255299"},nil},
		{Bookinfo{"Topological invariants of plane curves and caustics","Arnold, V. I. (Vladimir Igorevich)","9787040517057"},nil},
		{Bookinfo{"The art of translation","Liu, Xiaodong","9787544657624"},nil},
		{Bookinfo{"Speaking for ourselves","Nikolajeva, Maria","9787544656511"},nil},
	}
	for _,i:=range info{
		err:=AddBook(i.book)
		if err!=nil{
			t.Errorf("should be nil,get %v",err)
		}
	}
}

func TestRemove(t *testing.T) {
	ConnectDB()
	Init()
	//err1:=errors.New("no this book,can not remove")
	isbns:=[]struct{
		isbn string
		err error
	}{
		{"9787521301137",nil},
		{"9787521301090",nil},
		{"1234567891011",errnosuchbook},
	}
	for _,i:=range isbns{
		err:=Remove(i.isbn)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
	}
}

//as for 3 kinds of searches are similar,I just test two of them
func TestSearchbytitle(t *testing.T) {
	ConnectDB()
	Init()
	//err1:=errors.New("no this book")
	cs:= []struct {
		title string
		isbn string
		err error
	}{
		{"The art of translation","9787544657624",nil},
		{"Speaking for ourselves","9787544656511",nil},
		{"whatever","1234567891011",errnothisbook},
	}
	for _,i:=range cs{
		book,err:=Searchbytitle(i.title)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
		for _,j:=range book{
			if j.ISBN!=i.isbn{
				fmt.Println("wrong match when searching books by title")
			}
		}
	}
}

func TestSearchbyauthor(t *testing.T) {
	ConnectDB()
	Init()
	//err1:=errors.New("no this book")
	cs:= []struct {
		author string
		isbn string
		err error
	}{
		{"Liu, Xiaodong","9787544657624",nil},
		{"Nikolajeva, Maria","9787544656511",nil},
		{"whatever","1234567891011",errnothisbook},
	}
	for _,i:=range cs{
		book,err:=Searchbyauthor(i.author)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
		for _,j:=range book{
			if j.ISBN!=i.isbn{
				fmt.Println("wrong match when searching books by author")
			}
		}
	}
}

func TestAdduser(t *testing.T) {
	ConnectDB()
	Init()
	cs:=[]struct{
		usertype int
		name string
	}{
		{1,"aadmin"},
		{2,"cjy"},
		{3,"zsk"},
	}
	for _,i:=range cs{
		_,err:=Adduser(i.usertype,i.name)
		if err!=nil{
			t.Errorf("should be %v,get %v",nil,err)
		}
	}
}

func TestBorrow(t *testing.T) {
	ConnectDB()
	Init()
	//err1:=errors.New("can find the book")
	//err2:=errors.New("the book has no copy available")
	cs:=[]struct{
		user_id int
		book_id int
		err error
	}{
		{4,6,nil},
		{5,6,nil},
		{6,15,errcantfind},
		{4,5,errbooknoamount},
	}
	for _,i:=range cs{
		_,err:=Borrow(i.user_id,i.book_id)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
	}
}

func TestReturnbook(t *testing.T) {
	ConnectDB()
	Init()
	//err1:=errors.New("no this record")
	//err2:=errors.New("already returned")
	cs:=[]struct{
		record_id int
		err error
	}{
		{2,nil},
		{5,nil},
		{100,errnothisrecord},
		{7,erralready},
	}
	for _,i:=range cs{
		err:=Returnbook(i.record_id)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
	}
}

func TestExtendddl(t *testing.T) {
	ConnectDB()
	Init()
	//err1:=errors.New("no this record")
	//err2:=errors.New("already returned")
	//err3:=errors.New("already extend over 3 times")
	//err4:=errors.New("already overdue")
	cs:=[]struct{
		record_id int
		err error
	}{
		{1,nil},
		{2,nil},
		{100,errnothisrecord},
		{7,erralready},
		{3,errextd3},
		{4,erroverdue},
	}
	for _,i:=range cs{
		err:=Extendddl(i.record_id)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
	}
}

//use this to test both showhistory and getrecord
func TestShowhistory(t *testing.T) {
	ConnectDB()
	Init()
	cs:=[]struct{
		user_id int
		err error
	}{
		{4,nil},
		{5,nil},
		{6,nil},
	}
	for _,i:=range cs{
		_,err:=Showhistory(i.user_id)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
	}
}

func TestUseridcheck(t *testing.T) {
	ConnectDB()
	Init()
	//err1:=errors.New("no such user")
	cs:=[]struct{
		id int
		err error
	}{
		{1,nil},
		{4,nil},
		{100,errnosuchuser},
	}
	for _,i:=range cs{
		_,err:=Useridcheck(i.id)
		if err!=i.err{
			t.Errorf("should be %v,get %v",i.err,err)
		}
	}
}
