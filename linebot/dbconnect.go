package linebot

import (
	"fmt"
	"time"

	"github.com/churikawit/gosrc/sql"
)

var (
	Conn *sql.PostgresConnector
)

func init() {
	// Conn = sql.NewPostgresConnector(server, port, user, password, database)
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func RegisterUser(user_id string, bot_id int, display_name, picture_url, status_message string) {
	Conn.OpenConnection()
	defer Conn.CloseConnection()

	queryString := "select count(*) from register_user where user_id=$1 and bot_id=$2"
	var reader *sql.DataReader = Conn.Query(queryString, user_id, bot_id)
	defer reader.Close()

	reader.Read()
	datacell := reader.GetValue2(0)
	s, ok := datacell.(int32)
	if !ok {
		fmt.Printf("RegisterUser: get count(*) error\n")
	}

	if s <= 0 {
		insert_script := `insert into register_user(user_id, bot_id, display_name, picture_url, status_message, create_time) 
			values($1,$2,$3,$4,$5,$6)`
		Conn.Query(insert_script, user_id, bot_id, display_name, picture_url, status_message, getCurrentTime())
	} else {
		// update
		update_script := `update register_user
			set display_name = $1, picture_url = $2, status_message = $3
			where user_id=$4 and bot_id=$5`
		Conn.Query(update_script, display_name, picture_url, status_message, user_id, bot_id)
	}
}

func GetIntentStage(user_id string, bot_id int) string {
	Conn.OpenConnection()
	defer Conn.CloseConnection()

	queryString := "select intent_stage from register_user where user_id=$1 and bot_id=$2 limit 1"
	var reader *sql.DataReader = Conn.Query(queryString, user_id, bot_id)
	defer reader.Close()

	if reader.Read() {
		datacell := reader.GetValue2(0)
		s, ok := datacell.(string)
		if !ok {
			fmt.Printf("RegisterUser: get count(*) error\n")
			return ""
		}
		return s
	}
	return ""
}

func SetIntentStage(user_id string, bot_id int, intent string) {
	Conn.OpenConnection()
	defer Conn.CloseConnection()

	queryString := "update register_user set intent_stage=$3 where user_id=$1 and bot_id=$2"
	var reader *sql.DataReader = Conn.Query(queryString, string(user_id), bot_id, intent)
	defer reader.Close()
}

func LogEvent(source_type, source_userid, source_groupid, source_roomid string, bot_id int, event_type, event_body string) {
	Conn.OpenConnection()
	defer Conn.CloseConnection()

	queryString := `insert into log_event(source_type, source_userid, source_groupid, source_roomid,
		bot_id, event_type, event_body) values ($1,$2,$3,$4,$5,$6,$7)`
	var reader *sql.DataReader = Conn.Query(queryString, source_type, source_userid, source_groupid,
		source_roomid, bot_id, event_type, event_body)
	defer reader.Close()
}
