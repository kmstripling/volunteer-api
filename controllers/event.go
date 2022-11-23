package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pluralsight/webservice/models"
)

type eventController struct {
	eventIDPattern *regexp.Regexp
}



func (uc eventController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //Temporary
	w.Header().Set("Access-Control-Allow-Headers", "*") //Temporary
	w.Header().Set("Access-Control-Allow-Methods" , "*") //Temporary

	if r.URL.Path == "/events"{
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
		matches := uc.eventIDPattern.FindStringSubmatch(r.URL.Path)
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

func (uc *eventController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetEvents(), w)
}

func (uc *eventController) get(id int, w http.ResponseWriter) {
	u, err := models.GetEventByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(u, w)
}

func (uc *eventController) post(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse event object"))
		return
	}

	u, err = models.AddEvent(u)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	encodeResponseAsJSON(u, w)
}

func (uc *eventController) put(id int, w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse event Object"))
		return
	}
	if id != u.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted event must match ID in URL"))
		return
	}
	u, err = models.UpdateEvent(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(u, w)

}

func (uc *eventController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveEventByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (uc *eventController) parseRequest(r *http.Request) (models.Event, error) {

	dec := json.NewDecoder(r.Body)
	var p models.Event
	err := dec.Decode(&p)
	if err != nil {
		return models.Event{}, err
	}

	return p, nil
}

func newEventController() *eventController {
	return &eventController{
		eventIDPattern: regexp.MustCompile(`^/events/(\d+)/?`),
	}
}
