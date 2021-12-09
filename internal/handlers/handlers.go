package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/trcwrx78/bookings/internal/config"
	"github.com/trcwrx78/bookings/internal/driver"
	"github.com/trcwrx78/bookings/internal/forms"
	"github.com/trcwrx78/bookings/internal/helpers"
	"github.com/trcwrx78/bookings/internal/models"
	"github.com/trcwrx78/bookings/internal/render"
	"github.com/trcwrx78/bookings/internal/repository"
	"github.com/trcwrx78/bookings/internal/repository/dbrepo"
)

// The repository used by our handlers
var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// Creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// Sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Renders the make a reservation page
func(m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// Handles the posting of a registration form
func(m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}
	
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Email: r.Form.Get("email"),
		Phone: r.Form.Get("phone"),
		StartDate: startDate,
		EndDate: endDate,
		RoomID: roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		 data := make(map[string]interface{})
		 data["reservation"] = reservation

		 render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}	

	err = m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Renders the a room page
func(m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Renders the a room page
func(m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Renders the search availabilty page
func(m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availabilty.page.tmpl", &models.TemplateData{})
}

// Handles a post to the search-availability
func(m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")



	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK bool `json:"ok"`
	Message string `json:"message"`
	
}
// Handles request for avaialbility and send JSON response
func(m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK: false,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Renders the search availabilty page
func(m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Renders the search availabilty page
func(m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}