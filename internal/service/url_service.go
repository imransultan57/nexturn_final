package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"nexturn_final/internal/logger"
)

const codeLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLService struct {
	db   *pgxpool.Pool
	logg *logger.Logger
}

func NewURLService(db *pgxpool.Pool, logg *logger.Logger) *URLService {
	return &URLService{db: db, logg: logg}
}

func (s *URLService) Shorten(ctx context.Context, original string, expiresAt time.Time, ip string) (string, error) {
	code := s.generateCode()
	_, err := s.db.Exec(ctx, `
        INSERT INTO urls (original, shortcode, created_at, expires_at, hits)
        VALUES ($1, $2, NOW(), $3, 0)
    `, original, code, expiresAt)
	if err != nil {
		return "", err
	}
	s.logAction(ctx, code, ip, "shorten")
	return code, nil
}

func (s *URLService) Resolve(ctx context.Context, code, ip string) (string, error) {
	var original string
	var expiresAt time.Time
	err := s.db.QueryRow(ctx, `
        SELECT original, expires_at FROM urls WHERE shortcode=$1
    `, code).Scan(&original, &expiresAt)
	if err != nil {
		return "", errors.New("not found")
	}
	if time.Now().After(expiresAt) {
		return "", errors.New("expired")
	}
	s.db.Exec(ctx, `UPDATE urls SET hits = hits + 1 WHERE shortcode=$1`, code)
	s.logAction(ctx, code, ip, "redirect")
	return original, nil
}

func (s *URLService) Stats(ctx context.Context, code string) (int, time.Time, time.Time, error) {
	var hits int
	var createdAt, expiresAt time.Time
	err := s.db.QueryRow(ctx, `
        SELECT hits, created_at, expires_at FROM urls WHERE shortcode=$1
    `, code).Scan(&hits, &createdAt, &expiresAt)
	if err != nil {
		return 0, time.Time{}, time.Time{}, errors.New("not found")
	}
	return hits, createdAt, expiresAt, nil
}

func (s *URLService) generateCode() string {
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *URLService) logAction(ctx context.Context, code, ip, action string) {
	s.db.Exec(ctx, `
        INSERT INTO access_logs (url_id, accessed_at, ip, action)
        SELECT id, NOW(), $1, $2 FROM urls WHERE shortcode=$3
    `, ip, action, code)
}
