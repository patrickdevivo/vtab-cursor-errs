package vtab

import (
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type module struct {
	intarray []int
}

type vtab struct {
	intarray []int
}

type vtabCursor struct {
	vTab  *vtab
	index int
}

func (m module) Create(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	err := c.DeclareVTab("CREATE TABLE x(test TEXT)")
	if err != nil {
		return nil, err
	}
	return &vtab{m.intarray}, nil
}

func (m module) EponymousOnlyModule() {}

func (m module) Connect(c *sqlite3.SQLiteConn, args []string) (sqlite3.VTab, error) {
	return m.Create(c, args)
}

func (m module) DestroyModule() {}

func (v *vtab) BestIndex(cst []sqlite3.InfoConstraint, ob []sqlite3.InfoOrderBy) (*sqlite3.IndexResult, error) {
	used := make([]bool, 0, len(cst))
	for range cst {
		used = append(used, false)
	}

	// TODO uncomment the below line to see a failure
	// return nil, fmt.Errorf("silly error, does appear and cause failure")

	return &sqlite3.IndexResult{
		Used:           used,
		IdxNum:         0,
		IdxStr:         "test-index",
		AlreadyOrdered: true,
		EstimatedCost:  100,
		EstimatedRows:  200,
	}, nil
}

func (v *vtab) Disconnect() error {
	return nil
}

func (v *vtab) Destroy() error {
	return nil
}

func (v *vtab) Open() (sqlite3.VTabCursor, error) {
	return &vtabCursor{v, 0}, nil
}

func (vc *vtabCursor) Close() error {
	return nil
}

func (vc *vtabCursor) Filter(idxNum int, idxStr string, vals []interface{}) error {
	vc.index = 0
	return nil
}

func (vc *vtabCursor) Next() error {
	vc.index++
	if vc.index > 2 {
		// the following error does seem to halt execution of the cursor, but does not indicate
		// an error in the execution of the query
		return fmt.Errorf("silly error, doesn't seem to bubble out")
	}
	return nil
}

func (vc *vtabCursor) EOF() bool {
	return vc.index >= len(vc.vTab.intarray)
}

func (vc *vtabCursor) Column(c *sqlite3.SQLiteContext, col int) error {
	if col != 0 {
		return fmt.Errorf("column index out of bounds: %d", col)
	}
	c.ResultInt(vc.vTab.intarray[vc.index])
	return nil
}

func (vc *vtabCursor) Rowid() (int64, error) {
	return int64(vc.index), nil
}
