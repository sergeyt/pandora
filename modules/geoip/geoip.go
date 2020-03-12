// borrowed from https://github.com/d4l3k/go-sct/blob/master/geoip/geoip.go
package geoip

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sergeyt/pandora/modules/auth"
	"github.com/sergeyt/pandora/modules/send"
)

// RegisterAPI add a GET route to request GeoIP
func RegisterAPI(r chi.Router) {
	r = r.With(auth.Middleware)

	r.Get("/api/geoip", func(w http.ResponseWriter, r *http.Request) {
		t, err := LookupIP(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		send.JSON(w, t)
	})
}

// GeoIP payload.
type GeoIP struct {
	// The right side is the name of the JSON variable
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
	AreaCode    int     `json:"area_code"`
}

var (
	address  string
	err      error
	geo      GeoIP
	response *http.Response
	body     []byte
)

// LookupIP looks up the geolocation information for the specified address ("" for current host).
func LookupIP(address string) (*GeoIP, error) {
	// Use freegeoip.net to get a JSON response
	// There is also /xml/ and /csv/ formats available
	response, err = http.Get("https://freegeoip.net/json/" + address)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, &geo)
	if err != nil {
		return nil, err
	}

	return &geo, nil
}
