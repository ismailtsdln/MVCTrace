package detector

import (
	"mvctrace/internal/httpclient"
	"mvctrace/internal/result"
	"net/http"
	"strings"
)

// Detect performs all detection methods and aggregates results
func Detect(client *httpclient.Client, targetURL string) *result.Result {
	res := &result.Result{
		Target:   targetURL,
		IsMVC:    false,
		Version:  "",
		Evidence: []result.Evidence{},
	}

	// Perform detections
	evidence := []result.Evidence{}
	evidence = append(evidence, detectHeaders(client, targetURL)...)
	evidence = append(evidence, detectHTML(client, targetURL)...)
	evidence = append(evidence, detectRoutes(client, targetURL)...)
	evidence = append(evidence, detectErrorPage(client, targetURL)...)
	evidence = append(evidence, detectStaticFiles(client, targetURL)...)

	res.Evidence = evidence

	// Aggregate confidence
	totalConf := 0
	versionCandidates := map[string]int{}
	for _, ev := range evidence {
		totalConf += ev.Confidence
		if strings.Contains(ev.Description, "MVC Version") {
			// Extract version from description
			if idx := strings.Index(ev.Description, "MVC Version "); idx != -1 {
				ver := strings.TrimSpace(ev.Description[idx+12:])
				if dotIdx := strings.Index(ver, " "); dotIdx != -1 {
					ver = ver[:dotIdx]
				}
				versionCandidates[ver]++
			}
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

	// Determine version
	maxCount := 0
	for ver, count := range versionCandidates {
		if count > maxCount {
			maxCount = count
			res.Version = ver
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
			Description: "X-AspNetMvc-Version header detected: " + ver,
			Confidence:  90,
		})
		evidence = append(evidence, result.Evidence{
			Description: "MVC Version " + ver + " inferred from header",
			Confidence:  0, // Already counted
		})
	}

	if ver := resp.Header.Get("X-AspNet-Version"); ver != "" {
		evidence = append(evidence, result.Evidence{
			Description: "X-AspNet-Version header detected: " + ver,
			Confidence:  20, // Indicates .NET, but not specifically MVC
		})
	}

	if powered := resp.Header.Get("X-Powered-By"); strings.Contains(powered, "ASP.NET") {
		evidence = append(evidence, result.Evidence{
			Description: "X-Powered-By indicates ASP.NET",
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
			Description: "MVC validation attributes (data-val) found in HTML",
			Confidence:  60,
		})
	}

	if strings.Contains(body, "__MVCFormValidation") {
		evidence = append(evidence, result.Evidence{
			Description: "__MVCFormValidation script reference found",
			Confidence:  70,
		})
	}

	if strings.Contains(body, "jquery.validate.unobtrusive") {
		evidence = append(evidence, result.Evidence{
			Description: "jquery.validate.unobtrusive script found",
			Confidence:  50,
		})
	}

	if strings.Contains(body, "System.Web.Mvc") {
		evidence = append(evidence, result.Evidence{
			Description: "System.Web.Mvc reference in HTML",
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
				Description: "Default MVC route " + route + " responded with 200",
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
				Description: "MVC-style 404 error page detected",
				Confidence:  30,
			})
		}
	}

	// Check for stack traces indicating MVC
	if strings.Contains(body, "System.Web.Mvc") {
		evidence = append(evidence, result.Evidence{
			Description: "MVC stack trace in error page",
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
				Description: "MVC static file " + file + " found",
				Confidence:  20,
			})
		}
	}

	return evidence
}
