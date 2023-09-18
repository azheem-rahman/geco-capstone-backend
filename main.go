package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type user struct {
	Account_id   int    `json:"account_id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Account_Type string `json:"account_type"`
}

type emailPassword struct {
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

	// Routes related to user account login and creation
	router.POST("/login", login)
	router.GET("/accounts", getAccounts)
	router.POST("/new-account", postAccount)
	// router.PATCH("/update-account-password", updateAccountPassword)
	// router.DELETE("/delete-account", deleteAccount) only Admin
	// router.PATCH("/update-account-details", updateAccountDetails)
	// router.GET("/accounts-details", getAccountsDetails)
	// router.POST("/new-account-details", postAccountDetails)

	// Routes related to orders
	// router.GET("/orders", getOrders)
	// router.POST("/new-order", postOrder)

	router.Run("localhost:8080")
}

func login(c *gin.Context) {
	// Get email and password from request body
	var reqBody emailPassword
	var accountFoundInDB user
	// var emailFound bool
	// var passwordMatched bool

	// Returns Error HTTP Bad Request 400 if unable to read from request body
	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to read request body"})
		return
	}

	// Look up email of login account in database
	if err := db.QueryRow("SELECT * FROM accounts WHERE email=?", reqBody.Email).Scan(&accountFoundInDB.Account_id, &accountFoundInDB.Email, &accountFoundInDB.Password, &accountFoundInDB.Account_Type); err != nil {
		// if response returns a no row means email does not exist in database => account does not exist
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid email or password"})
			return
		}
	}

	// Compare password of login account with hashed password of account in database
	err := bcrypt.CompareHashAndPassword([]byte(accountFoundInDB.Password), []byte(reqBody.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": accountFoundInDB.Email,
		// expiration of token will be 1 day (24hours)
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get complete encoded JWT token as a string using secret key stored in .env file
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to create token"})
		return
	}

	// Set JWT token as cookie; cookie expires in 1 day (3600*24seconds)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorisation", tokenString, 3600*24, "", "", false, true)

	// Return HTTP OK 200
	c.JSON(http.StatusOK, gin.H{})
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
	// Get email and password for new account from request body
	var reqBody emailPassword
	var emailFound bool

	// Returns Error HTTP Bad Request 400 if unable to read from request body
	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to read request body"})
		return
	}

	// Check if email already taken
	if err := db.QueryRow("SELECT * FROM accounts WHERE email=?", reqBody.Email).Scan(&emailFound); err != nil {
		// if response returns a row means email already exists in database
		if err != sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Email taken"})
			return
		}
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to hash password"})
		return
	}

	// Create the user
	var newAccount emailPassword
	newAccount.Email = reqBody.Email
	newAccount.Password = string(hash)
	newAccount.Account_Type = reqBody.Account_Type

	_, err = db.Exec("INSERT INTO accounts (email, password, account_type) VALUES (?, ?, ?)", newAccount.Email, newAccount.Password, newAccount.Account_Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to create user"})

		return
	}

	// Respond
	c.IndentedJSON(http.StatusOK, gin.H{"Success": "Successfully created user"})
}
