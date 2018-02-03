package db

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"log"
	"database/sql"
	"github.com/kelseyhightower/envconfig"
	"TwitchSpy/config"
)

var conn *sqlx.DB

const (
	dbEnvPrefix = "db"
	driverName  = "postgres"
)

type TwitchGame struct {
	Name        string           `db:"name"`
	GameID      int              `db:"game_id"`
	GiantBombID sql.NullInt64    `db:"giantbomb_id"`
	Genres      []sql.NullString `db:"genres"`
	Aliases     []sql.NullString `db:"aliases"`
	Brief       sql.NullString   `db:"brief"`
}

type ClientToken struct {
	AccessToken  string
	RefreshToken string
	Expired      bool
}

type ClientState struct {
	ClientID     string         `db:"client_id"`
	ClientSecret string         `db:"client_secret"`
	AccessToken  sql.NullString `db:"access_token"`
	RefreshToken sql.NullString `db:"refresh_token"`
	Expired      bool           `db:"expired"`
}

func Connect(schema bool) {
	var dbConfig config.DBConfig
	if err := envconfig.Process(dbEnvPrefix, &dbConfig); err != nil {
		panic(err)
	}

	db, err := sqlx.Connect(driverName, dbConfig.ToString())
	if err != nil {
		panic(err)
	}
	conn = db

	if schema {
		execSchema()
	}
}

func Close() {
	conn.Close()
}

func GetClient() *ClientState {
	var meta ClientState
	query := "SELECT client_id, client_secret, access_token, refresh_token, expired FROM clients WHERE rid=$1"

	if err := conn.Get(&meta, query, 1); err != nil {
		panic(err)
	}

	return &meta
}

func UpdateClientToken(token ClientToken) error {
	fields := map[string]interface{}{
		"access":  token.AccessToken,
		"refresh": token.RefreshToken,
		"expired": token.Expired,
	}

	// Careful!!! will update all fields in clients because missing WHERE clause
	_, err := conn.NamedExec(
		`
		UPDATE clients
		SET (access_token, refresh_token, expired) = (:access, :refresh, :expired)
		`, fields)

	return err
}

func GameExists(gameID int) bool {
	var game int
	conn.Get(&game, `SELECT game_id FROM games WHERE game_id = $1`, gameID)
	return game != 0
}

func InsertGame(game *TwitchGame) int64 {
	res, err := conn.NamedExec(
		`
		INSERT INTO games
		(name, game_id, giantbomb_id)
		VALUES
		(:name, :game_id, :giantbomb_id) ON CONFLICT DO NOTHING
		`, game)

	if err != nil {
		panic(err)
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected
}

func execSchema() {
	schema :=
		// CREATE TABLE BEGIN
		`
		CREATE TABLE IF NOT EXISTS public.clients (
		client_id TEXT COLLATE pg_catalog."default" NOT NULL,
		client_secret TEXT COLLATE pg_catalog."default" NOT NULL,
		access_token TEXT COLLATE pg_catalog."default",
		refresh_token TEXT COLLATE pg_catalog."default",
		expired BOOLEAN NOT NULL DEFAULT FALSE,
		rid INTEGER NOT NULL DEFAULT nextval('client_rid_seq'::REGCLASS),
		CONSTRAINT client_pkey PRIMARY KEY (rid)
		)
		WITH (
			OIDS = FALSE
		)
		TABLESPACE pg_default;

		ALTER TABLE public.clients
			OWNER TO gera;
		`
		// CREATE TABLE END

	_, err := conn.MustExec(schema).RowsAffected()
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
}
