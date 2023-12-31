package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var salesChannelDB *sql.DB

//=============================//
// Struct for Sales Channel DB //

type ordersFromSales struct {
	DueDate             string `json:"due_date"`
	Completed           int    `json:"completed"`
	OrderLength         int    `json:"order_length"`
	OrderWidth          int    `json:"order_width"`
	OrderHeight         int    `json:"order_height"`
	OrderWeight         int    `json:"order_weight"`
	ConsigneeName       string `json:"consignee_name"`
	ConsigneeNumber     string `json:"consignee_number"`
	ConsigneeCountry    string `json:"consignee_country"`
	ConsigneeAddress    string `json:"consignee_address"`
	ConsigneePostal     string `json:"consignee_postal"`
	ConsigneeState      string `json:"consignee_state"`
	ConsigneeCity       string `json:"consignee_city"`
	ConsigneeProvince   string `json:"consignee_province"`
	ConsigneeEmail      string `json:"consignee_email"`
	PickupContactName   string `json:"pickup_contact_name"`
	PickupContactNumber string `json:"pickup_contact_number"`
	PickupCountry       string `json:"pickup_country"`
	PickupAddress       string `json:"pickup_address"`
	PickupPostal        string `json:"pickup_postal"`
	PickupState         string `json:"pickup_state"`
	PickupCity          string `json:"pickup_city"`
	PickupProvince      string `json:"pickup_province"`
}

// end Struct for Sales Channel DB //
//=================================//

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

