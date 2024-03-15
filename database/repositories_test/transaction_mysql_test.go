package repositories_test

import (
	"context"
	"crypto-challenge/config"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/docker/go-connections/nat"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TransactionMySqlIntTestSuite struct {
	suite.Suite
	terminateMySqlContainer *func()
	db                      *sql.DB
}

func (ts *TransactionMySqlIntTestSuite) SetupSuite() {
	dotenvFilePath, err := filepath.Abs(filepath.Join("..", "..", ".env"))
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.GetAppConfig(dotenvFilePath)

	mySqlC, terminateMySqlC, ctxMySqlC := setupMySqlContainer(cfg)

	ts.terminateMySqlContainer = terminateMySqlC

	endpoint, err := mySqlC.Endpoint(*ctxMySqlC, "")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.Database.User, cfg.Database.Password,
		endpoint, cfg.Database.DbName)

	ts.T().Log(dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		ts.T().Fatal("Failed opening connection to MySQL. ", err)
	}

	if err := db.Ping(); err != nil {
		ts.T().Fatal("Failed pinging MySQL. ", err)
	}

	ts.db = db
}

func (ts *TransactionMySqlIntTestSuite) TearDownSuite() {
	(*ts.terminateMySqlContainer)()
	ts.db.Close()
}

func (ts *TransactionMySqlIntTestSuite) TestHello() {
	ts.Equal("Hello", "Herou")
}

func TestTransactionMySqlIntTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionMySqlIntTestSuite))
}

func setupMySqlContainer(cfg *config.AppConfig, dsn string) (testcontainers.Container, *func(), *context.Context) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mysql@sha256:eeabfa5cd6a2091bf35eb9eae6ae48aab8231fd760f5a61cd0129df454333b1d",
		ExposedPorts: []string{"3306/tcp"},
		// (`/usr/sbin/mysqld: ready for connections\.`).AsRegexp()
		WaitingFor: wait.ForSQL(nat.Port(cfg.Database.Port), "mysql", func(host string) {
			return dsn
		}),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../../.docker/sql/create-transactions-table.sql",
				ContainerFilePath: "/docker-entrypoint-initdb.d/create-transactions-table.sql",
				FileMode:          0o755,
			},
		},
		Env: map[string]string{
			"MYSQL_USER":                 cfg.Database.User,
			"MYSQL_PASSWORD":             cfg.Database.Password,
			"MYSQL_DATABASE":             cfg.Database.DbName,
			"MYSQL_RANDOM_ROOT_PASSWORD": "yes",
		},
	}

	mySqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal("Could not start MySQL container.", err)
	}

	terminateMySqlC := func() {
		if err := mySqlC.Terminate(ctx); err != nil {
			log.Fatal("Could not stop MySQL container.", err)
		}
	}

	return mySqlC, &terminateMySqlC, &ctx
}
