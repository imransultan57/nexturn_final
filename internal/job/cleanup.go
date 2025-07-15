package job

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"nexturn_final/internal/logger"
)

func StartCleanupJob(db *pgxpool.Pool, logg *logger.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		<-ticker.C
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := db.Exec(ctx, `DELETE FROM urls WHERE expires_at < NOW()`)
		if err != nil {
			logg.Println("Cleanup job failed:", err)
		} else {
			logg.Println("Expired URLs cleaned up")
		}
		cancel()
	}
}
