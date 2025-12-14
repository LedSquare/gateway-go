package handlers

import (
	"gateway-go/internal/response"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/k0kubun/pp"
)

type Handler struct {
	ServiceURL     string
	AllowedHeaders []string
	Client         *http.Client
	Op             string
	version        string
}

func NewHandler(serviceURL string, allowedHeaders []string, op string, version string) *Handler {
	return &Handler{
		ServiceURL:     serviceURL,
		AllowedHeaders: allowedHeaders,
		Client:         &http.Client{}, // можно настроить таймауты и т.д.
		Op:             op,
		version:        version,
	}
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	subPath := chi.URLParam(r, "*")

	targetURL := h.ServiceURL +
		"/" + h.version +
		"/" + strings.Trim(subPath, "/")

	pp.Println(targetURL)
	// pp.Print(targetURL)

	// Создаём прокси-запрос с тем же методом и телом
	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, r.Body)
	if err != nil {
		render.JSON(w, r, response.Error("Invalid request", http.StatusBadRequest, h.Op))
		return
	}

	for key, values := range r.Header {
		if strings.ToLower(key) == "host" {
			continue
		}
		proxyReq.Header[key] = values
	}

	// Выполняем запрос
	resp, err := h.Client.Do(proxyReq)
	if err != nil {
		render.JSON(w, r, response.Error("Сервис недоступен", http.StatusServiceUnavailable, h.Op))
		return
	}
	defer resp.Body.Close()

	// Копируем разрешённые заголовки ответа
	for key, values := range resp.Header {
		if h.isHeaderAllowed(key) {
			w.Header()[key] = values
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (h *Handler) isHeaderAllowed(key string) bool {
	keyLower := strings.ToLower(key)
	for _, allowed := range h.AllowedHeaders {
		if keyLower == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}
