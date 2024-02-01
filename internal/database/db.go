package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/google/uuid"
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

	query := `CREATE TABLE IF NOT EXISTS storage(
    correlationId VARCHAR (255),
    userID VARCHAR (255),
    shortURLKey VARCHAR (255),
    originalURL VARCHAR (255)
    )`

	_, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("[Create Table] Не удалось создать таблицу в база данных: %q", err)
	}

	db.ExecContext(ctx, `CREATE UNIQUE INDEX originalURL_idx ON storage (originalURL)`)

	return DataBase{conn: db}
}

func (d DataBase) AddURL(ctx context.Context, userID uuid.UUID, shortURLKey string, originalURL []byte) error {

	query := `INSERT INTO storage (userID, shortURLKey, originalURL) VALUES($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	stmt, err := d.conn.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[PrepareContext] %s", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, shortURLKey, originalURL)
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

func (d DataBase) GetURL(ctx context.Context, userID uuid.UUID, shortURLKey string) ([]byte, bool) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	row := d.conn.QueryRowContext(
		ctx,
		"SELECT originalURL FROM storage WHERE userID = $1 AND shortURLKey = $2",
		userID,
		fmt.Sprint(config.Config.BaseAddress, "/", shortURLKey),
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

func (d DataBase) AddBatchURL(ctx context.Context, userID uuid.UUID, batchList []models.BatchURL) error {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := d.conn.Begin()
	if err != nil {
		log.Printf("Ошибка начала транзакции: %q", err)
	}

	query := `INSERT INTO storage (correlationId, userID, shortURLKey, originalURL) VALUES($1, $2, $3, $4)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[PrepareContext] %s", err)
	}
	defer stmt.Close()

	for _, v := range batchList {
		_, err = stmt.ExecContext(ctx, v.Correlation, userID, v.ShortURL, v.OriginalURL)
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

func (d DataBase) UserURLList(ctx context.Context, userID uuid.UUID) ([]models.ResponseUserURL, bool) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	rows, err := d.conn.QueryContext(
		ctx,
		"SELECT shortURLKey, originalURL FROM storage WHERE userID = $1",
		userID,
	)
	if err != nil {
		log.Printf("[QueryContext] Не удалось получить данные по userId: %q", err)
	}
	defer rows.Close()

	var userURLList []models.ResponseUserURL
	var userURL models.ResponseUserURL

	for rows.Next() {

		err = rows.Scan(&userURL.ShortURL, &userURL.OriginalURL)
		if err != nil {
			log.Printf("[rows Scan] Не удалось преобразовать данные: %q", err)
		}

		userURLList = append(userURLList, userURL)
	}
	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		log.Printf("[rows Err]: %q", err)
	}

	// если список пустой, то в базе нет записей
	if len(userURLList) == 0 {
		return userURLList, false
	}

	return userURLList, true

}
