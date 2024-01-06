package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"sms_gateway/broker"
	"sms_gateway/logger"
)

type RequestHandler struct {
	broker *broker.Activemq
	logger *logger.FileLogger
	host   string
	ip     string
	port   int
}

func NewRequestHandler() *RequestHandler {

	return &RequestHandler{
		broker: broker.NewMessageBroker(),
		logger: logger.GetLumberJack(),
		ip:     "localhost",
		port:   8080,
	}

}

func (h *RequestHandler) WithPort(port int) *RequestHandler {

	h.port = port
	return h
}

func (h *RequestHandler) WithHost(host string) *RequestHandler {

	h.host = host
	return h
}

func (h *RequestHandler) Request(w http.ResponseWriter, r *http.Request) {

	queryParametersURI := r.URL.RawQuery
	decodedString, _ := url.PathUnescape(queryParametersURI)
	queryParametersURI = url.PathEscape(decodedString)

	h.broker.Send("http_calls", queryParametersURI)
	h.logger.WriteLog("|HTTP|%s", queryParametersURI)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("001 OK"))
}

func (h *RequestHandler) Listen() {

	http.HandleFunc("/insms", h.Request)
	http.ListenAndServe(fmt.Sprintf("%s:%d", h.ip, h.port), nil)

}
