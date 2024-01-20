package database

import (
	"context"
	"database/sql"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type DataBase struct{}

func NewDataBase(ctx context.Context) DataBase {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS storage(correlationId VARCHAR (255), shortURLKey VARCHAR (255), originalURL VARCHAR (255))`)
	if err != nil {
		log.Printf("[Create Table] Не удалось создать таблицу в база данных: %q", err)
	}

	log.Print("Создана таблица для хранения УРЛ")

	return DataBase{}
}

func (d DataBase) AddURL(ctx context.Context, shortURLKey string, originalURL []byte) {

	log.Print("Вызван метод добавления урл")

	query := `INSERT INTO storage (shortURLKey, originalURL) VALUES($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}
	defer db.Close()

	//_, err = db.ExecContext(ctx, query, shortURLKey, originalURL)
	//if err != nil {
	//	log.Printf("[Insert into DB] Не удалось сделать запись в базу данных: %q", err)
	//}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[PrepareContext] %s", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, shortURLKey, originalURL)
	if err != nil {
		log.Printf("[Insert into DB] Не удалось сделать запись в базу данных: %q", err)
	}

}

func (d DataBase) GetURL(ctx context.Context, shortURLKey string) ([]byte, bool) {

	log.Print("Вызван метод получения урл")
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}
	defer db.Close()

	row := db.QueryRowContext(
		ctx,
		"SELECT originalURL FROM storage WHERE shortURLKey = $1", shortURLKey,
	)

	var originalURL sql.NullString
	err = row.Scan(&originalURL)
	if err != nil {
		log.Printf("[row Scan] Не удалось преобразовать данные: %q", err)
	}

	if originalURL.Valid {
		log.Printf("Оригинальный УРЛ: %q", originalURL.String)
		return []byte(originalURL.String), true
	}
	return nil, false

}

func (d DataBase) AddBatchURL(ctx context.Context, batchList *models.BatchList) {
	log.Print("Вызван метод AddBatchURL")

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}
	defer db.Close()

	// начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Ошибка начала транзакции: %q", err)
	}

	query := `INSERT INTO storage (correlationId, shortURLKey, originalURL) VALUES($1, $2, $3)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("[PrepareContext] %s", err)
	}
	defer stmt.Close()

	for _, v := range batchList.List {
		//все изменения записываются в транзакцию
		_, err = stmt.ExecContext(ctx, v.Correlation, v.ShortURL, v.OriginalURL)
		if err != nil {
			log.Printf("[Insert into DB] Не удалось сделать запись в базу данных: %q", err)
		}

		if err != nil {
			// если ошибка, то откатываем изменения
			tx.Rollback()
			log.Printf("Ошибка вставки в базу: %q", err)
			return
		}
	}
	// завершаем транзакцию
	tx.Commit()

}
