package handling

import (
	"net/http"

	ta "github.com/victorspringer/trapAdvisor/touristattraction"
	"github.com/victorspringer/trapAdvisor/trip"
)

// Service is the interface that provides handling methods.
type Service interface {
	StoreTrip(http.ResponseWriter, *http.Request)
	StoreTouristAttraction(http.ResponseWriter, *http.Request)
	FindTrip(http.ResponseWriter, *http.Request)
	FindTripByTravellerID(http.ResponseWriter, *http.Request)
	FindTouristAttraction(http.ResponseWriter, *http.Request)
	FindTouristAttractionByTripID(http.ResponseWriter, *http.Request)
	FindTouristAttractionByNamePart(http.ResponseWriter, *http.Request)
	FindMostVisitedTouristAttractions(http.ResponseWriter, *http.Request)
	FindBestRatedTouristAttractions(http.ResponseWriter, *http.Request)
}

type service struct {
	tripRepo trip.Repository
	taRepo   ta.Repository
}

// NewService creates a handling service with necessary dependencies.
func NewService(tr trip.Repository, tar ta.Repository) Service {
	return &service{tripRepo: tr, taRepo: tar}
}
