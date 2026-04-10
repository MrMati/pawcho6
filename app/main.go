package main

import (
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"os"
)

var version = "dev"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		ip := getServerIP()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Lab Docker Go</title>
  </head>
  <body>
    <main>
      <h1>Prosta aplikacja webowa</h1>
      <p><strong>Adres IP serwera:</strong> %s</p>
      <p><strong>Nazwa serwera:</strong> %s</p>
      <p><strong>Wersja aplikacji:</strong> %s</p>
    </main>
  </body>
</html>`, html.EscapeString(ip), html.EscapeString(hostname), html.EscapeString(version))
	})

	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		panic(err)
	}
}

func getServerIP() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("Failed to detect machine host name.", err.Error())
		return "127.0.0.1"
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		log.Println("Failed to lookup host IP.", err.Error())
		return "127.0.0.1"
	}

	fallback := ""
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			if fallback == "" {
				fallback = ipv4.String()
			}
			if ipv4[0] == 172 {
				return ipv4.String()
			}
		}
	}

	if fallback != "" {
		return fallback
	}

	log.Println("No IPv4 address found for host.", hostname)
	return "127.0.0.1"
}
