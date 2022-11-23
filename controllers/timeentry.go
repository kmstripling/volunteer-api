package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"github.com/pluralsight/webservice/models"
)

type timeentryController struct {
	timeentryIDPattern *regexp.Regexp
}



func (uc timeentryController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //Temporary
	w.Header().Set("Access-Control-Allow-Headers", "*") //Temporary
	w.Header().Set("Access-Control-Allow-Methods" , "*") //Temporary

	if r.URL.Path == "/timeentries"{
		switch r.Method {
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
				uc.getAll(w, r)
		case http.MethodPost:
		 		uc.post(w, r)
		default:
				w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := uc.timeentryIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}

		id, err := strconv.Atoi(matches[1])
		
		if err != nil{
			w.WriteHeader(http.StatusNotFound)
		}
		switch r.Method{
		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			uc.get(id, w)
		case http.MethodPut:
		 	uc.put(id, w, r)
		case http.MethodDelete:
			uc.delete(id, w)
		default:
			w.WriteHeader(http.StatusNotImplemented)

		}
	}

}

func (uc *timeentryController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetTimeentrys(), w)
}

func (uc *timeentryController) get(id int, w http.ResponseWriter) {
	u, err := models.GetTimeentryByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(u, w)
}

func (uc *timeentryController) post(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse timeentry object"))
		return
	}

	u, err = models.AddTimeentry(u)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	encodeResponseAsJSON(u, w)
}

func (uc *timeentryController) put(id int, w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse timeentry Object"))
		return
	}
	if id != u.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted timeentry must match ID in URL"))
		return
	}
	u, err = models.UpdateTimeentry(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(u, w)

}

func (uc *timeentryController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveTimeentryByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (uc *timeentryController) parseRequest(r *http.Request) (models.Timeentry, error) {

	fmt.Println("Parsing");

	dec := json.NewDecoder(r.Body)
	var p models.Timeentry
	err := dec.Decode(&p)
	if err != nil {
		return models.Timeentry{}, err
	}

	fmt.Println(p.ID)
	fmt.Println("TimeIn: " + p.TimeIn)
	fmt.Println("TimeOut: " + p.TimeOut)

	return p, nil
}

func newTimeentryController() *timeentryController {
	return &timeentryController{
		timeentryIDPattern: regexp.MustCompile(`^/timeentries/(\d+)/?`),
	}
}
