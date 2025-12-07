package proxy

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type Handler struct {
	ServiceURL     string
	AllowedHeaders []string
	Client         *http.Client // переиспользуемый клиент без пересоздания (в отличие от Laravel Octane)
}

func NewController(serviceURL string, allowedHeaders []string) *Handler {
	return &Handler{
		ServiceURL:     serviceURL,
		AllowedHeaders: allowedHeaders,
		Client:         &http.Client{}, // можно настроить таймауты и т.д.
	}
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	// Получаем new_route из query или заголовков
	newRoute := r.URL.Query().Get("new_route")
	if newRoute == "" {
		return fmt.Errorf("missing new_route")
	}

	targetURL := h.ServiceURL + "/" + strings.TrimPrefix(newRoute, "/")

	// Копируем query-параметры, исключая new_route
	query := r.URL.Query()
	query.Del("new_route")
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return err
	}
	parsedURL.RawQuery = query.Encode()

	// Фильтруем заголовки
	headers := h.filterRequestHeaders(r.Header)

	var body io.Reader
	if err := r.ParseMultipartForm(32 << 20); err == nil && r.MultipartForm != nil {
		// Есть файлы — формируем multipart
		body, headers, err = h.buildMultipartBody(r, headers)
		if err != nil {
			return err
		}
	} else {
		body = r.Body
	}

	proxyReq, err := http.NewRequestWithContext(ctx, r.Method, parsedURL.String(), body)
	if err != nil {
		return err
	}

	for k, v := range headers {
		proxyReq.Header[k] = v
	}

	resp, err := h.Client.Do(proxyReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filteredHeaders := h.filterResponseHeaders(resp.Header)

	for k, v := range filteredHeaders {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	return err
}

func (c *Handler) filterRequestHeaders(h http.Header) http.Header {
	filtered := make(http.Header)
	for key, values := range h {
		lowerKey := strings.ToLower(key)
		if lowerKey == "host" || lowerKey == "new_route" {
			continue
		}
		filtered[key] = values
	}
	return filtered
}

func (c *Handler) filterResponseHeaders(h http.Header) http.Header {
	filtered := make(http.Header)
	for key, values := range h {
		lowerKey := strings.ToLower(key)
		for _, allowed := range c.AllowedHeaders {
			if lowerKey == strings.ToLower(allowed) {
				filtered[key] = values
				break
			}
		}
	}
	return filtered
}

func (c *Handler) buildMultipartBody(r *http.Request, headers http.Header) (io.Reader, http.Header, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Копируем поля формы
	for key, values := range r.MultipartForm.Value {
		for _, value := range values {
			if err := writer.WriteField(key, value); err != nil {
				return nil, nil, err
			}
		}
	}

	// Копируем файлы
	for fieldName, fileHeaders := range r.MultipartForm.File {
		for _, fh := range fileHeaders {
			file, err := fh.Open()
			if err != nil {
				return nil, nil, err
			}
			defer file.Close()

			part, err := writer.CreateFormFile(fieldName, fh.Filename)
			if err != nil {
				return nil, nil, err
			}
			if _, err := io.Copy(part, file); err != nil {
				return nil, nil, err
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, nil, err
	}

	// Обновляем Content-Type с boundary
	headers.Set("Content-Type", writer.FormDataContentType())
	return body, headers, nil
}
