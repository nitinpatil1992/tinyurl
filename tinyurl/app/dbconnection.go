package app

import (
	"database/sql"
	"time"
	"tinyurl/app/model"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func CreateDBCon() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/tinyurl?parseTime=true")
	if err != nil {
		log.Warn(err.Error())
	} else {
		log.Info("db is connected")
	}
	//defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		log.Warn("MySQL db is not connected")
		log.Warn(err.Error())
	}
	return db
}

func InsertTinyURLs(shortString, longURL string, currentTime time.Time) error {
	const insertStatement = `
	INSERT INTO tiny_urls (
    	short_url, long_url, created_at
  	) VALUES (?, ?, ?)`

	stmtIns, err := Conn.Prepare(insertStatement)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(shortString, longURL, currentTime)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func SelectTinyURL(shortURL string) (model.TinyUrl, error) {
	var result model.TinyUrl

	stmtSelect, err := Conn.Prepare("SELECT * FROM tiny_urls WHERE short_url = ?")
	if err != nil {
		log.Warn(err.Error())
		return result, err
	}
	defer stmtSelect.Close()

	err = stmtSelect.QueryRow(shortURL).Scan(&result.ShortURL, &result.LongURL, &result.CreatedAt)

	//err = stmtSelect.QueryRow(shortURL).Scan(&result.ShortURL, &result.LongURL, &result.CreatedAt)
	if err != nil {
		log.Warn(err.Error())
		return result, err
	}
	log.Info("select tiny url :", result)
	return result, err
}

func InsertTinyURLVisits(shortURL string, currentTime time.Time) error {
	const insertStatement = `
	INSERT INTO tiny_url_visits (
	  short_url, visited_at
	) VALUES (?, ?)`

	stmtIns, err := Conn.Prepare(insertStatement)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(shortURL, time.Now())
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func CountTinyURLVisited(shortURL string) (int, error) {

	stmtSelect, err := Conn.Prepare("SELECT count(short_url) as count FROM tiny_url_visits WHERE short_url = ?")
	if err != nil {
		log.Warn(err.Error())
		return 0, err
	}
	defer stmtSelect.Close()

	var count int
	err = stmtSelect.QueryRow(shortURL).Scan(&count)

	if err != nil {
		log.Warn(err.Error())
		return 0, err
	}
	return count, err
}
