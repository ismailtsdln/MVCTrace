package result

// Evidence represents a piece of evidence found during detection
type Evidence struct {
	Description string `json:"description"`
	Source      string `json:"source"`          // Where was it found (header, html, route, etc)
	Value       string `json:"value,omitempty"` // The actual value found
	Confidence  int    `json:"confidence"`      // 0-100
}

// Result holds the detection outcome
type Result struct {
	Target        string     `json:"target"`
	IsMVC         bool       `json:"is_mvc"`
	Version       string     `json:"version,omitempty"`
	VersionSource string     `json:"version_source,omitempty"` // Where version was found
	Confidence    int        `json:"confidence"`
	Evidence      []Evidence `json:"evidence"`
}

// ConfidenceLevel returns a string representation of confidence
func (r *Result) ConfidenceLevel() string {
	if r.Confidence >= 70 {
		return "High"
	} else if r.Confidence >= 40 {
		return "Medium"
	}
	return "Low"
}
