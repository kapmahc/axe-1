package axe

// H hash data
type H map[string]interface{}

// K context key type
type K string

// HTTPError http error
type HTTPError struct {
	Message string
	Status  int
}

func (p *HTTPError) Error() string {
	return p.Message
}
