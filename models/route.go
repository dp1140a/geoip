package models

import "net/http"

type Route struct {
	Name        string           `json:"name"`
	Method      string           `json:"method"`
	Pattern     string           `json:"pattern"`
	Protected   bool             `json:"protected"`
	HandlerFunc http.HandlerFunc `json:"-"`
}
