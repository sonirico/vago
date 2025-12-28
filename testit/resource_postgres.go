package testit

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/ory/dockertest/v3"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sonirico/vago/lol"
)

func NewPostgresResource(
	migrationsPath string,
	logger lol.Logger,
	envarURL string,
	migrateFunc MigrateFunc,
	envFunc SetEnvFunc,
) *Resource {
	log := logger.WithField("resource", "postgres")
	if envarURL == "" {
		envarURL = "BROCK_POSTGRES_URL"
	}
	return &Resource{
		RunOptions: &dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "16.2",
			Env: []string{
				"POSTGRES_PASSWORD=secret",
				"POSTGRES_USER=user_name",
				"POSTGRES_DB=dbname",
				"listen_addresses = '*'",
			},
			Hostname:     "postgres",
			ExposedPorts: []string{"5432/tcp"},
		},
		HostConfig: &docker.HostConfig{
			PublishAllPorts: true,
			Mounts: []docker.HostMount{
				{
					Type:   "tmpfs",
					Target: "/var/lib/postgresql/data",
				},
			}},
		RetryFunc: func(dockerhost string, resource *dockertest.Resource) retryFunc {
			databaseUrl := fmt.Sprintf(
				"postgres://user_name:secret@%s:%s/dbname?sslmode=disable",
				dockerhost,
				resource.GetPort("5432/tcp"),
			)

			log.Infof("Connecting to database on url: %s", databaseUrl)

			return func() error {
				db, err := sql.Open("postgres", databaseUrl)
				if err != nil {
					return err
				}
				defer db.Close()
				return db.Ping()
			}

		},
		MigrateFunc: migrateFunc,
		SetEnvFunc:  envFunc,
	}
}

func CopyDirectory(scrDir, dest string) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		_, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			continue
		case os.ModeSymlink:
			continue
		default:
			if err := CopyFile(sourcePath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyFile(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}
