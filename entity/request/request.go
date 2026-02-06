package request

type Request struct {
	Category   string `json:"category"`
	Action     string `json:"action"`
	Parameters map[string]string
}
