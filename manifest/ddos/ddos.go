package ddos

type Manifest struct {
	Spec   Spec `json:"spec,omitempty"`
	Status Spec `json:"status,omitempty"`
}

type Spec struct {
	Url   string `json:"url"`
	Start bool   `json:"start"`
}
