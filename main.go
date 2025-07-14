package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// encoding functions
func htmlEncode(input string, times int) string {
	out := input
	for i := 0; i < times; i++ {
		out = html.EscapeString(out)
	}
	return out
}

func urlEncode(input string, times int) string {
	out := input
	for i := 0; i < times; i++ {
		out = url.QueryEscape(out)
	}
	return out
}

func jsEscape(input string, times int) string {
	out := input
	for i := 0; i < times; i++ {
		out = strings.ReplaceAll(out, "\\", "\\\\")
		out = strings.ReplaceAll(out, "'", "\\'")
		out = strings.ReplaceAll(out, `"`, `\\\"`)
	}
	return out
}

func unicodeEscape(input string, times int) string {
	out := input
	for t := 0; t < times; t++ {
		var buf strings.Builder
		for _, r := range out {
			buf.WriteString(fmt.Sprintf("\\u%04x", r))
		}
		out = buf.String()
	}
	return out
}

func base64Encode(input string, times int) string {
	out := input
	for i := 0; i < times; i++ {
		out = base64.StdEncoding.EncodeToString([]byte(out))
	}
	return out
}

func prompt(promptText string) string {
	fmt.Print(promptText + ": ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return s.Text()
}

func parseRawRequest(raw string) (*http.Request, error) {
	r := bufio.NewReader(strings.NewReader(raw))
	req, err := http.ReadRequest(r)
	if err != nil {
		return nil, errors.New("failed to parse HTTP request: " + err.Error())
	}
	// rewind body
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.ContentLength = int64(len(bodyBytes))
	}
	return req, nil
}

func main() {
	// ASCII art footer
	fmt.Println(`   ___      __    _               __   _  __  
  / _ )    / /   (_)   ___    ___/ /  | |/_/  
 / _  |   / /   / /   / _ \\  / _  /  _>  <   
/____/   /_/   /_/   /_//_/  \\,_/  /_/|_|      
                                            
      BlindX by progprnv
`)

	// 1) Paste raw POST call
	raw := prompt("Paste ur POST req call from burp (end with an empty line)")
	for {
		line := prompt("")
		if line == "" {
			break
		}
		raw += "\r\n" + line
	}

	req, err := parseRawRequest(raw)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}

	// Collect parameters
	params := []string{}
	for {
		p := prompt("Enter parameter to inject payload into")
		params = append(params, p)
		x := prompt("Do you have any other parameter to inject payload? (Y/N)")
		if strings.ToLower(x) != "y" {
			break
		}
	}

	// Payload
	payload := prompt("Enter your payload")

	// Encoding options
	fmt.Println("Type of encoding for payload:")
	types := []string{
		"1) HTML encode x1", "2) HTML encode x2", "3) HTML encode x3",
		"4) URL encode x1", "5) URL encode x2", "6) URL encode x3",
		"7) JS escape x1", "8) JS escape x2", "9) JS escape x3",
		"10) Unicode escape x1", "11) Unicode escape x2", "12) Unicode escape x3",
		"13) Base64 encode x1", "14) Base64 encode x2", "15) Base64 encode x3",
		"16) All variants (1-15)", "17) No encode (original)"}
	for _, t := range types {
		fmt.Println(t)
	}
	choice := prompt("Select option number (1-17)")

	// Prepare payload variants
	oList := []string{}

	sel, _ := strconv.Atoi(choice)
	switch {
	case sel >= 1 && sel <= 3:
		oList = append(oList, htmlEncode(payload, sel))
	case sel >= 4 && sel <= 6:
		oList = append(oList, urlEncode(payload, sel-3))
	case sel >= 7 && sel <= 9:
		oList = append(oList, jsEscape(payload, sel-6))
	case sel >= 10 && sel <= 12:
		oList = append(oList, unicodeEscape(payload, sel-9))
	case sel >= 13 && sel <= 15:
		oList = append(oList, base64Encode(payload, sel-12))
	case sel == 16:
		for i := 1; i <= 15; i++ {
			switch {
			case i <= 3:
				oList = append(oList, htmlEncode(payload, i))
			case i <= 6:
				oList = append(oList, urlEncode(payload, i-3))
			case i <= 9:
				oList = append(oList, jsEscape(payload, i-6))
			case i <= 12:
				oList = append(oList, unicodeEscape(payload, i-9))
			case i <= 15:
				oList = append(oList, base64Encode(payload, i-12))
			}
		}
	case sel == 17:
		oList = append(oList, payload)
	default:
		fmt.Println("Invalid choice, exiting.")
		return
	}

	// Header injection
	headers := map[string]string{}
	x := prompt("Need any additional or existing header to inject payload? (Y/N)")
	if strings.ToLower(x) == "y" {
		for {
			headerName := prompt("Enter header name")
			headerVal := prompt("Enter header value (use {{payload}})")
			headers[headerName] = headerVal
			more := prompt("More headers? (Y/N)")
			if strings.ToLower(more) != "y" {
				break
			}
		}
	}

	client := &http.Client{}

	// Send each variant
	for idx, ep := range oList {
		r2 := req.Clone(req.Context())
		// update body
		vals := url.Values{}
		if r2.Body != nil {
			b, _ := io.ReadAll(r2.Body)
			vals, _ = url.ParseQuery(string(b))
		}
		for _, p := range params {
			vals.Set(p, ep)
		}
		bodyStr := vals.Encode()
		r2.Body = io.NopCloser(strings.NewReader(bodyStr))
		r2.ContentLength = int64(len(bodyStr))

		// update headers
		for hn, hv := range headers {
			r2.Header.Set(hn, strings.ReplaceAll(hv, "{{payload}}", ep))
		}

		resp, err := client.Do(r2)
		if err != nil {
			fmt.Printf("[%d] Error: %v\n", idx+1, err)
			continue
		}
		fmt.Printf("[%d] %s -> %s\n", idx+1, r2.URL.String(), resp.Status)
		resp.Body.Close()
	}
}
