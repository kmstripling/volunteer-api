package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pluralsight/webservice/models"
)

type personController struct {
	personIDPattern *regexp.Regexp
}



func (uc personController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //Temporary
	w.Header().Set("Access-Control-Allow-Headers", "*") //Temporary
	w.Header().Set("Access-Control-Allow-Methods" , "*") //Temporary

	if r.URL.Path == "/persons"{
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
		matches := uc.personIDPattern.FindStringSubmatch(r.URL.Path)
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

func (uc *personController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetPersons(), w)
}

func (uc *personController) get(id int, w http.ResponseWriter) {
	u, err := models.GetPersonByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(u, w)
}

func (uc *personController) post(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse person object"))
		return
	}

	u, err = models.AddPerson(u)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	encodeResponseAsJSON(u, w)
}

func (uc *personController) put(id int, w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse person Object"))
		return
	}
	if id != u.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted person must match ID in URL"))
		return
	}
	u, err = models.UpdatePerson(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(u, w)

}

func (uc *personController) delete(id int, w http.ResponseWriter) {
	err := models.RemovePersonByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (uc *personController) parseRequest(r *http.Request) (models.Person, error) {

	dec := json.NewDecoder(r.Body)
	var p models.Person
	err := dec.Decode(&p)
	if err != nil {
		return models.Person{}, err
	}

	return p, nil
}

func newPersonController() *personController {
	return &personController{
		personIDPattern: regexp.MustCompile(`^/persons/(\d+)/?`),
	}
}
