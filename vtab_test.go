package vtab

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/mattn/go-sqlite3"
)

func TestVTab(t *testing.T) {
	sql.Register("sqlite3_vtab", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.CreateModule("test_vtab", &module{intarray: []int{
				10, 20, 30, 40, 50,
			}})
		},
	})

	db, err := sql.Open("sqlite3_vtab", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from test_vtab")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var val int
		rows.Scan(&val)
		fmt.Println(val)
	}

}
