package database

import (
	"context"
	"database/sql"
	"github.com/Jackalgit/BuildShortURL/cmd/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type DataBase struct {
	//Connect *sql.DB
}

func NewDataBase(ctx context.Context) DataBase {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS storage_url(id int primary key, shortURLKey text, originalURL text)")
	if err != nil {
		log.Printf("[Create DB] Не удалось создать таблицу в база данных: %q", err)
	}
	log.Print("Создана таблицы для хранения УРЛ")
	//rows, err := res.RowsAffected()
	//if err != nil {
	//	log.Printf("Ошибка в получении количества строк: %q", err)
	//}
	//last, err := res.LastInsertId()
	//if err != nil {
	//	log.Printf("Ошибка в получении LastInsertId: %q", err)
	//}
	//
	//log.Printf("Количество строк: %d", rows)
	//log.Printf("LastInsertId: %d", last)

	return DataBase{
		//Connect: db,
	}
}

func (d DataBase) AddURL(ctx context.Context, shortURLKey string, originalURL []byte) {

	log.Print("Вызван метод добавления урл")

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", config.Config.DatabaseDSN)
	if err != nil {
		log.Printf("[Open DB] Не удалось установить соединение с базой данных: %q", err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctx,
		"INSERT INTO storage_url(shortURLKey, originalURL)"+" VALUES(?,?)", shortURLKey, originalURL)
	if err != nil {
		log.Printf("[Insert into DB] Не удалось сделать запись в базу данных: %q", err)
	}
	//rows, err := res.RowsAffected()
	//if err != nil {
	//	log.Printf("Ошибка в получении количества строк: %q", err)
	//}
	//last, err := res.LastInsertId()
	//if err != nil {
	//	log.Printf("Ошибка в получении LastInsertId: %q", err)
	//}
	//
	//log.Printf("Количество строк: %d", rows)
	//log.Printf("LastInsertId: %d", last)

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
		"SELECT originalURL FROM storage_url WHERE shortURLKey = ?", shortURLKey,
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
