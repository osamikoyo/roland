package request

type Request struct{
	Category string `json:"category"`
	Parameters struct{
		Query string `json:"query"`
		Action string `json:"action"`
	}
}