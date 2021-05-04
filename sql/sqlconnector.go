package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb"
)

type DbConnector interface {
	OpenConnection()
	CloseConnection()
	Query(QueryString string) *DataReader
	NonQuery(QueryString string)
	// BulkCopy(src_reader DataReader, String dest_table, ArrayList Mapping)
}

type MssqlConnector struct {
	db         *sql.DB
	ServerName string
	ServerPort string
	Username   string
	Password   string
	Database   string
	Instance   string

	m_connection_str string
}

func NewMssqlConnector(host, port, user, password, dbname, instance string) *MssqlConnector {
	conn := &MssqlConnector{
		ServerName: host,
		ServerPort: port,
		Username:   user,
		Password:   password,
		Database:   dbname,
		Instance:   instance,
	}

	return conn
}

func (conn *MssqlConnector) OpenConnection() error {
	query := url.Values{}
	// query.Add("app name", "MyAppName")
	query.Add("database", conn.Database)
	query.Add("connection+timeout", "60")

	var u *url.URL
	u = &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(conn.Username, conn.Password),
		// Host:   fmt.Sprintf("%s:%s", conn.ServerName, conn.ServerPort),
		// Host:     fmt.Sprintf("%s", conn.ServerName),
		Path:     conn.Instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	if conn.Instance == "" && conn.ServerPort != "" {
		u.Host = fmt.Sprintf("%s:%s", conn.ServerName, conn.ServerPort)
	} else {
		u.Host = fmt.Sprintf("%s", conn.ServerName)
	}

	// Connect to multi instance, need SQL Server Browser service to be running
	fmt.Printf("mssql conn:%s\n", u.String())
	var err error
	conn.db, err = sql.Open("sqlserver", u.String())
	// conn.db, err = sql.Open("sqlserver", "sqlserver://report_user:password@localhost/MSSQL2017?database=local")
	if err != nil {
		return err
	}
	err = checkVersion(conn.db)
	if err != nil {
		return err
	}
	return nil
}

func (conn *MssqlConnector) CloseConnection() {
	conn.db.Close()
}

func (conn *MssqlConnector) Query(QueryString string, args ...interface{}) *DataReader {
	rows, err := conn.db.Query(QueryString, args...)
	if err != nil {
		// log.Fatal(err)
		log.Printf("%s\n", err)
		return nil
	}
	return CreateDataReader(rows)
}

func (conn *MssqlConnector) NonQuery(QueryString string, args ...interface{}) error {
	_, err := conn.db.Query(QueryString, args...)
	if err != nil {
		// log.Fatal(err)
		log.Printf("%s\n", err)
		return err
	}

	/*
		var ctx context.Context
		result, err := conn.db.ExecContext(ctx, QueryString)
		if err != nil {
			log.Fatal(err)
		}
		_, err = result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		} */
	return nil
}

func checkVersion(db *sql.DB) error {
	ctx := context.Background()

	var err error
	err = db.Ping()
	if err != nil {
		log.Printf("Ping database failed: %v\n", err)
		return err
	}
	/*
		err := db.PingContext(ctx)
		if err != nil {
			// log.Fatal("Ping database failed:", err.Error())
			log.Printf("Ping database failed\n")
			return err
		} */

	var version string
	err = db.QueryRowContext(ctx, "SELECT @@version").Scan(&version)
	if err != nil {
		// log.Fatal("Scan failed:", err.Error())
		log.Printf("Scan failed\n")
		return err
	}
	// log.Printf("%s\n", version)
	return nil
}
