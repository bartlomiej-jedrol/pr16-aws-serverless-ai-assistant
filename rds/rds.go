package rds

import (
	"context"
	"fmt"
	"time"

	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertNote(ctx context.Context, pool *pgxpool.Pool, title, content string) error {
	function := "InsertNote"

	query := `
    INSERT INTO notes (title, content)
    VALUES ($1, $2)
    RETURNING id, created_at`

	var id int
	var createdAt time.Time
	err := pool.QueryRow(ctx, query, title, content).Scan(&id, &createdAt)
	if err != nil {
		iLog.Error("failed to execute query row", query, err, configuration.ServiceName, function)
		return fmt.Errorf("failed to execute query row: %v", err)
	}

	fmt.Printf("Inserted note with id: %d, created at: %v", id, createdAt)
	return nil
}
