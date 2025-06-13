package ticketbot

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func initDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", "bot.db")
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	db.MustExec(boardSettingsSchema)
	db.MustExec(usersTableSchema)

	return db, nil
}
