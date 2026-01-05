package detector

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ismailtsdln/mvctrace/internal/httpclient"
	"github.com/ismailtsdln/mvctrace/internal/result"
)

// Detect performs all detection methods and aggregates results
func Detect(client *httpclient.Client, targetURL string) *result.Result {
	res := &result.Result{
		Target:        targetURL,
		IsMVC:         false,
		Version:       "",
		VersionSource: "",
		Evidence:      []result.Evidence{},
	}

	// Perform detections
	evidence := []result.Evidence{}
	evidence = append(evidence, detectHeaders(client, targetURL)...)
	evidence = append(evidence, detectHTML(client, targetURL)...)
	evidence = append(evidence, detectRoutes(client, targetURL)...)
	evidence = append(evidence, detectErrorPage(client, targetURL)...)
	evidence = append(evidence, detectStaticFiles(client, targetURL)...)

	res.Evidence = evidence

	// Aggregate confidence and track version source
	totalConf := 0
	versionCandidates := map[string]string{} // version -> source
	for _, ev := range evidence {
		totalConf += ev.Confidence
		if strings.Contains(ev.Description, "MVC Version") && ev.Value != "" {
			versionCandidates[ev.Value] = ev.Source
		}
	}

	if totalConf > 100 {
		totalConf = 100
	}
	res.Confidence = totalConf

	// Determine if MVC
	if totalConf >= 40 {
		res.IsMVC = true
	}

	// Determine version - prefer the first one found with highest confidence
	for _, ev := range evidence {
		if strings.Contains(ev.Description, "MVC Version") && ev.Value != "" {
			res.Version = ev.Value
			res.VersionSource = ev.Source
			break
		}
	}

	return res
}

// detectHeaders checks for MVC-specific headers
func detectHeaders(client *httpclient.Client, url string) []result.Evidence {
	resp, err := client.Get(url)
	if err != nil {
		return []result.Evidence{}
	}
	defer resp.Body.Close()

	evidence := []result.Evidence{}

	if ver := resp.Header.Get("X-AspNetMvc-Version"); ver != "" {
		evidence = append(evidence, result.Evidence{
			Description: fmt.Sprintf("MVC Version %s detected", ver),
			Source:      "HTTP Header: X-AspNetMvc-Version",
			Value:       ver,
			Confidence:  90,
		})
	}

	if ver := resp.Header.Get("X-AspNet-Version"); ver != "" {
		evidence = append(evidence, result.Evidence{
			Description: fmt.Sprintf(".NET Framework version detected: %s", ver),
			Source:      "HTTP Header: X-AspNet-Version",
			Value:       ver,
			Confidence:  20,
		})
	}

	if powered := resp.Header.Get("X-Powered-By"); strings.Contains(powered, "ASP.NET") {
		evidence = append(evidence, result.Evidence{
			Description: "ASP.NET technology stack identified",
			Source:      "HTTP Header: X-Powered-By",
			Value:       powered,
			Confidence:  10,
		})
	}

	return evidence
}

// detectHTML checks for MVC-specific HTML markers
func detectHTML(client *httpclient.Client, url string) []result.Evidence {
	body, resp, err := client.GetBody(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return []result.Evidence{}
	}

	evidence := []result.Evidence{}

	if strings.Contains(body, `data-val="true"`) {
		evidence = append(evidence, result.Evidence{
			Description: "MVC validation attributes detected",
			Source:      "HTML Body: data-val attribute",
			Value:       `data-val="true"`,
			Confidence:  60,
		})
	}

	if strings.Contains(body, "__MVCFormValidation") {
		evidence = append(evidence, result.Evidence{
			Description: "MVC form validation script found",
			Source:      "HTML Body: __MVCFormValidation script",
			Value:       "__MVCFormValidation",
			Confidence:  70,
		})
	}

	if strings.Contains(body, "jquery.validate.unobtrusive") {
		evidence = append(evidence, result.Evidence{
			Description: "jQuery unobtrusive validation library referenced",
			Source:      "HTML Body: Script reference",
			Value:       "jquery.validate.unobtrusive",
			Confidence:  50,
		})
	}

	if strings.Contains(body, "System.Web.Mvc") {
		evidence = append(evidence, result.Evidence{
			Description: "System.Web.Mvc namespace referenced",
			Source:      "HTML Body: Script or attribute",
			Value:       "System.Web.Mvc",
			Confidence:  80,
		})
	}

	return evidence
}

// detectRoutes tests common MVC routes
func detectRoutes(client *httpclient.Client, baseURL string) []result.Evidence {
	routes := []string{"/Home/Index", "/Account/Login", "/Home/About"}
	evidence := []result.Evidence{}

	for _, route := range routes {
		url := strings.TrimSuffix(baseURL, "/") + route
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			evidence = append(evidence, result.Evidence{
				Description: "MVC default route is accessible",
				Source:      fmt.Sprintf("HTTP Route: %s", route),
				Value:       fmt.Sprintf("HTTP %d", resp.StatusCode),
				Confidence:  40,
			})
		}
	}

	return evidence
}

// detectErrorPage checks error pages for MVC indicators
func detectErrorPage(client *httpclient.Client, baseURL string) []result.Evidence {
	url := strings.TrimSuffix(baseURL, "/") + "/nonexistent-path-12345"
	body, resp, err := client.GetBody(url)
	if err != nil {
		return []result.Evidence{}
	}

	evidence := []result.Evidence{}

	if resp.StatusCode == http.StatusNotFound {
		if strings.Contains(body, "The resource cannot be found") || strings.Contains(body, "Server Error in '/' Application") {
			evidence = append(evidence, result.Evidence{
				Description: "MVC error page detected",
				Source:      "HTTP Error Page: /nonexistent-path-12345",
				Value:       "404 Error Message Pattern",
				Confidence:  30,
			})
		}
	}

	// Check for stack traces indicating MVC
	if strings.Contains(body, "System.Web.Mvc") {
		evidence = append(evidence, result.Evidence{
			Description: "MVC stack trace in error page",
			Source:      "HTTP Error Page: /nonexistent-path-12345",
			Value:       "System.Web.Mvc Stack Trace",
			Confidence:  50,
		})
	}

	return evidence
}

// detectStaticFiles checks for MVC-specific static files
func detectStaticFiles(client *httpclient.Client, baseURL string) []result.Evidence {
	files := []string{"/Content/Site.css", "/Scripts/jquery-1.10.2.js", "/bundles/jquery"}
	evidence := []result.Evidence{}

	for _, file := range files {
		url := strings.TrimSuffix(baseURL, "/") + file
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			evidence = append(evidence, result.Evidence{
				Description: "MVC static file structure detected",
				Source:      fmt.Sprintf("Static File: %s", file),
				Value:       fmt.Sprintf("HTTP %d", resp.StatusCode),
				Confidence:  20,
			})
		}
	}

	return evidence
}
