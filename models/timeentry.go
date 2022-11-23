package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"


	_ "github.com/go-sql-driver/mysql"
)

type Timeentry struct {
	ID        		int 	`json:"id"`
	RegistrationID	int		`json:"registration_id"`
	TimeIn  		string	`json:"time_in"`
	TimeOut			string	`json:"time_out"`
}

func GetTimeentrys() []*Timeentry {

	var timeentrys []*Timeentry

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  results, err := db.Query("SELECT id, registration_id, time_in, time_out FROM timeentry")
	  if err != nil {
		panic(err.Error())
	  }

	for results.Next() {
		var timeentry Timeentry

		err = results.Scan(&timeentry.ID, &timeentry.RegistrationID, &timeentry.TimeIn, &timeentry.TimeOut)
		if err != nil {
		  panic(err.Error())
		}

		timeentrys = append(timeentrys, &timeentry)

		}

	  return timeentrys

}

func AddTimeentry(u Timeentry) (Timeentry, error) {

	if u.ID != 0 {
		return Timeentry{}, errors.New("new Timeentry must not include id or it must be")
	}

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("INSERT INTO timeentry (registration_id) VALUES (%d);", u.RegistrationID)
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

func GetTimeentryByID(id int) (Timeentry, error) {

	var timeentry Timeentry

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")
	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  stmt := fmt.Sprintf("SELECT id, registration_id, time_in, time_out  FROM timeentry WHERE id = %d", id)
	  result := db.QueryRow(stmt)
	  err = result.Scan(&timeentry.ID, &timeentry.RegistrationID, &timeentry.TimeIn, &timeentry.TimeOut)
	  
	  switch {
		case err == sql.ErrNoRows:
			return Timeentry{}, fmt.Errorf("Timeentry with ID '%d' not found",id)
		case err != nil:
			panic(err.Error())		
	 	 }

	return timeentry, nil

}

func UpdateTimeentry(u Timeentry) (Timeentry, error) {

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("update timeentry set time_in = \"%s\", time_out = \"%s\" where id = %d;", u.TimeIn, u.TimeOut, u.ID)

	  fmt.Println(sql)

	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  count, err := results.RowsAffected()
		if err != nil {
		panic(err.Error())
	  }

	  if count == 0 {
		return Timeentry{}, fmt.Errorf("No updates made to TimeEntry with ID '%v'", u.ID)
	  }

	return u, nil
}

func RemoveTimeentryByID(id int) error {

	db, err := sql.Open("mysql", "root:mypass@tcp(172.18.0.2:3306)/volunteer?parseTime=true")

	if err != nil {
		log.Print(err.Error())
	  }
	  defer db.Close()

	  sql := fmt.Sprintf("delete from timeentry where id = %d;", id)

	  results, err := db.Exec(sql)
	  if err != nil {
		panic(err.Error())
	  }

	  count, err := results.RowsAffected()
		if err != nil {
		panic(err.Error())
	  }

	  
	  if count == 0 {
		return fmt.Errorf("Timeentry with ID '%v' not found", id)
	  }

	return nil
}