package controladores

import (
	"io"
	"net/http"
	"os"
)

// ---- Proxy GET /api/logs/query?id=X&limite=N → back ----
func (c *Controlador) APILogsQuery(w http.ResponseWriter, r *http.Request) {
	backURL := os.Getenv("BACK_URL")
	if backURL == "" {
		backURL = "http://back:10001"
	}

	id := r.URL.Query().Get("id")
	limite := r.URL.Query().Get("limite")
	if limite == "" {
		limite = "10"
	}

	resp, err := http.Get(backURL + "/api/logs/query?id=" + id + "&limite=" + limite)
	if err != nil {
		http.Error(w, "Backend no disponible", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// ---- Proxy GET /api/logs/recolectar?id=X → back (síncrono, espera a que termine) ----
func (c *Controlador) APILogsRecolectar(w http.ResponseWriter, r *http.Request) {
	backURL := os.Getenv("BACK_URL")
	if backURL == "" {
		backURL = "http://back:10001"
	}

	id := r.URL.Query().Get("id")
	url := backURL + "/api/logs/recolectar"
	if id != "" {
		url += "?id=" + id
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Backend no disponible", http.StatusServiceUnavailable)
		return
	}
	resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
}

// ---- Proxy GET /api/logs/stream/{id} → back (SSE) ----
func (c *Controlador) APILogsStream(w http.ResponseWriter, r *http.Request) {
	backURL := os.Getenv("BACK_URL")
	if backURL == "" {
		backURL = "http://back:10001"
	}

	id := r.PathValue("id")

	req, err := http.NewRequestWithContext(r.Context(), "GET", backURL+"/api/logs/stream/"+id, nil)
	if err != nil {
		http.Error(w, "Error creando request", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Backend no disponible", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			if ok {
				flusher.Flush()
			}
		}
		if err != nil {
			break
		}
	}
}
