package repository_integration

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/redblood-pixel/pastebin/pkg/logger"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbURL string

func TestMain(m *testing.M) {

	ctx := context.Background()
	logger.Init("test")
	container, url := startPostgresContainer(ctx)
	defer container.Terminate(ctx)
	dbURL = url
	setUpContainer(ctx)

	m.Run()

}

func waitForDB(dbURL string, timeout time.Duration) error {
	start := time.Now()
	for time.Since(start) < timeout {
		db, err := sql.Open("postgres", dbURL)
		if err == nil {
			if err = db.Ping(); err == nil {
				db.Close()
				return nil
			}
			db.Close()
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("database not ready after %v", timeout)
}

func startPostgresContainer(ctx context.Context) (testcontainers.Container, string) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:14",
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		ExposedPorts: []string{"5432/tcp"},

		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithStartupTimeout(60 * time.Second).
			WithPollInterval(time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		slog.Error("failed to start container", "err", err.Error())
	}
	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")
	url := fmt.Sprintf("postgres://test:test@%s:%s/test?sslmode=disable", host, port.Port())
	fmt.Println(url)
	return postgresC, url
}

func setUpContainer(ctx context.Context) {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("failed to get path")
	}
	migrationsPath := filepath.Join(filepath.Dir(path), "..", "..", "..", "migrations")
	fmt.Println(migrationsPath)

	if err := waitForDB(dbURL, 30*time.Second); err != nil {
		fmt.Printf("Database not ready: %v", err)
		return
	}
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		fmt.Println(err)
	}
	if err = m.Up(); err != nil {
		fmt.Println(err)
	}
	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec(ctx, `
    -- Вставка пользователей
    INSERT INTO users (name, email, password_hashed) VALUES
    ('userr1', 'user1@example.com', '$2a$10$xJwL5v9zZ1hBp5UQ6YQZQe9v6ZG9Jt8Xe6ZG9Jt8Xe6ZG9Jt8Xe6ZG'),
    ('userr2', 'user2@example.com', '$2a$10$xJwL5v9zZ1hBp5UQ6YQZQe9v6ZG9Jt8Xe6ZG9Jt8Xe6ZG9Jt8Xe6ZG'),
    ('userr3', 'user3@example.com', '$2a$10$xJwL5v9zZ1hBp5UQ6YQZQe9v6ZG9Jt8Xe6ZG9Jt8Xe6ZG9Jt8Xe6ZG'),
    ('userr4', 'user4@example.com', '$2a$10$xJwL5v9zZ1hBp5UQ6YQZQe9v6ZG9Jt8Xe6ZG9Jt8Xe6ZG9Jt8Xe6ZG'),
	('userr5', 'user5@email11.com', '$2a$10$xJwL5v9zZ1hBp5UQ6YQZQe9v6ZG9Jt8Xe6ZG9Jt8Xe6ZG9Jt8Xe6ZG');

    -- Вставка паст
    INSERT INTO pastes (title, expires_at, visibility, user_id) VALUES
    ('Public paste 1', '2030-01-01 00:00:00', 'public', 1),
    ('Private paste 1', '2030-01-01 00:00:00', 'private', 1),
    ('Temporary paste', NOW() + INTERVAL '1 day', 'public', 2),
    ('Important notes', '2030-01-01 00:00:00', 'private', 2),
    ('Code snippet', NOW() + INTERVAL '7 days', 'public', 3),
    ('Personal diary', '2030-01-01 00:00:00', 'private', 3),
    ('Work documentation', NOW() + INTERVAL '30 days', 'public', 4),
    ('Secret project', '2030-01-01 00:00:00', 'private', 4),
    ('Public API docs', '2030-01-01 00:00:00', 'public', 1),
    ('Backup codes', NOW() + INTERVAL '14 days', 'private', 2),
    ('Tutorial', '2030-01-01 00:00:00', 'public', 3),
    ('Configuration', '2030-01-01 00:00:00', 'private', 4),
    ('Meeting notes', NOW() + INTERVAL '2 days', 'public', 1),
    ('Ideas', '2030-01-01 00:00:00', 'private', 2),
    ('Recipes', '2030-01-01 00:00:00', 'public', 3),
    ('Passwords', NOW() + INTERVAL '90 days', 'private', 4),
    ('Public FAQ', '2030-01-01 00:00:00', 'public', 1),
    ('Private thoughts', '2030-01-01 00:00:00', 'private', 2),
    ('Team docs', NOW() + INTERVAL '180 days', 'public', 3),
    ('Financial info', '2030-01-01 00:00:00', 'private', 4);
	
	INSERT INTO pastes (title, created_at, expires_at, visibility, user_id) VALUES
	('x1', NOW() - INTERVAL '10 minutes', '2030-01-01 00:00:00', 'public', 5),
	('x2', NOW() - INTERVAL '2 days', '2030-01-01 00:00:00', 'public', 5),
	('x3', NOW() - INTERVAL '7 days', '2030-01-01 00:00:00', 'public', 5),
	('x4', NOW() - INTERVAL '30 days', '2030-01-01 00:00:00', 'public', 5),
	('x5', NOW() - INTERVAL '4 hours', '2030-01-01 00:00:00', 'public', 5),
	('x6', NOW() - INTERVAL '17 hours', '2030-01-01 00:00:00', 'public', 5);
`)
	if err != nil {
		fmt.Printf("failed to seed test data: %s", err.Error())
	}

	var count int
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM pastes").Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	if count != 20 {
		fmt.Printf("Ожидалось 2 записи, получено %d", count)
	}
	fmt.Println("migrations set up ended")
}
