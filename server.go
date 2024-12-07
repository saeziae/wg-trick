package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/saeziae/wg-trick/config"
)

var configFilePath string

func getListOfRoute(config *config.WireGuardConfig) string {
	var listOfRoute []string
	listOfRoute = append(listOfRoute, config.Interface.Mask)
	for _, peer := range config.Peers {
		if peer.IsGateway {
			listOfRoute = append(listOfRoute, peer.AllowedIPs)
		}
	}
	return strings.Join(listOfRoute, ",")
}

func getClientIP(config *config.WireGuardConfig, publicKey string) string {
	for _, peer := range config.Peers {
		if peer.PublicKey == publicKey {
			return strings.Split(peer.AllowedIPs, "/")[0]
		}
	}
	return ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	pubKey := strings.Replace(chi.URLParam(r, "pubkey"), "_", "/", -1)
	configFile, err := config.ReadWireGuardConfig(configFilePath)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Fail to load configuration")
		log.Fatal(err)
		return
	}
	clientIP := getClientIP(configFile, pubKey)
	if clientIP == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Public key is not found")
		return
	}

	fmt.Print(configFile.Interface.ListenPort)
	clientPeer := config.Peer{
		PublicKey:           configFile.Interface.PublicKey,
		AllowedIPs:          getListOfRoute(configFile),
		Endpoint:            configFile.Interface.Endpoint + ":" + configFile.Interface.ListenPort,
		PersistentKeepalive: "25",
	}
	clientConfigini := config.GeneratePeerIni(clientPeer, clientIP)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Fail to generate configuration")
		return
	}
	w.Header().Set("Content-Type", "application/x-ini")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(clientConfigini))

}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/wg-trick/v1", func(r chi.Router) {
		r.Get("/{pubkey}", handler)
	})
	var listenAddr string
	flag.StringVar(&listenAddr, "l", "127.0.0.1:18964", "server listen address")
	flag.StringVar(&configFilePath, "c", "/etc/wireguard/wg0.conf", "path to config file")
	flag.Parse()
	log.Println("Starting server on ", listenAddr)
	err := http.ListenAndServe(listenAddr, r)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
