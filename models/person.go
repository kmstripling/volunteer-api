package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"


	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	ID        int 		`json:"id"`
	FirstName string	`json:"firstname"`
	LastName  string	`json:"lastname"`
}

func GetPersons() []*Person {

	var persons []*Person

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  results, err := db.Query("SELECT id, FirstName, LastName FROM person")
	  if err != nil {
		panic(err.Error())
	  }

	for results.Next() {
		var person Person

		err = results.Scan(&person.ID, &person.FirstName, &person.LastName)
		if err != nil {
		  panic(err.Error())
		}

		persons = append(persons, &person)

		}

	  return persons

}

func AddPerson(u Person) (Person, error) {

	if u.ID != 0 {
		return Person{}, errors.New("new Person must not include id or it must be")
	}

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("INSERT INTO person (FirstName, LastName) VALUES (\"%s\",\"%s\");", u.FirstName, u.LastName)
                      
	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  lastId, err := results.LastInsertId()
	  if err != nil {
		  log.Fatal(err)
	  }

	  u.ID = int(lastId)

	return u, nil
}

func GetPersonByID(id int) (Person, error) {

	var person Person

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")
	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  stmt := fmt.Sprintf("SELECT id, FirstName, LastName FROM person WHERE id = %d", id)
	  result := db.QueryRow(stmt)
	  err = result.Scan(&person.ID, &person.FirstName, &person.LastName)
	  
	  switch {
		case err == sql.ErrNoRows:
			return Person{}, fmt.Errorf("Person with ID '%d' not found",id)
		case err != nil:
			panic(err.Error())		
	 	 }

	return person, nil

}

func UpdatePerson(u Person) (Person, error) {

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("update person set firstname = \"%s\", lastname = \"%s\" where id = %d;", u.FirstName, u.LastName, u.ID)

	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  count, err := results.RowsAffected()
		if err != nil {
		panic(err.Error())
	  }

	  if count == 0 {
		return Person{}, fmt.Errorf("No updates made to Person with ID '%v'", u.ID)
	  }


	return u, nil
}

func RemovePersonByID(id int) error {

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("delete from person where id = %d;", id)

	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  count, err := results.RowsAffected()
		if err != nil {
		panic(err.Error())
	  }

	  
	  if count == 0 {
		return fmt.Errorf("Person with ID '%v' not found", id)
	  }

	return nil
}