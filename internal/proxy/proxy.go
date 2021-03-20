package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type routeMap struct {
	Host 	string `json:"host"`
	Backend string `json:"backend"`
}

var IngressRule routeMap

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/routeMap", updateRouteMap)
	mux.HandleFunc("/", reverseProxy)

	log.Printf("Now serving simpleingress-simpleproxy\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func makeHTTPError(w http.ResponseWriter, r*http.Request, statusCode int , statustext, msg string) {
	if msg == "" {
		http.Error(w, statustext, statusCode)
	}
	http.Error(w, msg, statusCode)
	log.Printf("%d %s %s %s %s", statusCode, statustext, r.Method, r.Host, msg)
}

func updateRouteMap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		makeHTTPError(w,r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), "")
		return	
	}

	if r.Body == http.NoBody {
		makeHTTPError(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "missing request body")
		return
	}
	
	body, err := ioutil.ReadAll(r.Body)
	if body != nil {
		defer r.Body.Close()
	}
	if err != nil {
		makeHTTPError(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error())
		return
	}

	var rm routeMap
	err = json.Unmarshal(body, &rm)
	if err != nil {
		makeHTTPError(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error())
		return
	}
	IngressRule.Backend = rm.Backend

	log.Printf("%d %s %s %s", http.StatusOK, http.StatusText(http.StatusOK), r.Method, r.Host)
	fmt.Fprint(w, "ingress rule has been set successfully")
}

func reverseProxy(w http.ResponseWriter, r *http.Request) {
	if IngressRule.Backend == "" {
		http.Error(w, "proxy error: ingress rule not set", http.StatusBadRequest)
		log.Printf("%d %s %s %s %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest), r.Method, r.Host, "proxy error: ingress rule not set")
		return
	}
	target, _ := url.Parse(fmt.Sprintf("http://%s", IngressRule.Backend))

	proxy := httputil.NewSingleHostReverseProxy(target)		
	proxy.ServeHTTP(w, r)
	log.Printf("%s %s", r.Method, r.Host)
}
