package main

import (
	// "database/sql"
	// "log"
	// _ "github.com/go-sql-driver/mysql"
	"net/http"

	"github.com/pluralsight/webservice/controllers"
)

// type Person struct {
// 	id       int    `json:"id"`
// 	firstname string `json:"firstname"`
// 	lastname string `json:"lastname"`
//   }

func main() {

	// db, err := sql.Open("mysql", "root:mypass@tcp(localhost:3306)/volunteer?parseTime=true")

	// if err != nil {
	// 	log.Print(err.Error())
	//   }
	//   defer db.Close()

	//   results, err := db.Query("SELECT id, firstname, lastname FROM person")
	//   if err != nil {
	// 	panic(err.Error())
	//   }

	// for results.Next() {
	// 	var person Person
	// 	err = results.Scan(&person.id, &person.firstname, &person.lastname)
	// 	if err != nil {
	// 	  panic(err.Error())
	// 	}
	// 	log.Printf(person.firstname + " " + person.lastname)
	//   }

	controllers.RegisterControllers()
	http.ListenAndServe(":3000", nil)

}
