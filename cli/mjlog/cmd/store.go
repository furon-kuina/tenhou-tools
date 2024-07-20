package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type YonmaGame struct {
	Id         int       `db:"id"`
	Time       time.Time `db:"time"`
	Url        string    `db:"url"`
	Table      int       `db:"table"`
	Aka        bool      `db:"aka"`
	Nashinashi bool      `db:"nashinashi"`
	Tonpuu     bool      `db:"tonpuu"`
	Player1    string    `db:"player1"`
	Player2    string    `db:"player2"`
	Player3    string    `db:"player3"`
	Player4    string    `db:"player4"`
}

type SanmaGame struct {
	Id         int       `db:"id"`
	Time       time.Time `db:"time"`
	Url        string    `db:"url"`
	Table      int       `db:"table"`
	Aka        bool      `db:"aka"`
	Nashinashi bool      `db:"nashinashi"`
	Tonpuu     bool      `db:"tonpuu"`
	Player1    string    `db:"player1"`
	Player2    string    `db:"player2"`
	Player3    string    `db:"player3"`
}

func AddGame(ctx context.Context, gi *Game, db *sqlx.DB) error {
	var err error
	if len(gi.PlayerID) == 3 {
		sql := `INSERT INTO sanma_game
		(time, url, "table", aka, nashinashi, tonpuu, player1, player2, player3)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
		_, err = db.ExecContext(ctx, sql, gi.Time, gi.LogURL, gi.GameType.Table, gi.GameType.Aka, gi.GameType.NashiNashi, gi.GameType.Tonpuu, gi.PlayerID[0], gi.PlayerID[1], gi.PlayerID[2])
	} else {
		sql := `INSERT INTO yonma_game
			(time, url, "table", aka, nashinashi, tonpuu, player1, player2, player3, player4)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`
		_, err = db.ExecContext(ctx, sql, gi.Time, gi.LogURL, gi.GameType.Table, gi.GameType.Aka, gi.GameType.NashiNashi, gi.GameType.Tonpuu, gi.PlayerID[0], gi.PlayerID[1], gi.PlayerID[2], gi.PlayerID[3])
	}
	return err
}

type Games []*Game

func ListYonmaGame(ctx context.Context, db *sqlx.DB) ([]*YonmaGame, error) {
	sql := `SELECT * FROM yonma_game;`
	games := make([]*YonmaGame, 0)
	if err := db.SelectContext(ctx, &games, sql); err != nil {
		return nil, err
	}
	return games, nil
}

func ListSanmaGame(ctx context.Context, db *sqlx.DB) ([]*SanmaGame, error) {
	sql := `SELECT * FROM sanma_game;`
	games := make([]*SanmaGame, 0)
	if err := db.SelectContext(ctx, &games, sql); err != nil {
		return nil, err
	}
	return games, nil
}

type Config struct {
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort     string `env:"DB_PORT" envDefault:"35432"`
	DBName     string `env:"DB_NAME" envDefault:"tenhou"`
}

func New(ctx context.Context, cfg Config) (*sqlx.DB, func(), error) {
	dst := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("postgres", dst)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	dbx := sqlx.NewDb(db, "postgres")
	return dbx, func() { _ = db.Close() }, nil
}
