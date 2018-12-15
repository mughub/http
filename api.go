package http

import (
	"encoding/json"
	"github.com/mughub/mughub/bare"
	"github.com/mughub/mughub/db"
	"github.com/spf13/viper"
	"net/http"
)

// RegisterAPIEndpoint registers the GraphQL API endpoint with the provided router.
func RegisterAPIEndpoint(r bare.Router, cfg *viper.Viper) {
	domain := cfg.GetString("domain")
	if domain == "" {
		panic("api domain must be provided")
	}

	api := r.Host(domain)
	api.Methods("GET", "POST").Path("/graphql")

	api.HandlerFunc(gqlHandler)
}

// req represents an incoming GraphQL request
type gqlReq struct {
	Query         string
	OperationName string
	Variables     map[string]interface{}
}

// getReq extracts a GraphQL request from a http.Request
func getReq(req *http.Request) (r *gqlReq, err error) {
	r = new(gqlReq)

	if req.Method == http.MethodGet {
		err = req.ParseForm()
		if err != nil {
			return
		}

		q := req.URL.Query()
		r.Query = q.Get("query")
		r.OperationName = q.Get("operationName")

		vStr := q.Get("variables")
		if vStr != "" {
			r.Variables = make(map[string]interface{})
			err = json.Unmarshal([]byte(vStr), r.Variables)
		}
	}

	if req.Method == http.MethodPost {
		err = json.NewDecoder(req.Body).Decode(r)
	}

	return
}

// gqlHandler handles GraphQL requests
func gqlHandler(w http.ResponseWriter, req *http.Request) {
	// Get GraphQL request from http.Request
	r, err := getReq(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute GraphQL request against database
	res := db.Do(req.Context(), r.Query, r.Variables)

	// Marshal GraphQL response
	b, err := res.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	_, err = w.Write(b)
	if err != nil {
		// TODO: Handle this write
	}
}
