package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"github.com/wfchiang/davic"
)

var OPT_HI2YOU = []interface{}{}
var OPT_DAVIC = []interface{}{}

const (
	KEY_STORE_HTTP_REQUEST  = "http-reqt"
	KEY_STORE_HTTP_RESPONSE = "http-resp"
)

// ==== 
// Recovery function 
// ====
func recoverFromPanic (http_resp http.ResponseWriter, id_service string) {
	if r := recover(); r != nil {
		err_message := fmt.Sprintf("%v", r)
		log.Println(fmt.Sprintf("[%s] %s", id_service, err_message))
		fmt.Fprintf(http_resp, err_message)
	}
}

// ==== 
// Variable initializations 
// ====
func initOperations () {
	// Hi2You
	resp_hi2you := map[string]interface{}{"name":nil, "message":"Hi!"}
	opt_resp_init := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_STORE_WRITE, KEY_STORE_HTTP_RESPONSE, resp_hi2you}
	opt_reqt_read := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_STORE_READ, KEY_STORE_HTTP_REQUEST}
	opt_resp_read := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_STORE_READ, KEY_STORE_HTTP_RESPONSE}
	opt_obj_read_name := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_OBJ_READ, opt_reqt_read, []interface{}{"name"}}
	opt_obj_update_name := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_OBJ_UPDATE, opt_resp_read, []interface{}{"name"}, opt_obj_read_name}
	opt_resp_update := []interface{}{davic.SYMBOL_OPT_MARK, davic.OPT_STORE_WRITE, KEY_STORE_HTTP_RESPONSE, opt_obj_update_name}
	OPT_HI2YOU = []interface{}{opt_resp_init, opt_resp_update} 
}

// ====
// Handlers 
// ====
func homepageHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	log.Println("Homepage is Hit!")
	fmt.Fprintf(http_resp, "Rest-mocker is Here!")
}

func echoHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "echo")

	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		panic("Failed to read the request body")
	}

	reqt_body := string(bytes_reqt_body)
	log.Println("Echo is Hit! Body: " + reqt_body)

	fmt.Fprintf(http_resp, reqt_body)
}

func hi2youHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "hi2you")

	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		panic("Failed to read the request body")
	}

	reqt_body := string(bytes_reqt_body)
	log.Println("Hi2You is Hit! Body: " + reqt_body)

	// Initialize the Davic environment 
	reqt_obj := davic.CreateObjFromBytes(bytes_reqt_body)
	env := davic.CreateNewEnvironment()
	env.Store = map[string]interface{}{KEY_STORE_HTTP_REQUEST:reqt_obj}
	
	// Prepare the operation 
	log.Println("Execute OPT_HI2YOU")
	env = davic.Execute(env, OPT_HI2YOU)
	obj_resp := env.Store[KEY_STORE_HTTP_RESPONSE]

	resp_body, err := json.Marshal(obj_resp) 
	if err != nil {
		panic(fmt.Sprintf("Response marshalling failed: %v", err))
	} 
	
	fmt.Fprintf(http_resp, string(resp_body))
}

func davicSetHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "davic/set")

	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		panic("Failed to read the request body")
	}

	var draft_opt_davic []interface{}
	json.Unmarshal(bytes_reqt_body, &draft_opt_davic)

	OPT_DAVIC = draft_opt_davic
	log.Println("Davic/Set is Hit!")

	log.Println(fmt.Sprintf("[davic] Operation: %v", OPT_DAVIC))
}

func davicGoHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "davic/go")

	bytes_reqt_body, err := ioutil.ReadAll(http_reqt.Body)
	if err != nil {
		panic("Failed to read the request body")
	}

	reqt_body := string(bytes_reqt_body)
	log.Println("Davic/Go is Hit! Body: " + reqt_body)

	// Initialize the Davic environment 
	reqt_obj := davic.CreateObjFromBytes(bytes_reqt_body)
	env := davic.CreateNewEnvironment()
	env.Store = map[string]interface{}{KEY_STORE_HTTP_REQUEST:reqt_obj}

	// Prepare the operation 
	env = davic.Execute(env, OPT_DAVIC)
	davic_result := env.Store[KEY_STORE_HTTP_RESPONSE]
	log.Println(fmt.Sprintf("[davic] Result: %v", davic_result))

	resp_body, err := json.Marshal(davic_result) 
	if err != nil {
		panic(fmt.Sprintf("Response marshalling failed: %v", err))
	} 
	
	fmt.Fprintf(http_resp, string(resp_body))
}

func getoptHandler (http_resp http.ResponseWriter, http_reqt *http.Request) {
	defer recoverFromPanic(http_resp, "getopt")

	id_service := mux.Vars(http_reqt)["id"]

	log.Println(fmt.Sprintf("Getopt is Hit! service id: %v", id_service))

	if (strings.Compare(id_service, "hi2you") == 0) {
		resp_body, err := json.Marshal(OPT_HI2YOU)
		if err != nil {
			panic(fmt.Sprintf("Service [%v] has a bad operation...", id_service))
		}
		fmt.Fprintf(http_resp, string(resp_body))

	} else if (strings.Compare(id_service, "davic") == 0) {
		resp_body, err := json.Marshal(OPT_DAVIC)
		if err != nil {
			panic(fmt.Sprintf("Service [%v] has a bad operation...", id_service))
		}
		fmt.Fprintf(http_resp, string(resp_body))
	
	} else {
		panic(fmt.Sprintf("Service [%v] is not defined...", id_service))
	}
}

// ====
// Main
// ====
func main () {
	log.Println("Initialize operations...")
	initOperations()

	log.Println("Starting rest-mocker...")
	mux_router := mux.NewRouter()

	mux_router.HandleFunc("/", homepageHandler)
	mux_router.HandleFunc("/echo", echoHandler)
	mux_router.HandleFunc("/hi2you", hi2youHandler).Methods("POST");
	mux_router.HandleFunc("/davic/go", davicGoHandler).Methods("POST");
	mux_router.HandleFunc("/davic/set", davicSetHandler).Methods("POST");
	mux_router.HandleFunc("/getopt/{id}", getoptHandler).Methods("GET"); 
	
	http.Handle("/", mux_router)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}