package services

import (
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// Defines who should process the request
func (s *Service) Proxy(w http.ResponseWriter, r *http.Request) {
	if websocket.IsWebSocketUpgrade(r) {
		s.wsProxy(w, r)
	} else {
		s.httpProxy(w, r)
	}
}

// Http proxy handler
func (s *Service) httpProxy(w http.ResponseWriter, r *http.Request) {
	proxyUrl, err := s.getServiceUrl(r.Host, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	parsed, err := url.Parse(proxyUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	r.URL.Scheme = parsed.Scheme
	r.URL.Host = parsed.Host

	proxy := &http.Transport{
		Proxy: http.ProxyURL(parsed),
	}
	proxyReq, err := proxy.RoundTrip(r)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer proxyReq.Body.Close()

	copyHeader(w.Header(), proxyReq.Header)
	w.WriteHeader(proxyReq.StatusCode)
	io.Copy(w, proxyReq.Body)
}

// Websockets proxy handler
func (s *Service) wsProxy(w http.ResponseWriter, r *http.Request) {
	proxyUrl, err := s.getServiceUrl(r.Host, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	u := url.URL{Scheme: "ws", Host: proxyUrl[7:], Path: r.URL.Path}

	connBackend, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		http.Error(w, "Could not connect to backend", http.StatusInternalServerError)
		return
	}
	defer connBackend.Close()

	upgrader := websocket.Upgrader{}
	connClient, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer connClient.Close()

	go func() {
		for {
			_, msg, err := connBackend.ReadMessage()
			if err != nil {
				return
			}
			connClient.WriteMessage(websocket.TextMessage, msg)
		}
	}()

	for {
		_, msg, err := connClient.ReadMessage()
		if err != nil {
			return
		}
		connBackend.WriteMessage(websocket.TextMessage, msg)
	}
}

// Copyes headers to new response
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
