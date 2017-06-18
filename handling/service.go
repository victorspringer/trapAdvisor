package handling

import (
	"net/http"

	"github.com/victorspringer/trapAdvisor/friendship"
	ta "github.com/victorspringer/trapAdvisor/touristattraction"
	"github.com/victorspringer/trapAdvisor/traveller"
	"github.com/victorspringer/trapAdvisor/trip"
)

// Service is the interface that provides handling methods.
type Service interface {
	Health(http.ResponseWriter, *http.Request)
	StoreTrip(http.ResponseWriter, *http.Request)
	StoreTouristAttraction(http.ResponseWriter, *http.Request)
	FindTraveller(http.ResponseWriter, *http.Request)
	FindFriendshipByTravellerID(http.ResponseWriter, *http.Request)
	FindTrip(http.ResponseWriter, *http.Request)
	FindTripByTravellerID(http.ResponseWriter, *http.Request)
	FindTouristAttraction(http.ResponseWriter, *http.Request)
	FindTouristAttractionByTripID(http.ResponseWriter, *http.Request)
	FindTouristAttractionByNamePart(http.ResponseWriter, *http.Request)
	FindMostVisitedTouristAttractions(http.ResponseWriter, *http.Request)
	FindBestRatedTouristAttractions(http.ResponseWriter, *http.Request)
}

type service struct {
	travRepo traveller.Repository
	fRepo    friendship.Repository
	tripRepo trip.Repository
	taRepo   ta.Repository
}

// NewService creates a handling service with necessary dependencies.
func NewService(
	travr traveller.Repository,
	fr friendship.Repository,
	tr trip.Repository,
	tar ta.Repository,
) Service {
	return &service{
		travRepo: travr,
		fRepo:    fr,
		tripRepo: tr,
		taRepo:   tar,
	}
}