type userDetails struct {
	Detail_id  int    `json:"detail_id"`
	Account_id int    `json:"account_id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type userFirstLastNameEmail struct {
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type email struct {
	Email string `json:"email"`
}

type orderIdAccountId struct {
	OrderId   int `json:"order_id"`
	AccountId int `json:"account_id"`
}

type orderIdCompleted struct {
	OrderId   int `json:"order_id"`
	Completed int `json:"completed"`
}

type order struct {
	OrderId   int           `json:"order_id"`
	AccountId sql.NullInt64 `json:"account_id,omitempty"`
	// AccountId           int    `json:"account_id"`
	OrderLength         int    `json:"order_length"`
	OrderWidth          int    `json:"order_width"`
	OrderHeight         int    `json:"order_height"`
	OrderWeight         int    `json:"order_weight"`
	ConsigneeName       string `json:"consignee_name"`
	ConsigneeNumber     string `json:"consignee_number"`
	ConsigneeCountry    string `json:"consignee_country"`
	ConsigneeAddress    string `json:"consignee_address"`
	ConsigneePostal     string `json:"consignee_postal"`
	ConsigneeState      string `json:"consignee_state"`
	ConsigneeCity       string `json:"consignee_city"`
	ConsigneeProvince   string `json:"consignee_province"`
	ConsigneeEmail      string `json:"consignee_email"`
	PickupContactName   string `json:"pickup_contact_name"`
	PickupContactNumber string `json:"pickup_contact_number"`
	PickupCountry       string `json:"pickup_country"`
	PickupAddress       string `json:"pickup_address"`
	PickupPostal        string `json:"pickup_postal"`
	PickupState         string `json:"pickup_state"`
	PickupCity          string `json:"pickup_city"`
	PickupProvince      string `json:"pickup_province"`
	DueDate             string `json:"due_date"`
	Completed           int    `json:"completed"`
}

type orderWithoutId struct {
	AccountId           int    `json:"account_id"`
	OrderLength         int    `json:"order_length"`
	OrderWidth          int    `json:"order_width"`
	OrderHeight         int    `json:"order_height"`
	OrderWeight         int    `json:"order_weight"`
	ConsigneeName       string `json:"consignee_name"`
	ConsigneeNumber     string `json:"consignee_number"`
	ConsigneeCountry    string `json:"consignee_country"`
	ConsigneeAddress    string `json:"consignee_address"`
	ConsigneePostal     string `json:"consignee_postal"`
	ConsigneeState      string `json:"consignee_state"`
	ConsigneeCity       string `json:"consignee_city"`
	ConsigneeProvince   string `json:"consignee_province"`
	ConsigneeEmail      string `json:"consignee_email"`
	PickupContactName   string `json:"pickup_contact_name"`
	PickupContactNumber string `json:"pickup_contact_number"`
	PickupCountry       string `json:"pickup_country"`
	PickupAddress       string `json:"pickup_address"`
	PickupPostal        string `json:"pickup_postal"`
	PickupState         string `json:"pickup_state"`
	PickupCity          string `json:"pickup_city"`
	PickupProvince      string `json:"pickup_province"`
	DueDate             string `json:"due_date"`
	Completed           int    `json:"completed"`
}

type newOrderFromFrontend struct {
	AccountId           int    `json:"account_id"`
	OrderLength         int    `json:"order_length"`
	OrderWidth          int    `json:"order_width"`
	OrderHeight         int    `json:"order_height"`
	OrderWeight         int    `json:"order_weight"`
	ConsigneeName       string `json:"consignee_name"`
	ConsigneeNumber     string `json:"consignee_number"`
	ConsigneeCountry    string `json:"consignee_country"`
	ConsigneeAddress    string `json:"consignee_address"`
	ConsigneePostal     string `json:"consignee_postal"`
	ConsigneeState      string `json:"consignee_state"`
	ConsigneeCity       string `json:"consignee_city"`
	ConsigneeProvince   string `json:"consignee_province"`
	ConsigneeEmail      string `json:"consignee_email"`
	PickupContactName   string `json:"pickup_contact_name"`
	PickupContactNumber string `json:"pickup_contact_number"`
	PickupCountry       string `json:"pickup_country"`
	PickupAddress       string `json:"pickup_address"`
	PickupPostal        string `json:"pickup_postal"`
	PickupState         string `json:"pickup_state"`
	PickupCity          string `json:"pickup_city"`
	PickupProvince      string `json:"pickup_province"`
	DueDate             string `json:"due_date"`
	Completed           int    `json:"completed"`
}

func setupSalesChannelDBConnection() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "Password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "capstonesaleschanneldb",
	}

	var err error
	salesChannelDB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := salesChannelDB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to Sales Channel DB!")
}

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
	setupSalesChannelDBConnection()
	setupDBConnection()

	router := gin.Default()

	// To enable CORS Support
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Route to pull orders from sales channel DB
	router.GET("/pull-orders-from-sales-channel", getOrdersFromSales)

	// Routes related to user account login and creation
	router.POST("/login", login)
	router.GET("/accounts", getAccounts)
	router.POST("/new-account", postAccount)
	router.GET("/is-logged-in", auth, accountIsLoggedIn)
	// router.PATCH("/update-account-password", updateAccountPassword)
	// router.DELETE("/delete-account", deleteAccount) only Admin
	// router.PATCH("/update-account-details", updateAccountDetails)
	router.POST("/new-account-details", postAccountDetails)
	router.GET("/account-details", getAccountDetails)

	// Routes related to orders
	router.GET("/orders", getOrders)
	router.POST("/new-order", postOrder)
	router.PATCH("/update-order-status", updateOrderStatus)
	router.PATCH("/assign-order", assignOrder)

	router.Run("localhost:8080")
}

func getOrdersFromSales(c *gin.Context) {
	var orders []ordersFromSales

	// Get rows of orders with all details from capstone
	rows, err := salesChannelDB.Query("SELECT due_date, completed, order_length,order_width,order_height,order_weight,consignee_name,consignee_number,consignee_country,consignee_address,consignee_postal,consignee_state,consignee_city,consignee_province,consignee_email,pickup_contact_name,pickup_contact_number,pickup_country,pickup_address,pickup_postal,pickup_state,pickup_city,pickup_province FROM orders JOIN order_details ON order_details.order_id = orders.order_id JOIN consignee_details ON consignee_details.order_id = orders.order_id JOIN pickup_details ON pickup_details.order_id = orders.order_id")
	// if err from getting rows of orders from DB, return HTTP Bad Request 400
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to retrieve orders from DB"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var currentOrder ordersFromSales
		// scan each row of order and save to currentOrder
		if err := rows.Scan(
			&currentOrder.DueDate,
			&currentOrder.Completed,
			&currentOrder.OrderLength,
			&currentOrder.OrderWidth,
			&currentOrder.OrderHeight,
			&currentOrder.OrderWeight,
			&currentOrder.ConsigneeName,
			&currentOrder.ConsigneeNumber,
			&currentOrder.ConsigneeCountry,
			&currentOrder.ConsigneeAddress,
			&currentOrder.ConsigneePostal,
			&currentOrder.ConsigneeState,
			&currentOrder.ConsigneeCity,
			&currentOrder.ConsigneeProvince,
			&currentOrder.ConsigneeEmail,
			&currentOrder.PickupContactName,
			&currentOrder.PickupContactNumber,
			&currentOrder.PickupCountry,
			&currentOrder.PickupAddress,
			&currentOrder.PickupPostal,
			&currentOrder.PickupState,
			&currentOrder.PickupCity,
			&currentOrder.PickupProvince); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to save orders from DB"})
			return
		}
		// add currentOrder to orders slice
		orders = append(orders, currentOrder)
	}

	// INSERT each order from orders slice into capstonedb (client's db)
	for _, value := range orders {
		_, err := db.Exec("INSERT INTO orders (order_length, order_width, order_height, order_weight, consignee_name, consignee_number, consignee_country, consignee_address, consignee_postal, consignee_state, consignee_city, consignee_province, consignee_email, pickup_contact_name, pickup_contact_number, pickup_country, pickup_address, pickup_postal, pickup_state, pickup_city, pickup_province, due_date, completed) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", value.OrderLength, value.OrderWidth, value.OrderHeight, value.OrderWeight, value.ConsigneeName, value.ConsigneeNumber, value.ConsigneeCountry, value.ConsigneeAddress, value.ConsigneePostal, value.ConsigneeState, value.ConsigneeCity, value.ConsigneeProvince, value.ConsigneeEmail, value.PickupContactName, value.PickupContactNumber, value.PickupCountry, value.PickupAddress, value.PickupPostal, value.PickupState, value.PickupCity, value.PickupProvince, value.DueDate, value.Completed)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to create order in DB"})
			return
		}
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "successfully pulled orders from sales channel", "orders": orders})
}

func login(c *gin.Context) {
	// Get email and password from request body
	var reqBody emailPassword
	var accountFoundInDB user
	var accountDetailsFoundInDB userDetails

	// Returns Error HTTP Bad Request 400 if unable to read from request body
	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to read request body"})
		return
	}

	// Look up email of login account in accounts database and retrieve account id, email, password, account_type
	rows, err := db.Query("SELECT * FROM accounts WHERE email=?", reqBody.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to retrieve login credentials in database"})
		return
	}
	if rows.Next() {
		// Account found with email provided
		rows.Scan(&accountFoundInDB.Account_id, &accountFoundInDB.Email, &accountFoundInDB.Password, &accountFoundInDB.Account_Type)
	} else {
		// No Account found with email provided
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid email or password"})
		return
	}

	// Look up email of login account in accounts_details database and retrieve first_name, last_name
	rows, err = db.Query("SELECT * FROM accounts_details WHERE account_id=?", accountFoundInDB.Account_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to retrieve account details in database"})
		return
	}
	if rows.Next() {
		// Account details found with email provided
		rows.Scan(&accountDetailsFoundInDB.Detail_id, &accountDetailsFoundInDB.Account_id, &accountDetailsFoundInDB.First_name, &accountDetailsFoundInDB.Last_name)
	} else {
		// No Account details found with email provided
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "No account details found"})
		return
	}

	// Compare password of login account with hashed password of account in database
	err = bcrypt.CompareHashAndPassword([]byte(accountFoundInDB.Password), []byte(reqBody.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid email or password"})
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
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to create token"})
		return
	}

	// Set JWT token as cookie; cookie expires in 1 day (3600*24seconds)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorisation", tokenString, 3600*24, "", "", false, true)

	// Return HTTP OK 200
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Login successful", "firstName": accountDetailsFoundInDB.First_name, "lastName": accountDetailsFoundInDB.Last_name, "accountType": accountFoundInDB.Account_Type, "accessToken": tokenString})
}

func accountIsLoggedIn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "User account already logged in",
	})
}

func auth(c *gin.Context) {
	// Get cookie from request body
	tokenString, err := c.Cookie("Authorisation")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode and validate cookie
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Cookie successfully validated, user account has access

		// Check cookie expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			// Abort if cookie expired (current time greater than cookie expiration)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find user account with token sub
		var foundAccount user
		if err := db.QueryRow("SELECT * FROM accounts WHERE email=?", claims["sub"]).Scan(&foundAccount.Account_id, &foundAccount.Email, &foundAccount.Password, &foundAccount.Account_Type); err != nil {
			// account email not found
			if err == sql.ErrNoRows {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		}

		// Attach user account to request
		c.Set("user", foundAccount)

		// Continue
		c.Next()
	} else {
		// Cookie not successfully validated => abort
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func getAccounts(c *gin.Context) {
	var accounts []user

	// Get rows of accounts from DB
	rows, err := db.Query("SELECT * FROM accounts")
	// if err from getting rows of accounts from DB, return HTTP Bad Request 400
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to retrieve accounts from DB"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var account user
		// scan each row of accounts and save to account
		if err := rows.Scan(&account.Account_id, &account.Email, &account.Password, &account.Account_Type); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to save accounts from DB"})
			return
		}
		// add account to accounts slice
		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "successfully retrieved profiles from DB", "profiles": accounts})
}

func postAccount(c *gin.Context) {
	// Get email and password for new account from request body
	var reqBody emailPassword
	// var emailFound bool

	// Returns Error HTTP Bad Request 400 if unable to read from request body
	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to read request body"})
		return
	}

	// Check if email already taken
	rows, err := db.Query("SELECT * FROM accounts WHERE email=?", reqBody.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to check if email taken in database"})
		return
	}

	if rows.Next() {
		// Email already taken in database
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Email taken"})
		return
	} else {
		// Email available to use for new account

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to hash password"})
			return
		}

		// Create the user
		var newAccount emailPassword
		newAccount.Email = reqBody.Email
		newAccount.Password = string(hash)
		newAccount.Account_Type = reqBody.Account_Type

		_, err = db.Exec("INSERT INTO accounts (email, password, account_type) VALUES (?, ?, ?)", newAccount.Email, newAccount.Password, newAccount.Account_Type)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to create user"})
			return
		}

		// Respond
		c.IndentedJSON(http.StatusOK, gin.H{"status": http.StatusOK, "newAccountEmailCreated": newAccount.Email})
	}
}

func postAccountDetails(c *gin.Context) {
	// Get account details for new account from request body
	var reqBody userFirstLastNameEmail
	var accountID int

	// Return Error HTTP Bad Request 400 if unable to read from request body
	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to read request body"})
	}

	// Find and Save Account ID from database using Account Email provided in request body
	rows, err := db.Query("SELECT account_id FROM accounts WHERE email=?", reqBody.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to check if account exists in database"})
		return
	}

	if rows.Next() {
		// AccountID found
		rows.Scan(&accountID)

		// Create new account detail for new account
		var newAccountDetail userDetails
		newAccountDetail.Account_id = accountID
		newAccountDetail.First_name = reqBody.First_name
		newAccountDetail.Last_name = reqBody.Last_name

		// Save new account detail to DB
		_, err = db.Exec("INSERT INTO accounts_details (account_id, first_name, last_name) VALUES (?, ?, ?)", newAccountDetail.Account_id, newAccountDetail.First_name, newAccountDetail.Last_name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to create account details"})
			return
		}

		// Respond
		c.IndentedJSON(http.StatusOK, gin.H{"status": http.StatusOK, "newAccountDetailCreated": reqBody.Email})
	} else {
		// No AccountID found
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "No Account Found"})
		return
	}
}

func getAccountDetails(c *gin.Context) {
	var reqBody email
	var accountDetails userDetails
	var accountID int

	// Return Error HTTP Bad Request 400 if unable to read from request body
	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to read request body"})
	}

	// Find and Save Account ID from database using Account Email provided in request body
	if err := db.QueryRow("SELECT account_id FROM accounts WHERE email=?", reqBody.Email).Scan(&accountID); err != nil {
		// if response returns no rows means no account found in database
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No Account Found"})
			return
		}
	}

	// Find and Save account details from DB
	if err := db.QueryRow("SELECT * FROM accounts_details WHERE account_id=?", accountID).Scan(&accountDetails.Detail_id, &accountDetails.Account_id, &accountDetails.First_name, &accountDetails.Last_name); err != nil {
		// if response returns no rows means no account details found in database
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No Account Details Found"})
			return
		}
	}

	// Respond
	c.IndentedJSON(http.StatusOK, accountDetails)
}

func postOrder(c *gin.Context) {
	// Get order details from request body
	var reqBody newOrderFromFrontend

	// Returns Error HTTP Bad Request 400 if unable to read from request body
	if c.ShouldBindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to read request body", "reqBody": reqBody})
		return
	}

	// Create the order in database
	var newOrder orderWithoutId
	newOrder.AccountId = reqBody.AccountId
	newOrder.OrderLength = reqBody.OrderLength
	newOrder.OrderWidth = reqBody.OrderWidth
	newOrder.OrderHeight = reqBody.OrderHeight
	newOrder.OrderWeight = reqBody.OrderWeight
	newOrder.ConsigneeName = reqBody.ConsigneeName
	newOrder.ConsigneeNumber = reqBody.ConsigneeNumber
	newOrder.ConsigneeCountry = reqBody.ConsigneeCountry
	newOrder.ConsigneeAddress = reqBody.ConsigneeAddress
	newOrder.ConsigneePostal = reqBody.ConsigneePostal
	newOrder.ConsigneeState = reqBody.ConsigneeState
	newOrder.ConsigneeCity = reqBody.ConsigneeCity
	newOrder.ConsigneeProvince = reqBody.ConsigneeProvince
	newOrder.ConsigneeEmail = reqBody.ConsigneeEmail
	newOrder.PickupContactName = reqBody.PickupContactName
	newOrder.PickupContactNumber = reqBody.PickupContactNumber
	newOrder.PickupCountry = reqBody.PickupCountry
	newOrder.PickupAddress = reqBody.PickupAddress
	newOrder.PickupPostal = reqBody.PickupPostal
	newOrder.PickupState = reqBody.PickupState
	newOrder.PickupCity = reqBody.PickupCity
	newOrder.PickupProvince = reqBody.PickupProvince
	newOrder.DueDate = reqBody.DueDate
	newOrder.Completed = reqBody.Completed

	_, err := db.Exec("INSERT INTO orders (account_id, order_length, order_width, order_height, order_weight, consignee_name, consignee_number, consignee_country, consignee_address, consignee_postal, consignee_state, consignee_city, consignee_province, consignee_email, pickup_contact_name, pickup_contact_number, pickup_country, pickup_address, pickup_postal, pickup_state, pickup_city, pickup_province, due_date, completed) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", newOrder.AccountId, newOrder.OrderLength, newOrder.OrderWidth, newOrder.OrderHeight, newOrder.OrderWeight, newOrder.ConsigneeName, newOrder.ConsigneeNumber, newOrder.ConsigneeCountry, newOrder.ConsigneeAddress, newOrder.ConsigneePostal, newOrder.ConsigneeState, newOrder.ConsigneeCity, newOrder.ConsigneeProvince, newOrder.ConsigneeEmail, newOrder.PickupContactName, newOrder.PickupContactNumber, newOrder.PickupCountry, newOrder.PickupAddress, newOrder.PickupPostal, newOrder.PickupState, newOrder.PickupCity, newOrder.PickupProvince, newOrder.DueDate, newOrder.Completed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to create order"})
		return
	}

	// Respond
	c.IndentedJSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "New Order Successfully Created", "newOrderCreated": newOrder})
}

func getOrders(c *gin.Context) {
	var orders []order

	// Get rows of orders from DB
	rows, err := db.Query("SELECT * FROM orders")
	// if err from getting rows of orders from DB, return HTTP Bad Request 400
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to retrieve orders from DB"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var currentOrder order
		// scan each row of order and save to currentOrder
		if err := rows.Scan(
			&currentOrder.OrderId,
			&currentOrder.AccountId,
			&currentOrder.OrderLength,
			&currentOrder.OrderWidth,
			&currentOrder.OrderHeight,
			&currentOrder.OrderWeight,
			&currentOrder.ConsigneeName,
			&currentOrder.ConsigneeNumber,
			&currentOrder.ConsigneeCountry,
			&currentOrder.ConsigneeAddress,
			&currentOrder.ConsigneePostal,
			&currentOrder.ConsigneeState,
			&currentOrder.ConsigneeCity,
			&currentOrder.ConsigneeProvince,
			&currentOrder.ConsigneeEmail,
			&currentOrder.PickupContactName,
			&currentOrder.PickupContactNumber,
			&currentOrder.PickupCountry,
			&currentOrder.PickupAddress,
			&currentOrder.PickupPostal,
			&currentOrder.PickupState,
			&currentOrder.PickupCity,
			&currentOrder.PickupProvince,
			&currentOrder.DueDate,
			&currentOrder.Completed); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to save orders from DB"})
			return
		}
		// add currentOrder to orders slice
		orders = append(orders, currentOrder)
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "successfully retrieved orders from DB", "orders": orders})
}

func updateOrderStatus(c *gin.Context) {
	var reqBody orderIdCompleted

	// Returns Error HTTP Bad Request 400 if unable to read from request body
	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to read request body"})
		return
	}

	// Find order based on Order ID and update status
	rows, err := db.Query("UPDATE orders SET completed=? WHERE order_id=?", reqBody.Completed, reqBody.OrderId)
	// if err in updating order, return HTTP Bad Request 400
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to update orders in database"})
		return
	} else if rows != nil {
		// Respond
		c.IndentedJSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Order Updated Successfully", "orderUpdated": reqBody.OrderId})
	}
}

func assignOrder(c *gin.Context) {
	var reqBody orderIdAccountId

	// Returns Error HTTP Bad Request 400 if unable to read from request body
	if c.BindJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to read request body"})
		return
	}

	// Find order based on Order ID and update status
	rows, err := db.Query("UPDATE orders SET account_id=? WHERE order_id=?", reqBody.AccountId, reqBody.OrderId)
	// if err in updating order, return HTTP Bad Request 400
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Failed to assign orders in database"})
		return
	} else if rows != nil {
		// Respond
		c.IndentedJSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Order Assigned Successfully", "orderUpdated": reqBody.OrderId})
	}
}
