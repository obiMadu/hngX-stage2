package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type User struct {
	id        string
	slackname string
	fullname  string
	email     string
}

func main() {
	//main func
	//connect to sql database
	//define db conn string
	//connectionString := "username:password@tcp(ip:port)/db"
	connectionString := os.Getenv("CONN_STRING")
	if connectionString == "" {
		host := os.Getenv("MYSQL_HOST")
		username := os.Getenv("MYSQL_USERNAME")
		password := os.Getenv("MYSQL_PASSWORD")
		dbname := os.Getenv("MYSQL_DBNAME")

		connectionString = username + ":" + password + "@tcp(" + host + ")/" + dbname
	}
	fmt.Println(connectionString)
	var err error

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Error validating sql.Open arguments")
		panic(err.Error())
	}
	defer db.Close() //close connection. Best practice?

	//verify connection to database
	err = db.Ping()
	if err != nil {
		fmt.Println("Unable verify connection to DB with db.Ping()")
		panic(err.Error())
	}

	fmt.Println("Connection to database succesful!")

	//create url handlers
	route := mux.NewRouter()
	route.HandleFunc("/api", createUser)
	route.HandleFunc("/api/{slackname}", readHandler).Methods("GET")
	route.HandleFunc("/api/{slackname}", updateHandler).Methods("PATCH")
	route.HandleFunc("/api/{slackname}", putHandler).Methods("PUT")
	route.HandleFunc("/api/{slackname}", deleteHandler).Methods("DELETE")

	//start http server
	fmt.Println("*** Server Listening ***")
	http.ListenAndServe(":80", route)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// indicate /api path is running
	fmt.Println("*** CREATING A USER ***")

	// get and validate values from request
	slackname := r.FormValue("name")
	fullname := r.FormValue("fullname")
	email := r.FormValue("email")

	if slackname == "" {
		w.Header().Set("Content-Type", "text/plain") //set text header
		w.WriteHeader(206)
		fmt.Fprintf(w, "Slackname is a required parameter! User not created!!") //return error message
		return
	}

	insertStatement := "INSERT INTO `db`.`Users` (`slackname`, `fullname`, `email`) VALUES (?, ?, ?);"

	ins, err := db.Prepare(insertStatement)
	if err != nil {
		fmt.Println("Something went wrong preparing the sql statement", err)
		w.Header().Set("Content-Type", "text/plain")       //set text header
		fmt.Fprintf(w, "Error creating user. Check logs.") //return error message
		return
	}

	defer ins.Close()

	_, err = ins.Exec(slackname, fullname, email)
	if err != nil {
		fmt.Println("Error running DB insert statement", err)
		w.WriteHeader(206)
		w.Header().Set("Content-Type", "text/plain")                         //set text header
		fmt.Fprintf(w, "Error! Slackname has been taken. Try another name.") //return error message
		return
	}

	fmt.Fprintf(w, "User with slackname:"+slackname+" created succesfully!")
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	//running read handler
	fmt.Println("*** READING A USER ***")

	// Extract the slackname from the path.
	params := mux.Vars(r)
	slackname := params["slackname"]
	fmt.Println("Route variable: " + slackname)

	// run query
	query := "SELECT `id`,`slackname`,`fullname`,`email` FROM `db`.`Users` WHERE `slackname` = ?;"
	res := db.QueryRow(query, slackname)

	var user User
	err := res.Scan(&user.id, &user.slackname, &user.fullname, &user.email)
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("User does not exist"))
		return
	}

	fmt.Fprintln(w, "id: "+user.id)
	fmt.Fprintln(w, "slackname: "+user.slackname)
	fmt.Fprintln(w, "fullname: "+user.fullname)
	fmt.Fprintln(w, "email: "+user.email)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	//
	// indicate /api path is running
	fmt.Println("*** UPDATING A USER ***")

	// Extract the slackname from the path.
	params := mux.Vars(r)
	slackname0 := params["slackname"]
	fmt.Println("Route variable: " + slackname0)

	// get and validate values from request
	slackname := r.FormValue("name")
	fullname := r.FormValue("fullname")
	email := r.FormValue("email")

	//check if slackname exists in DB
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM `db`.`Users` WHERE `slackname` = ?);", slackname0).Scan(&exists)
	if err != nil {
		fmt.Println("Unable to determine if username exists", err.Error())
		fmt.Fprintf(w, "Unable to determine if username exists")
		return
	}
	//fmt.Println("exits", exists)
	if exists == false {
		w.Write([]byte("Username does not exist!"))
		return
	}

	if slackname == slackname0 || slackname == "" {
		updateStatement := "UPDATE `db`.`Users` SET `fullname`=?, `email`=? WHERE slackname=?;"

		ins, err := db.Prepare(updateStatement)
		if err != nil {
			fmt.Println("Something went wrong preparing the sql statement", err)
			w.Header().Set("Content-Type", "text/plain")       //set text header
			fmt.Fprintf(w, "Error creating user. Check logs.") //return error message
			return
		}

		_, err = ins.Exec(fullname, email, slackname0)

		if err != nil {
			fmt.Println("Error running DB update statement", err)
			w.WriteHeader(206)
			w.Header().Set("Content-Type", "text/plain") //set text header
			fmt.Fprintf(w, "Error! Creating update.")    //return error message
			return
		}

		defer ins.Close()

		fmt.Fprintf(w, "User updated succesfully!")

	} else {
		updateStatement := "UPDATE `db`.`Users` SET `slackname`=?, `fullname`=?, `email`=? WHERE slackname=?;"

		ins, err := db.Prepare(updateStatement)

		if err != nil {
			fmt.Println("Something went wrong preparing the sql statement", err)
			w.Header().Set("Content-Type", "text/plain")       //set text header
			fmt.Fprintf(w, "Error creating user. Check logs.") //return error message
			return
		}

		_, err = ins.Exec(slackname, fullname, email, slackname0)

		if err != nil {
			fmt.Println("Error running DB update statement", err)
			w.WriteHeader(501)
			w.Header().Set("Content-Type", "text/plain") //set text header
			fmt.Fprintf(w, "Error! updating user.")      //return error message
			return
		}

		defer ins.Close()
		fmt.Fprintf(w, "User with slackname:"+slackname0+" Updated to "+slackname+" succesfully!")

	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// indicate /api path is running
	fmt.Println("*** DELETING A USER ***")

	// Extract the slackname from the path.
	params := mux.Vars(r)
	slackname := params["slackname"]
	fmt.Println("Route variable: " + slackname)

	//check if slackname exists in DB
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM `db`.`Users` WHERE `slackname` = ?);", slackname).Scan(&exists)
	if err != nil {
		fmt.Println("Unable to determine if username exists", err.Error())
		fmt.Fprintf(w, "Unable to determine if username exists")
		return
	}
	//fmt.Println("exits", exists)
	if exists == false {
		w.Write([]byte("Username does not exist!"))
		return
	}

	deleteStatement := "DELETE FROM `db`.`Users` WHERE `slackname` = ?;"

	ins, err := db.Prepare(deleteStatement)
	if err != nil {
		fmt.Println("Something went wrong preparing the sql statement", err)
		w.Header().Set("Content-Type", "text/plain") //set text header
		fmt.Fprintf(w, "Error deleting user.")       //return error message
		return
	}

	defer ins.Close()

	_, err = ins.Exec(slackname)
	if err != nil {
		fmt.Println("Error running DB delete statement", err)
		w.WriteHeader(501)
		w.Header().Set("Content-Type", "text/plain") //set text header
		fmt.Fprintf(w, "Error deleting user.")       //return error message
		return
	}

	fmt.Fprintf(w, "User with slackname:"+slackname+" deleted succesfully!")
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Please use PATCH for updates.")
}
