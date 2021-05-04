package sql

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/lib/pq"
)

type PostgresConnector struct {
	db            *sql.DB
	ServerName    string
	ServerPort    string
	Username      string
	Password      string
	Database      string
	ConnectionStr string
}

func NewPostgresConnector(host, port, user, password, dbname string) *PostgresConnector {
	conn := &PostgresConnector{
		ServerName:    host,
		ServerPort:    port,
		Username:      user,
		Password:      password,
		Database:      dbname,
		ConnectionStr: "",
	}
	return conn
}

func NewPostgresConnector2(connStr string) *PostgresConnector {
	conn := &PostgresConnector{
		ServerName:    "",
		ServerPort:    "",
		Username:      "",
		Password:      "",
		Database:      "",
		ConnectionStr: connStr,
	}
	return conn
}

func (conn *PostgresConnector) OpenConnection() error {
	query := url.Values{}
	query.Add("database", conn.Database)
	var u *url.URL
	u = &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(conn.Username, conn.Password),
		Host:     fmt.Sprintf("%s", conn.ServerName),
		RawQuery: query.Encode(),
	}

	if conn.ConnectionStr == "" {
		conn.ConnectionStr = u.String()
	}
	var err error
	conn.db, err = sql.Open("postgres", conn.ConnectionStr)
	if err != nil {
		return err
	}

	// check if connected
	err = conn.db.Ping()
	if err != nil {
		fmt.Printf("postgres: cannot connect\n")
		return err
	}
	fmt.Printf("postgres: connected\n")

	return nil
}

func (conn *PostgresConnector) CloseConnection() {
	conn.db.Close()
}

func (conn *PostgresConnector) Query(QueryString string, args ...interface{}) *DataReader {
	rows, err := conn.db.Query(QueryString, args...)
	if err != nil {
		// log.Fatal(err)
		log.Printf("%s\n", err)
		return nil
	}
	return CreateDataReader(rows)
}

func (conn *PostgresConnector) NonQuery(QueryString string, args ...interface{}) error {
	_, err := conn.db.Query(QueryString, args...)
	if err != nil {
		return err
	}
	return nil
}
