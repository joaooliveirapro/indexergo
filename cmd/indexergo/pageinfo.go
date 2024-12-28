package indexergo

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type PageInfo struct {
	HTTPResponse  HTTPResponse   `json:"httpResponse"`
	HTMLTags      map[string]int `json:"htmlTags"`
	ContentTokens map[string]int `json:"contentTokens"`
	Timestamp     string         `json:"timestamp"`
}

func NewPageInfo(url string) (*PageInfo, error) {
	// HTTP request to URL
	response, err := GetHTML(url)
	if err != nil {
		return nil, fmt.Errorf("[error] Couldn't get response for %s", response.URL)
	}
	// Build PageInfo based on response data
	pageInfo := PageInfo{
		Timestamp:     time.Now().Format("02-01-2006 15:04:05"),
		HTTPResponse:  response,
		HTMLTags:      map[string]int{},
		ContentTokens: map[string]int{},
	}
	return &pageInfo, nil
}

type HTTPResponse struct {
	StatusCode      int      `json:"statusCode"`
	URL             string   `json:"url"`
	Redirected      bool     `json:"redirected"`
	HTML            string   `json:"-"` // don't include
	RedirectHistory []string `json:"redirectsHistory"`
}

func GetHTML(url string) (HTTPResponse, error) {
	var HttpResponse HTTPResponse
	HttpResponse.URL = url

	// Create a request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		e := fmt.Errorf("[debug] error creating request: %v", err)
		return HttpResponse, e
	}
	// Create a new HTTP Client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			HttpResponse.Redirected = true
			// Append previous URLs to the redirect history
			for _, r := range via {
				HttpResponse.RedirectHistory = append(HttpResponse.RedirectHistory, r.URL.String())
			}
			HttpResponse.RedirectHistory = append(HttpResponse.RedirectHistory, req.URL.String())
			return nil
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		e := fmt.Errorf("[debug] error making request: %v", err)
		return HttpResponse, e
	}
	defer resp.Body.Close()
	// Save status code
	HttpResponse.StatusCode = resp.StatusCode

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Errorf("[debug] error reading response body: %v", err)
		return HttpResponse, e
	}
	// Response is OK
	if resp.StatusCode != http.StatusOK {
		e := fmt.Errorf("[debug] HTTP code error %d", resp.StatusCode)
		return HttpResponse, e
	}
	HttpResponse.HTML = string(body)
	return HttpResponse, nil
}
