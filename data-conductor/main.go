package main

import (
	"fmt"
	"log"
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"io/ioutil"
	"github.com/wfchiang/davic"
)

var LIST_OPT_DAVIC = []interface{}{}
var DATA_HEROS = map[string]interface{}{}
var DATA_POWERS = map[string]interface{}{}

/*
Utilities
*/ 
func loadData (file_path string) map[string]interface{} {
	raw_data, err := ioutil.ReadFile(file_path)
	if (err != nil) {
		panic(fmt.Sprintf("Failed to load data file %v : %v", file_path, err))
	}

	var data map[string]interface{}

	err = json.Unmarshal(raw_data, &data)
	if (err != nil) {
		panic(fmt.Sprintf("Data unmarshaling failed: %v", err))
	}

	return data
}

func recoverFromPanic (http_resp http.ResponseWriter, id_service string) {
	if r := recover(); r != nil {
		err_message := fmt.Sprintf("%v", r)
		log.Println(fmt.Sprintf("[%s] %s", id_service, err_message))
		http_resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(http_resp, err_message)
	}
}

/* 
Mock REST services
*/ 
func getHeroHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	parameters := http_reqt.URL.Query()
	id, err := strconv.Atoi(parameters.Get("id"))
	if (err != nil) {
		panic(fmt.Sprintf("Invalid id %v", err))
	}

	list_heros := davic.CastInterfaceToArray(DATA_HEROS["data"])
	found_hero := false 

	for _, v := range list_heros {
		hero := davic.CastInterfaceToObj(v)
		hero_id := davic.CastInterfaceToNumber(hero["id"])
		if (hero_id == float64(id)) {
			string_hero := string(davic.MarshalInterfaceToBytes(hero))
			fmt.Fprintf(http_resp, string_hero)
			found_hero = true
			break
		}
	}

	if (!found_hero) {
		fmt.Fprintf(http_resp, "{}")
	}
}

func getPowerHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	parameters := http_reqt.URL.Query()
	id, err := strconv.Atoi(parameters.Get("id"))
	if (err != nil) {
		panic(fmt.Sprintf("Invalid id %v", err))
	}

	list_powers := davic.CastInterfaceToArray(DATA_POWERS["data"])
	found_power := false
	
	for _, v := range list_powers {
		power := davic.CastInterfaceToObj(v)
		power_id := davic.CastInterfaceToNumber(power["id"])
		if (power_id == float64(id)) {
			string_power := string(davic.MarshalInterfaceToBytes(v))
			fmt.Fprintf(http_resp, string_power)
			found_power = true 
		}
	}

	if (!found_power) {
		fmt.Fprintf(http_resp, "{}")
	}
}

/* 
Davic handlers 
*/
func davicUnsetHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "davic/unset")

	if (len(LIST_OPT_DAVIC) > 0) {
		popped_opt := davic.CastInterfaceToArray(LIST_OPT_DAVIC[len(LIST_OPT_DAVIC)-1])
		LIST_OPT_DAVIC = LIST_OPT_DAVIC[0:len(LIST_OPT_DAVIC)-1]
		log.Println(fmt.Sprintf("[davic unset] Popped Operation: %v", popped_opt))
	}
	
	log.Println("Davic/Unset is Hit!")
}

func davicSetHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "davic/set")

	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		panic("Failed to read the request body")
	}

	var opt_davic []interface{}
	json.Unmarshal(bytes_reqt_body, &opt_davic)

	LIST_OPT_DAVIC = append(LIST_OPT_DAVIC, opt_davic)
	log.Println("Davic/Set is Hit!")

	log.Println(fmt.Sprintf("[davic set] Operation List: %v", LIST_OPT_DAVIC))
}

func davicGoHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "davic/go")

	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		panic("Failed to read the request body")
	}

	// Initialize the Davic environment 
	reqt_obj := davic.CreateObjFromBytes(bytes_reqt_body)
	env := davic.CreateNewEnvironment()
	env.Store = reqt_obj

	// Prepare the operation 
	env = davic.Execute(env, LIST_OPT_DAVIC)
	davic_result := env.Store
	log.Println(fmt.Sprintf("[davic] Result: %v", davic_result))

	// Acquire the result 
	resp_body, err := json.Marshal(davic_result) 
	if err != nil {
		panic(fmt.Sprintf("Response marshalling failed: %v", err))
	} 
	
	fmt.Fprintf(http_resp, string(resp_body))
}

/* 
Main 
*/ 
func main () {
	// Load data 
	log.Println("Load Data to Cache")
	DATA_HEROS = loadData("./data/heros.json")
	DATA_POWERS = loadData("./data/powers.json")

	// Setup routes 
	log.Println("Starting data-conductor")
	mux_router := mux.NewRouter()
	mux_router.HandleFunc("/get-hero", getHeroHandler).Methods("GET")
	mux_router.HandleFunc("/get-power", getPowerHandler).Methods("GET")
	mux_router.HandleFunc("/davic/set", davicSetHandler).Methods("POST");
	mux_router.HandleFunc("/davic/unset", davicUnsetHandler).Methods("GET");
	mux_router.HandleFunc("/davic/go", davicGoHandler).Methods("POST");

	http.Handle("/", mux_router)

	// Start the server 
	log.Fatal(http.ListenAndServe(":8080", nil))
}