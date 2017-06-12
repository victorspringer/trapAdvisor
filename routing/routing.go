package routing

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/victorspringer/trapAdvisor/handling"
	"github.com/victorspringer/trapAdvisor/persistence"
)

type route struct {
	Method      string
	Pattern     string
	Name        string
	HandlerFunc http.HandlerFunc
}

// Router initializer.
func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	hSvc := handling.NewService(
		persistence.NewTripRepository(), persistence.NewTouristAttractionRepository(),
	)

	routes := []route{
		route{"POST", "/v1/trip/store", "StoreTrip", hSvc.StoreTrip},
		route{"POST", "/v1/ta/store", "StoreTouristAttraction", hSvc.StoreTouristAttraction},
		route{"GET", "/v1/trip/find/{id}", "FindTrip", hSvc.FindTrip},
		route{"GET", "/v1/trip/find/traveller/{id}", "FindTripByTravellerID", hSvc.FindTripByTravellerID},
		route{"GET", "/v1/ta/find/{id}", "FindTouristAttraction", hSvc.FindTouristAttraction},
		route{"GET", "/v1/ta/find/trip/{id}", "FindTouristAttractionByTripID", hSvc.FindTouristAttractionByTripID},
		route{"GET", "/v1/ta/find/name_part/{namePart}", "FindTouristAttractionByNamePart", hSvc.FindTouristAttractionByNamePart},
		route{"GET", "/v1/ta/most_visited", "FindMostVisitedTouristAttractions", hSvc.FindMostVisitedTouristAttractions},
		route{"GET", "/v1/ta/best_rated", "FindBestRatedTouristAttractions", hSvc.FindBestRatedTouristAttractions},
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	return router
}

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}
