package main

import "sync"

type Column struct {
	Name     string
	Datatype string
}

type Row struct {
	Values map[string]interface{}
}

type ForeignKey struct {
	Column           string
	ReferencedTable  string
	ReferencedColumn string
}

type Table struct {
	mu          sync.Mutex
	Name        string
	Columns     []Column
	Rows        []Row
	PrimaryKey  string
	ForeignKeys []ForeignKey
}

func NewTable(name string, columns []Column, primaryKey string) *Table {
	return &Table{
		Name:       name,
		Columns:    columns,
		PrimaryKey: primaryKey,
	}
}

func (t *Table) InsertRow(row Row) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Rows = append(t.Rows, row)
}
func (t *Table) SelectRows(columnName string, value interface{}) []Row {
	t.mu.Lock()
	defer t.mu.Unlock()

	var result []Row
	for _, row := range t.Rows {
		if row.Values[columnName] == value {
			result = append(result, row)
		}
	}
	return result
}

func (t *Table) DeleteRows(columnName string, value interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	newRows := []Row{}
	for _, row := range t.Rows {
		if row.Values[columnName] != value {
			newRows = append(newRows, row)
		}
	}
	t.Rows = newRows
}

func (t *Table) AddForeignKey(fk ForeignKey) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ForeignKeys = append(t.ForeignKeys, fk)
}

type Database struct {
	mu     sync.Mutex
	Tables map[string]*Table
}

func NewDatabase() *Database {
	return &Database{
		Tables: make(map[string]*Table),
	}
}

func (d *Database) CreateTable(name string, columns []Column, primaryKey string) {
	d.mu.Lock()
	d.mu.Unlock()
	d.Tables[name] = NewTable(name, columns, primaryKey)
}
