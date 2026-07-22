package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/devi/booklet/internal/testutil"
	"gorm.io/gorm"
)

var (
	testDB        *gorm.DB
	testContainer *testutil.PostgresContainer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := testutil.SetupPostgresContainer(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup postgres container: %v\n", err)
		os.Exit(1)
	}
	testContainer = container

	db, err := testutil.NewTestDB(container)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create test db: %v\n", err)
		if termErr := container.Terminate(ctx); termErr != nil {
			fmt.Fprintf(os.Stderr, "terminate container: %v\n", termErr)
		}
		os.Exit(1)
	}
	testDB = db

	code := m.Run()

	if termErr := container.Terminate(ctx); termErr != nil {
		fmt.Fprintf(os.Stderr, "terminate container: %v\n", termErr)
	}
	os.Exit(code)
}
