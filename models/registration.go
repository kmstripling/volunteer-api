package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"


	_ "github.com/go-sql-driver/mysql"
)

type Registration struct {
	ID        		int 	`json:"id"`
	RegDatetime		string	`json:"regdatetime"`
	EventID	 		int		`json:"event_id"`
	EventName		string	`json:"event_name"`
	EventStartDateTime	string `json:"event_startdatetime"`
	EventEndDateTime	string `json:"event_enddatetime"`     
	VolunteerID  	int		`json:"volunteer_id"`
	VolunteerFirstName	string	`json:"volunteer_firstname"`
	VolunteerLastName	string	`json:"volunteer_lastname"`
}

func GetRegistrations(event_id int, volunteer_id int) []*Registration {

	params := []string{}
	var registrations []*Registration

	//var sqlStr = "SELECT id, event_id, volunteer_id, regdatetime FROM registration %s"
	var sqlStr = "select r.id, r.regdatetime, e.id AS \"event_id\", e.name AS \"event_name\", e.startdatetime AS \"event_startdatetime\", e.enddatetime AS \"event_enddatetime\", p.id AS \"volunteer_id\", p.firstname AS \"volunteer_firstname\", p.lastname AS \"volunteer_lastname\" from registration r inner join person p on r.volunteer_id = p.id inner join event e on e.id = r.event_id %s"
	
	var whereClauseStr = ""

	if event_id != 0 {
	params = append(params, fmt.Sprintf("e.id = %d", event_id))
	}
	if volunteer_id != 0 {
	params = append(params, fmt.Sprintf("p.id = %d", volunteer_id))
	}

	for i, v := range params {

	if i == 0 {
	whereClauseStr = whereClauseStr + "where " + v
	} else if i > 0 {
	whereClauseStr = whereClauseStr + " and " + v
	}
	}

	log.Printf(sqlStr, whereClauseStr);

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  results, err := db.Query(fmt.Sprintf(sqlStr, whereClauseStr))
	  if err != nil {
		panic(err.Error())
	  } 

	for results.Next() {
		var registration Registration

		err = results.Scan(&registration.ID, &registration.RegDatetime, &registration.EventID, &registration.EventName, &registration.EventStartDateTime, &registration.EventEndDateTime, &registration.VolunteerID, &registration.VolunteerFirstName, &registration.VolunteerLastName)
		if err != nil {
		  panic(err.Error())
		}

		registrations = append(registrations, &registration)

		}

	  return registrations

}

func AddRegistration(u Registration) (Registration, error) {

	if u.ID != 0 {
		return Registration{}, errors.New("new Registration must not include id or it must be")
	}

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("INSERT INTO registration (event_id, volunteer_id) VALUES (%d,%d);", u.EventID, u.VolunteerID)
	  fmt.Println(sql)
      
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

func GetRegistrationByID(id int) (Registration, error) {

	var registration Registration

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")
	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  stmt := fmt.Sprintf("SELECT id, event_id, volunteer_id, regdatetime FROM registration WHERE id = %d", id)
	  result := db.QueryRow(stmt)
	  err = result.Scan(&registration.ID, &registration.EventID, &registration.VolunteerID, &registration.RegDatetime)
	  
	  switch {
		case err == sql.ErrNoRows:
			return Registration{}, fmt.Errorf("Registration with ID '%d' not found",id)
		case err != nil:
			panic(err.Error())		
	 	 }

	return registration, nil

}

func UpdateRegistration(u Registration) (Registration, error) {

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("update registration set event_id = %d, volunteer_id = %d where id = %d;", u.EventID, u.VolunteerID, u.ID)

	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  count, err := results.RowsAffected()
		if err != nil {
		panic(err.Error())
	  }

	  if count == 0 {
		return Registration{}, fmt.Errorf("No updates made to Registration with ID '%v'", u.ID)
	  }

	return u, nil
}

func RemoveRegistrationByID(id int) error {

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("delete from registration where id = %d;", id)

	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  count, err := results.RowsAffected()
		if err != nil {
		panic(err.Error())
	  }

	  
	  if count == 0 {
		return fmt.Errorf("Registration with ID '%v' not found", id)
	  }

	return nil
}