package db

import (
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	"fmt"
	"log"
	"database/sql"
)

var conn *sqlx.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "gera"
	password = "DVA3HFDL=y#tA#:m)WcKuKnU"
	dbname   = "twitchspy"
)

type TwitchGame struct {
	Name        string           `db:"name"`
	Gameid      int              `db:"game_id"`
	Giantbombid sql.NullInt64    `db:"giantbomb_id"`
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

func Connect(dialect string, schema bool) {
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, "disable")
	db, err := sqlx.Connect(dialect, psql)
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

func InsertGame(game *TwitchGame) {
	_, err := conn.NamedExec(
		`
		INSERT INTO games
		(name, game_id, giantbomb_id)
		VALUES
		(:name, :game_id, :giantbomb_id)
		`, game)

	if err != nil {
		panic(err)
	}
}

func execSchema() {
	schema := `CREATE TABLE IF NOT EXISTS public.clients (
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
		OWNER TO gera;`

	_, err := conn.MustExec(schema).RowsAffected()
	if err != nil {
		log.Printf("%s\n", err.Error())
	}
}
