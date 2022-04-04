package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

// Our connection struct
type connection struct {
	host     string
	port     string
	user     string
	password string
	database string
	db       *sql.DB
	closed   bool
}

// Our response struct
type result struct {
	content  *sql.Rows
	complete bool
	error    error
}

/**
// *** README: What happens here depends on your use-case; you'll need to write that yourself ***
// Our parsed result struct
type parsedData struct {
}

// Presumably create & read our result struct
func (res result) read() parsedData {
	data := parsedData{}
	// Insert your code here
	return data
}
*/

// Send a query to the database
func (con connection) sendQuery(queryArg string) result {
	res, err := con.db.Query(queryArg)

	if err != nil {
		fmt.Println("Stopping application...\n- An error occurred during the '", queryArg, "' query")
		panic(err.Error())
	}

	fmtResult := result{content: res}
	fmt.Println("Attempting to network query response...")
	fmtResult.complete = true

	fmt.Println("Successful query handshake")
	return fmtResult
}

// Start our connection
func (con connection) create(ch chan connection) error {
	err := godotenv.Load("./config/.env")

	if err != nil {
		fmt.Println("Stopping application...")
		return errors.New("your config/.env folder failed to load")
	}

	envKeys := []string{"ENV_HOST", "ENV_PORT", "ENV_USER", "ENV_PASSWORD", "ENV_DATABASE"}

	for i := 0; i < len(envKeys); i++ {
		if os.Getenv(envKeys[i]) == "<INSERT HERE>" {
			fmt.Println("Stopping application...\nyour config/.env file is either unconfigured or you have not provided usable environment variables\"")
			return errors.New("failed to load needed environment variables")
		}
	}

	con.host = os.Getenv("ENV_HOST")
	con.port = os.Getenv("ENV_PORT")
	con.user = os.Getenv("ENV_USER")
	con.password = os.Getenv("ENV_PASSWORD")
	con.database = os.Getenv("ENV_DATABASE")

	conArg := con.user + ":" + con.password + "@tcp(" + con.host + ":" + con.port + ")/" + con.database
	db, err := sql.Open("mysql", conArg)

	if err != nil {
		fmt.Println("Please verify your database can be connected to")
		return errors.New("failed to connect to database")
	}

	fmt.Println("Successfully established connection")
	con.db = db

	ch <- con
	return nil
}

// End our connection
func (con connection) exit() error {
	_, err := con.db.Query("EXIT")

	if err != nil {
		fmt.Println("Stopping application...")
		return errors.New("failed to safely exit from our database connection")
	}

	con.closed = true
	return nil
}
