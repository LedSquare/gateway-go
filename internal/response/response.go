package response

type JsonResponse struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Operation string `json:"operation,omitempty"`
}

func Error(msg string, status int, op string) JsonResponse {
	return JsonResponse{
		Message:   msg,
		Status:    status,
		Operation: op,
	}
}
