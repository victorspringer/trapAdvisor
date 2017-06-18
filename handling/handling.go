package handling

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/victorspringer/trapAdvisor/persistence"
	ta "github.com/victorspringer/trapAdvisor/touristattraction"
	"github.com/victorspringer/trapAdvisor/trip"
)

func (s *service) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func (s *service) StoreTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var t trip.Trip

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = r.Body.Close(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, &t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	repo := persistence.NewTripRepository()
	if err = repo.Store(&t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *service) StoreTouristAttraction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var t ta.TouristAttraction

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = r.Body.Close(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(body, &t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	repo := persistence.NewTouristAttractionRepository()
	if err = repo.Store(&t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *service) FindTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo := persistence.NewTripRepository()
	t, err := repo.Find(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) FindTripByTravellerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo := persistence.NewTripRepository()
	t, err := repo.FindByTravellerID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) FindTouristAttraction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo := persistence.NewTouristAttractionRepository()
	t, err := repo.Find(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) FindTouristAttractionByTripID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo := persistence.NewTouristAttractionRepository()
	t, err := repo.FindByTripID(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) FindTouristAttractionByNamePart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	namePart := mux.Vars(r)["namePart"]

	repo := persistence.NewTouristAttractionRepository()
	t, err := repo.FindByNamePart(namePart)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) FindMostVisitedTouristAttractions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	repo := persistence.NewTouristAttractionRepository()
	t, err := repo.FindMostVisited()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *service) FindBestRatedTouristAttractions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	repo := persistence.NewTouristAttractionRepository()
	t, err := repo.FindBestRated()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
