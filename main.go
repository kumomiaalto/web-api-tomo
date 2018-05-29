package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type AirlineTicket struct {
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Airline             string    `json:"airline"`
	Source              string    `json:"source"`
	Destination         string    `json:"destination"`
	FlightNumber        string    `json:"flight_number"`
	TicketNumber        string    `json:"ticket_number"`
	Seat                string    `json:"seat"`
	TicketClass         string    `json:"ticket_class"`
	BoardingTime        time.Time `json:"boarding_time"`
	DepartureTime       time.Time `json:"departure_time"`
	Terminal            string    `json:"terminal"`
	Gate                string    `json:"gate"`
	GateLocation        string    `json:"gate_location"`
	SecurityTimeMins    int       `json:"security_time"`
	ImmigrationTimeMins int       `json:"immigration_time"`
}

type PassengerLocation struct {
	TicketNumber string    `json:"ticket_number"`
	ReportType   string    `json:"report_type"`
	BeaconUUID   string    `json:"beacon_uuid"`
	BeaconMajor  string    `json:"beacon_major"`
	BeaconMinor  string    `json:"beacon_minor"`
	GeofenceID   string    `json:"geofence_id"`
	Latitude     string    `json:"latitude"`
	Longitude    string    `json:"longitude"`
	ReportTime   time.Time `json:"report_time"`
}

type Beacon struct {
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	Text          string `json:"text"`
	BeaconType    string `json:"beacon_type"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	NextBeacon    string `json:"next_beacon"`
	TimeToGate    string `json:"time_to_gate"`
	BoardingTime  string `json:"boarding_time"`
	DepartureTime string `json:"departure_time"`
	Route         string `json:"route"`
}

type BeaconRoute struct {
	Beacons []interface{} `json:"beacons"`
	Name    string        `json:"name"`
}

func getTicket(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var tickets []AirlineTicket
	q := datastore.NewQuery("AirlineTicket").Limit(1)
	if _, err := q.GetAll(ctx, &tickets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(tickets[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonBody)
}

func postTicket(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var ticket AirlineTicket
	err = json.Unmarshal(body, &ticket)
	if err != nil {
		panic(err)
	}

	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "AirlineTicket", ticket.TicketNumber, 0, nil)

	if _, err := datastore.Put(ctx, key, &ticket); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(w, "Something has gone wrong!")
		return
	}

	fmt.Println(w, "Updated successfully!")
}

func resetTicket(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var tickets []AirlineTicket
	q := datastore.NewQuery("AirlineTicket")
	if _, err := q.GetAll(ctx, &tickets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, ticket := range tickets {
		key := datastore.NewKey(ctx, "AirlineTicket", ticket.TicketNumber, 0, nil)
		err := datastore.Delete(ctx, key)
		if err != nil {
			panic(err)
		}
	}
}

func postLocation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var location PassengerLocation
	err = json.Unmarshal(body, &location)
	if err != nil {
		panic(err)
	}

	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "PassengerLocation", location.TicketNumber, 0, nil)

	if _, err := datastore.Put(ctx, key, &location); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(w, "Something has gone wrong!")
		return
	}

	fmt.Println(w, "Updated successfully!")
}

func getLocations(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var locations []PassengerLocation
	q := datastore.NewQuery("PassengerLocation").Order("-ReportTime")
	if _, err := q.GetAll(ctx, &locations); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(locations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonBody)
}

func resetPassengerLocations(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var locations []PassengerLocation
	q := datastore.NewQuery("PassengerLocation")
	if _, err := q.GetAll(ctx, &locations); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, location := range locations {
		key := datastore.NewKey(ctx, "PassengerLocation", location.TicketNumber, 0, nil)
		err := datastore.Delete(ctx, key)
		if err != nil {
			panic(err)
		}
	}
}

func postBeacon(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var beacon Beacon
	err = json.Unmarshal(body, &beacon)
	if err != nil {
		panic(err)
	}

	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "Beacon", beacon.Name, 0, nil)

	if _, err := datastore.Put(ctx, key, &beacon); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(w, "Something has gone wrong!")
		return
	}

	fmt.Println(w, "Updated successfully!")
}

func getAllBeacons(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var beacons []Beacon
	q := datastore.NewQuery("Beacon")
	if _, err := q.GetAll(ctx, &beacons); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(beacons)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonBody)
}

func resetBeacons(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var beacons []Beacon
	q := datastore.NewQuery("Beacon")
	if _, err := q.GetAll(ctx, &beacons); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, beacon := range beacons {
		key := datastore.NewKey(ctx, "Beacon", beacon.Name, 0, nil)
		err := datastore.Delete(ctx, key)
		if err != nil {
			panic(err)
		}
	}
}

func postBeaconRoute(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var beaconRoute BeaconRoute
	err = json.Unmarshal(body, &beaconRoute)
	if err != nil {
		panic(err)
	}

	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "BeaconRoute", beaconRoute.Name, 0, nil)

	if _, err := datastore.Put(ctx, key, &beaconRoute); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(w, "Something has gone wrong!")
		return
	}

	fmt.Println(w, "Updated successfully!")
}

func getAllBeaconRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var beaconRoutes []BeaconRoute
	q := datastore.NewQuery("BeaconRoute")
	if _, err := q.GetAll(ctx, &beaconRoutes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(beaconRoutes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonBody)
}

func resetBeaconRoutes(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var beaconRoutes []BeaconRoute
	q := datastore.NewQuery("BeaconRoute")
	if _, err := q.GetAll(ctx, &beaconRoutes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, beaconRoute := range beaconRoutes {
		key := datastore.NewKey(ctx, "Beacon", beaconRoute.Name, 0, nil)
		err := datastore.Delete(ctx, key)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	http.HandleFunc("/ticket/show", getTicket)
	http.HandleFunc("/ticket/update", postTicket)
	http.HandleFunc("/location/update", postLocation)
	http.HandleFunc("/locations/show", getLocations)
	http.HandleFunc("/beacons/show", getAllBeacons)
	http.HandleFunc("/beacon/update", postBeacon)
	http.HandleFunc("/beaconRoute/update", postBeaconRoute)
	http.HandleFunc("/beaconRoutes/show", getAllBeaconRoutes)
	http.HandleFunc("/reset/ticket", resetTicket)
	http.HandleFunc("/reset/locations", resetPassengerLocations)
	http.HandleFunc("/reset/beacons", resetBeacons)
	http.HandleFunc("/reset/beaconRoutes", resetBeaconRoutes)

	appengine.Main()
}
