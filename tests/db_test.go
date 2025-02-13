package tests

import (
	"context"
	"os"
	"testing"

	"github.com/prorok210/AvitoShop/internal/db"
)

func Test_connectDB(t *testing.T) {
	if os.Getenv("DATABASE_USER") == "" ||
		os.Getenv("DATABASE_PASSWORD") == "" ||
		os.Getenv("DATABASE_HOST") == "" ||
		os.Getenv("DATABASE_PORT") == "" ||
		os.Getenv("DATABASE_NAME") == "" {
		t.Skip("Skipping test because environment variables are not set")
	}

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
		db.DBConn.Close()
	}
}
