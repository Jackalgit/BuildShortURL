package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type DataBase struct {
	conn *sql.DB
}

func NewDataBase(ctx context.Context) DataBase {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}

	query := `CREATE TABLE IF NOT EXISTS storage(correlationId VARCHAR (255), shortURLKey VARCHAR (255), originalURL VARCHAR (255))`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("[Create Table] Не удалось создать таблицу в база данных: %q", err)
	}

	db.ExecContext(ctx, `CREATE UNIQUE INDEX originalURL_idx ON storage (originalURL)`)

	return DataBase{conn: db}
}

func (d DataBase) AddURL(ctx context.Context, shortURLKey string, originalURL []byte) error {

	query := `INSERT INTO storage (shortURLKey, originalURL) VALUES($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	stmt, err := d.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[PrepareContext] %s", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, shortURLKey, originalURL)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("[Insert into DB] Не удалось сделать запись в базу данных: %q", err)

			dupShortURLKey := d.GetShortURLinDB(ctx, originalURL)

			AddURLError := models.NewAddURLError(dupShortURLKey)

			return AddURLError
		}
	}
	return nil

}

func (d DataBase) GetURL(ctx context.Context, shortURLKey string) ([]byte, bool) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	//db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	//if err != nil {
	//	log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	//}
	//defer db.Close()

	row := d.conn.QueryRowContext(
		ctx,
		"SELECT originalURL FROM storage WHERE shortURLKey = $1", fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey),
	)

	var originalURL sql.NullString
	err := row.Scan(&originalURL)
	if err != nil {
		log.Printf("[row Scan] Не удалось преобразовать данные: %q", err)
	}

	if originalURL.Valid {
		log.Printf("Оригинальный УРЛ: %q", originalURL.String)
		return []byte(originalURL.String), true
	}
	return nil, false

}

func (d DataBase) AddBatchURL(ctx context.Context, batchList []models.BatchURL) error {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := d.conn.Begin()
	if err != nil {
		log.Printf("Ошибка начала транзакции: %q", err)
	}

	query := `INSERT INTO storage (correlationId, shortURLKey, originalURL) VALUES($1, $2, $3)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[PrepareContext] %s", err)
	}
	defer stmt.Close()

	for _, v := range batchList {
		_, err = stmt.ExecContext(ctx, v.Correlation, v.ShortURL, v.OriginalURL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				log.Printf("[Insert into DB] Не удалось сделать запись в базу данных: %q", err)

				dupShortURLKey := d.GetShortURLinDB(ctx, []byte(v.OriginalURL))

				AddURLError := models.NewAddURLError(dupShortURLKey)

				return AddURLError
			}
		}

		if err != nil {
			tx.Rollback()
			log.Printf("Ошибка записи в базу: %q", err)
		}
	}
	tx.Commit()
	return nil

}

func (d DataBase) GetShortURLinDB(ctx context.Context, originalURL []byte) string {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	row := d.conn.QueryRowContext(
		ctx,
		"SELECT shortURLKey FROM storage WHERE originalURL = $1", originalURL,
	)

	var shortURLKey sql.NullString
	err := row.Scan(&shortURLKey)
	if err != nil {
		log.Printf("[row Scan] Не удалось преобразовать данные: %q", err)
	}

	if shortURLKey.Valid {
		log.Printf("Оригинальный УРЛ: %q", shortURLKey.String)
		return shortURLKey.String
	}

	return ""

}
