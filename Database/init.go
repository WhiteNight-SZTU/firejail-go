package database

import (
    "database/sql"
    "fmt"
    "github.com/go-sql-driver/mysql"
)

var db *sql.DB
type databaseConfigure struct {
	username string `json:"username"`
	password string `json:"password"`
	dbname string `json:"dbname"`
}

func InitDB(Configure string) error {
    var err error
    db, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        return fmt.Errorf("打开数据库失败: %v", err)
    }

    if err = db.Ping(); err != nil {
        return fmt.Errorf("连接数据库失败: %v", err)
    }

    return nil
}

func readConfigure() databaseConfigure{
	jsonFile, err := os.Open("database.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configure databaseConfigure
	if err := json.Unmarshal(byteValue, &configure); err != nil {
		fmt.Println(err)
	}
	return configure
}