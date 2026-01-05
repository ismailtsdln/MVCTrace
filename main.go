package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mvctrace/internal/detector"
	"mvctrace/internal/httpclient"
	"os"
	"time"
)

const (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	reset  = "\033[0m"
	bold   = "\033[1m"
)

func getConfidenceColor(conf int) string {
	if conf >= 70 {
		return green
	} else if conf >= 40 {
		return yellow
	}
	return red
}

func main() {
	var (
		jsonOutput = flag.Bool("json", false, "Output in JSON format")
		timeout    = flag.Duration("timeout", 10*time.Second, "Request timeout")
		proxy      = flag.String("proxy", "", "HTTP proxy URL")
		silent     = flag.Bool("silent", false, "Minimal output")
	)
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Usage: mvctrace <url>")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	targetURL := flag.Arg(0)

	client := httpclient.NewClient(*timeout, *proxy)
	result := detector.Detect(client, targetURL)

	if *jsonOutput {
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonData))
		return
	}

	if *silent {
		if result.IsMVC {
			fmt.Printf("%sMVC detected%s (%s%s%s confidence)\n", green, reset, getConfidenceColor(result.Confidence), result.ConfidenceLevel(), reset)
		} else {
			fmt.Printf("%sNot MVC%s\n", red, reset)
		}
		return
	}

	// Human-readable output
	fmt.Printf("%sTarget:%s %s\n", blue+bold, reset, result.Target)
	if result.IsMVC {
		fmt.Printf("%sFramework:%s ASP.NET MVC\n", green+bold, reset)
		if result.Version != "" {
			fmt.Printf("%sMVC Version:%s %s (%s%s%s Confidence)\n", green+bold, reset, result.Version, getConfidenceColor(result.Confidence), result.ConfidenceLevel(), reset)
		} else {
			fmt.Printf("%sMVC Detected%s (%s%s%s Confidence)\n", green+bold, reset, getConfidenceColor(result.Confidence), result.ConfidenceLevel(), reset)
		}
	} else {
		fmt.Printf("%sFramework:%s Not ASP.NET MVC\n", red+bold, reset)
		fmt.Printf("%sConfidence:%s %s%s%s\n", red+bold, reset, getConfidenceColor(result.Confidence), result.ConfidenceLevel(), reset)
	}

	if len(result.Evidence) > 0 {
		fmt.Printf("%sEvidence:%s\n", yellow+bold, reset)
		for _, ev := range result.Evidence {
			fmt.Printf("  %s•%s %s\n", yellow, reset, ev.Description)
		}
	}
}
