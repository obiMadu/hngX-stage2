package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//main func
	//connect to sql database
	db, err := sql.Open("sql", "conn-string")
	if err != nil {
		fmt.Println("Error validating sql.Open arguments")
		panic(err.Error())
	}
	defer db.Close() //close connection. Best practice?

	//verify connection to database
	err = db.Ping()
	if err != nil {
		fmt.Print("Unable verify connection to DB with db.Ping")
		panic(err.Error())
	}

	fmt.Println("Connection to database succesful!")

	//create url handlers
	// http.HandleFunc("/api", createUser)

	//start http server
	// http.ListenAndServe(":8080", nil)

}

// func createUser(w *http.ResponseWriter, r *http.Request) {
// 	//
// }
