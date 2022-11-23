package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"


	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	ID        		int 	`json:"id"`
	Name	 		string	`json:"name"`
	StartDatetime  	string	`json:"startdatetime"`
	EndDatetime		string	`json:"enddatetime"`
}

func GetEvents() []*Event {

	var events []*Event

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  results, err := db.Query("SELECT id, name, startdatetime, enddatetime FROM event")
	  if err != nil {
		panic(err.Error())
	  }

	for results.Next() {
		var event Event

		err = results.Scan(&event.ID, &event.Name, &event.StartDatetime, &event.EndDatetime)
		if err != nil {
		  panic(err.Error())
		}

		events = append(events, &event)

		}

	  return events

}

func AddEvent(u Event) (Event, error) {

	if u.ID != 0 {
		return Event{}, errors.New("new Event must not include id or it must be")
	}

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("INSERT INTO event (name, startdatetime, enddatetime) VALUES (\"%s\",\"%s\",\"%s\");", u.Name, u.StartDatetime, u.EndDatetime)
      
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

func GetEventByID(id int) (Event, error) {

	var event Event

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")
	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  stmt := fmt.Sprintf("SELECT id, name, startdatetime, enddatetime FROM event WHERE id = %d", id)
	  result := db.QueryRow(stmt)
	  err = result.Scan(&event.ID, &event.Name, &event.StartDatetime, &event.EndDatetime)
	  
	  switch {
		case err == sql.ErrNoRows:
			return Event{}, fmt.Errorf("Event with ID '%d' not found",id)
		case err != nil:
			panic(err.Error())		
	 	 }

	return event, nil

}

func UpdateEvent(u Event) (Event, error) {

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("update event set name = \"%s\", startdatetime = \"%s\", enddatetime = \"%s\" where id = %d;", u.Name, u.StartDatetime, u.EndDatetime, u.ID)

	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  count, err := results.RowsAffected()
		if err != nil {
		panic(err.Error())
	  }

	  if count == 0 {
		return Event{}, fmt.Errorf("No updates made to Event with ID '%v'", u.ID)
	  }

	return u, nil
}

func RemoveEventByID(id int) error {

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("delete from event where id = %d;", id)

	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  count, err := results.RowsAffected()
		if err != nil {
		panic(err.Error())
	  }

	  
	  if count == 0 {
		return fmt.Errorf("Event with ID '%v' not found", id)
	  }

	return nil
}