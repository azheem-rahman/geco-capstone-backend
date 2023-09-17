package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type user struct {
	Account_id   int    `json:"account_id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Account_Type string `json:"account_type"`
}

// type usersDetails struct {
// 	Detail_id  int    `json:"detail_id"`
// 	Account_id int    `json:"account_id"`
// 	First_name string `json:"first_name"`
// 	Last_name  string `json:"last_name"`
// }

func setupDBConnection() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "Password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "capstonedb",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("DB Connected!")
}

func main() {
	setupDBConnection()

	router := gin.Default()

	router.GET("/accounts", getAccounts)
	router.POST("/new-account", postAccount)

	router.Run("localhost:8080")
}

func getAccounts(c *gin.Context) {
	accounts, err := getAccountsFromDB()
	if err != nil {
		fmt.Println(err.Error())
	}

	c.IndentedJSON(http.StatusOK, accounts)
}

func getAccountsFromDB() ([]user, error) {
	var accounts []user

	rows, err := db.Query("SELECT * FROM accounts")
	if err != nil {
		return nil, fmt.Errorf("getAccountsFromDB %v", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var account user
		if err := rows.Scan(&account.Account_id, &account.Email, &account.Password, &account.Account_Type); err != nil {
			return nil, fmt.Errorf("getAccountsFromDB %v", err.Error())
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAccountsFromDB %v", err.Error())
	}

	return accounts, nil
}

func postAccount(c *gin.Context) {
	var newUser user

	c.BindJSON(&newUser)

	result, err := postAccountToDB(newUser)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.IndentedJSON(http.StatusOK, result)
}

func postAccountToDB(newUser user) (int64, error) {
	// check if email username is already taken -- account already exists
	var emailFound bool

	if err := db.QueryRow("SELECT * FROM accounts WHERE email=?", newUser.Email).Scan(&emailFound); err != nil {
		// email username already taken since found in database
		if err != sql.ErrNoRows {
			fmt.Println("Email taken!")
		}
	}

	result, err := db.Exec("INSERT INTO accounts (email, password, account_type) VALUES (?, ?, ?)", newUser.Email, newUser.Password, newUser.Account_Type)
	if err != nil {
		fmt.Println(err.Error())
	}

	newUserID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
	}

	return newUserID, nil

}
