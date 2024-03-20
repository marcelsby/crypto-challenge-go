package helpers

import (
	"context"
	"crypto-challenge/config"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupMySqlContainer(cfg *config.AppConfig, migrationsFolderPath string) (testcontainers.Container, *func(), *context.Context) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mysql@sha256:eeabfa5cd6a2091bf35eb9eae6ae48aab8231fd760f5a61cd0129df454333b1d",
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor: wait.ForSQL("3306/tcp", "mysql", func(host string, port nat.Port) string {
			dbCfg := mysql.Config{
				User:   cfg.Database.User,
				Passwd: cfg.Database.Password,
				Net:    port.Proto(),
				Addr:   fmt.Sprintf("%s:%d", host, port.Int()),
				DBName: cfg.Database.DbName,
			}

			return dbCfg.FormatDSN()
		}).WithPollInterval(time.Millisecond * 500),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      filepath.Join(migrationsFolderPath, "create-transactions-table.sql"),
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
