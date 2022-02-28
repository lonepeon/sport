package web

type Response struct {
	HTTPCode   int
	Layout     string
	LogMessage string
	Data       interface{}
	Template   string
}

func (r Response) Templates() []string {
	var templates []string
	if r.Layout != "" {
		templates = append(templates, r.Layout)
	}
	templates = append(templates, r.Template)

	return templates
}
