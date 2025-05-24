package data

import (
	"database/sql"
	"fmt"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error

	username := "root"
	password := "reakgo_user"
	host := "127.0.0.1"
	port := "3306"
	database := "magic_login_db"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username, password, host, port, database)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to MySQL:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Failed to ping MySQL:", err)
	}

	fmt.Println("✅ Connected to MySQL database!")
}
