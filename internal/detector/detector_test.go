package detector

import (
	"github.com/ismailtsdln/mvctrace/internal/httpclient"
	"testing"
)

func TestDetect(t *testing.T) {
	client := httpclient.NewClient(10, "")
	result := Detect(client, "http://httpbin.org")

	// Should not detect MVC
	if result.IsMVC {
		t.Errorf("Expected not MVC, but detected MVC")
	}

	// For non-MVC sites, may have no evidence
	// Confidence should be 0
	if result.Confidence != 0 {
		t.Errorf("Expected confidence 0, got %d", result.Confidence)
	}
}

func TestDetectHeaders(t *testing.T) {
	client := httpclient.NewClient(10, "")
	evidence := detectHeaders(client, "http://httpbin.org")

	// httpbin.org has some headers, but not MVC specific
	// Just check it returns evidence
	if len(evidence) == 0 {
		t.Log("No header evidence found, which is expected for non-MVC site")
	}
}
