// App de ejemplo para probar los tenants tipo "app" de kosmos: responde en PORT y escribe un
// contador en /data (el volumen), para comprobar que el FS de solo lectura y el PVC funcionan.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	var mu sync.Mutex
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		b, _ := os.ReadFile("/data/visitas")
		n, _ := strconv.Atoi(strings.TrimSpace(string(b)))
		n++
		if err := os.WriteFile("/data/visitas", []byte(strconv.Itoa(n)), 0o600); err != nil {
			http.Error(w, "no se pudo escribir en /data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "App de ejemplo en kosmos. Visitas: %d\n", n)
	})
	log.Printf("escuchando en :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
