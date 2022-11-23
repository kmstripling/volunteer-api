package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pluralsight/webservice/models"
)

type registrationController struct {
	registrationIDPattern *regexp.Regexp
}



func (uc registrationController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //Temporary
	w.Header().Set("Access-Control-Allow-Headers", "*") //Temporary
	w.Header().Set("Access-Control-Allow-Methods" , "*") //Temporary

	if r.URL.Path == "/registrations"{
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
		matches := uc.registrationIDPattern.FindStringSubmatch(r.URL.Path)
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

func (uc *registrationController) getAll(w http.ResponseWriter, r *http.Request) {

	event_id := 0;
	volunter_id := 0;

	param_event_id := r.URL.Query().Get("event_id")
	if param_event_id != "" {
		v_event_id, err := strconv.Atoi(param_event_id);
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not parse event_id"))
			return
		}
		event_id = v_event_id;
	}

	param_volunteer_id := r.URL.Query().Get("volunteer_id")
	if param_volunteer_id != "" {
		v_volunteer_id, err := strconv.Atoi(param_volunteer_id);
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not parse volunteer_id"))
			return
		}
		volunter_id = v_volunteer_id;
	} 
	
	encodeResponseAsJSON(models.GetRegistrations(event_id, volunter_id), w)
}

func (uc *registrationController) get(id int, w http.ResponseWriter) {
	u, err := models.GetRegistrationByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encodeResponseAsJSON(u, w)
}

func (uc *registrationController) post(w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse registration object"))
		return
	}

	u, err = models.AddRegistration(u)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	encodeResponseAsJSON(u, w)
}

func (uc *registrationController) put(id int, w http.ResponseWriter, r *http.Request) {
	u, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse registration Object"))
		return
	}
	if id != u.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID of submitted registration must match ID in URL"))
		return
	}
	u, err = models.UpdateRegistration(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(u, w)

}

func (uc *registrationController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveRegistrationByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (uc *registrationController) parseRequest(r *http.Request) (models.Registration, error) {

	dec := json.NewDecoder(r.Body)
	var p models.Registration
	err := dec.Decode(&p)
	if err != nil {
		return models.Registration{}, err
	}

	fmt.Println(p.VolunteerID)
	fmt.Println(p.EventID)

	return p, nil
}

func newRegistrationController() *registrationController {
	return &registrationController{
		registrationIDPattern: regexp.MustCompile(`^/registrations/(\d+)/?`),
	}
}
