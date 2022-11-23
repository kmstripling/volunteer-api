package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterControllers() {
	uc := newUserController()
	pc := newPersonController()
	ec := newEventController()
	rc := newRegistrationController()
	tc := newTimeentryController()

	http.Handle("/users", *uc)
	http.Handle("/users/", *uc)
	http.Handle("/persons", *pc)
	http.Handle("/persons/", *pc)
	http.Handle("/events", *ec)
	http.Handle("/events/", *ec)
	http.Handle("/registrations", *rc)
	http.Handle("/registrations/", *rc)
	http.Handle("/timeentries", *tc)
	http.Handle("/timeentries/", *tc)
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	
	enc := json.NewEncoder(w)
	enc.Encode(data)
}