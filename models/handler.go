package models

import "context"

type HandlerIFace interface {
	GetRoutes() []Route
	GetService() Service
	GetPrefix() string
}

type Handler struct {
	HandlerIFace
	Routes  []Route
	Prefix  string
	Service Service
	Context context.Context
}
