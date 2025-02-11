package tests

import (
	"context"
	"testing"

	"github.com/prorok210/AvitoShop/internal/db"
)

func Test_connectDB(t *testing.T) {
	if err := db.ConnectDB(); err != nil {
		t.Errorf("An error was expected but was not found")
	} else {
		if db.DBConn == nil {
			t.Errorf("DBConn is nil")
		}
		q, err := db.DBConn.Exec(context.Background(), "SELECT 1;")
		if err != nil || q.String() != "SELECT 1" {
			t.Errorf("Error executing request")
		}
	}
}
