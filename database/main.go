package main

import (
	"fmt"
	"sync"
)

// Column struct - Represents a column in a table
type Column struct {
	Name     string
	DataType string // e.g., INT, VARCHAR, etc.
}

// Row struct - Represents a single row in a table
type Row struct {
	Values map[string]interface{} // Column Name -> Value
}

// ForeignKey struct - Represents a foreign key relationship
type ForeignKey struct {
	Column           string // Column in current table
	ReferencedTable  string // Table being referenced
	ReferencedColumn string // Column in referenced table
}

// Table struct - Represents a table with columns, rows, and constraints
type Table struct {
	mu          sync.RWMutex
	Name        string
	Columns     []Column
	Rows        []Row
	PrimaryKey  string
	ForeignKeys []ForeignKey // Foreign key relationships
}

// NewTable - Create a new table with given columns
func NewTable(name string, columns []Column, primaryKey string) *Table {
	return &Table{
		Name:       name,
		Columns:    columns,
		PrimaryKey: primaryKey,
	}
}

// InsertRow - Inserts a new row into the table
func (t *Table) InsertRow(row Row) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Rows = append(t.Rows, row)
}

// SelectRows - Fetches rows based on a simple condition
func (t *Table) SelectRows(columnName string, value interface{}) []Row {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var result []Row
	for _, row := range t.Rows {
		if row.Values[columnName] == value {
			result = append(result, row)
		}
	}
	return result
}

// DeleteRows - Deletes rows based on a condition
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

// AddForeignKey - Adds a foreign key to the table
func (t *Table) AddForeignKey(fk ForeignKey) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ForeignKeys = append(t.ForeignKeys, fk)
}

// Database struct - Manages multiple tables
type Database struct {
	mu     sync.RWMutex
	Tables map[string]*Table
}

// NewDatabase - Creates a new empty database
func NewDatabase() *Database {
	return &Database{
		Tables: make(map[string]*Table),
	}
}

// CreateTable - Adds a new table to the database
func (db *Database) CreateTable(name string, columns []Column, primaryKey string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Tables[name] = NewTable(name, columns, primaryKey)
}

// InsertIntoTable - Inserts a row into a table by name
func (db *Database) InsertIntoTable(tableName string, row Row) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	table, exists := db.Tables[tableName]
	if exists {
		table.InsertRow(row)
	}
}

// SelectFromTable - Selects rows from a table by name
func (db *Database) SelectFromTable(tableName, columnName string, value interface{}) []Row {
	db.mu.RLock()
	defer db.mu.RUnlock()
	table, exists := db.Tables[tableName]
	if exists {
		return table.SelectRows(columnName, value)
	}
	return nil
}

// Sample usage
func main() {
	// Create a new database
	db := NewDatabase()

	// Define columns for the 'users' table
	userColumns := []Column{
		{"id", "INT"},
		{"name", "VARCHAR"},
		{"email", "VARCHAR"},
	}

	// Create the 'users' table with 'id' as the primary key
	db.CreateTable("users", userColumns, "id")

	// Insert a row into the 'users' table
	row := Row{Values: map[string]interface{}{
		"id":    1,
		"name":  "Alice",
		"email": "alice@example.com",
	}}
	db.InsertIntoTable("users", row)

	// Select rows where name is 'Alice'
	result := db.SelectFromTable("users", "name", "Alice")
	for _, row := range result {
		fmt.Println("User:", row.Values)
	}

	// Add foreign key example (e.g., orders table referencing users)
	orderColumns := []Column{
		{"order_id", "INT"},
		{"user_id", "INT"}, // Foreign key to users table
		{"amount", "FLOAT"},
	}
	db.CreateTable("orders", orderColumns, "order_id")
	ordersTable := db.Tables["orders"]
	ordersTable.AddForeignKey(ForeignKey{
		Column:           "user_id",
		ReferencedTable:  "users",
		ReferencedColumn: "id",
	})

	// Insert into orders
	db.InsertIntoTable("orders", Row{Values: map[string]interface{}{
		"order_id": 101, "user_id": 1, "amount": 250.0,
	}})

	// Query for orders of user 1
	orders := db.SelectFromTable("orders", "user_id", 1)
	fmt.Println("Orders for user 1:", orders)
}
