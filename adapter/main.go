package main

import "fmt"

type Database interface {
	GetData()
}

type MySQLDb struct {
}

func (m *MySQLDb) GetDataFromMySQL() {
	fmt.Println("my sql")
}

type MySQLAdapter struct {
	mySQLDb *MySQLDb
}

func NewMySQLAdapter(db *MySQLDb) *MySQLAdapter {
	return &MySQLAdapter{
		mySQLDb: db,
	}
}

func (m *MySQLAdapter) GetData() {
	m.mySQLDb.GetDataFromMySQL()
}

func main() {
	x := NewMySQLAdapter(&MySQLDb{})
	x.GetData()
}
