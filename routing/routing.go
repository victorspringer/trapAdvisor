package routing

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/victorspringer/trapAdvisor/authenticating"
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

	aSvc := authenticating.NewService()

	hSvc := handling.NewService(
		persistence.NewTravellerRepository(),
		persistence.NewFriendshipRepository(),
		persistence.NewTripRepository(),
		persistence.NewTouristAttractionRepository(),
	)

	routes := []route{
		route{"GET", "/health", "Health", hSvc.Health},
		route{"GET", "/login", "Login", aSvc.HandleFacebookLogin},
		route{"GET", "/auth_callback", "AuthCallback", aSvc.HandleFacebookCallback},
		route{"GET", "/logout", "Logout", aSvc.HandleFacebookLogout},

		route{"POST", "/v1/trip/store", "StoreTrip", aSvc.AuthMiddleware(hSvc.StoreTrip)},
		route{"POST", "/v1/ta/store", "StoreTouristAttraction", aSvc.AuthMiddleware(hSvc.StoreTouristAttraction)},

		route{"GET", "/v1/traveller/find/{id}", "FindTraveller", aSvc.AuthMiddleware(hSvc.FindTraveller)},
		route{"GET", "/v1/friendship/find/traveller/{id}", "FindFriendshipByTravellerID", aSvc.AuthMiddleware(hSvc.FindFriendshipByTravellerID)},
		route{"GET", "/v1/trip/find/{id}", "FindTrip", aSvc.AuthMiddleware(hSvc.FindTrip)},
		route{"GET", "/v1/trip/find/traveller/{id}", "FindTripByTravellerID", aSvc.AuthMiddleware(hSvc.FindTripByTravellerID)},
		route{"GET", "/v1/ta/find/{id}", "FindTouristAttraction", aSvc.AuthMiddleware(hSvc.FindTouristAttraction)},
		route{"GET", "/v1/ta/find/trip/{id}", "FindTouristAttractionByTripID", aSvc.AuthMiddleware(hSvc.FindTouristAttractionByTripID)},
		route{"GET", "/v1/ta/find/name_part/{namePart}", "FindTouristAttractionByNamePart", aSvc.AuthMiddleware(hSvc.FindTouristAttractionByNamePart)},
		route{"GET", "/v1/ta/most_visited", "FindMostVisitedTouristAttractions", aSvc.AuthMiddleware(hSvc.FindMostVisitedTouristAttractions)},
		route{"GET", "/v1/ta/best_rated", "FindBestRatedTouristAttractions", aSvc.AuthMiddleware(hSvc.FindBestRatedTouristAttractions)},
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = logger(handler, route.Name)
		handler = cors.Default().Handler(handler)

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	return router
}

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf("%v\t%v\t%v\t%v", r.Method, r.RequestURI, name, time.Since(start))
	})
}
