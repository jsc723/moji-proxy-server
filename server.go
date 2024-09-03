package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/gocolly/colly"
)

var (
	ApplicationID string
	MojiVersion   string
	SessionToken  string
	Version       = "1.4.7"
)

func main() {
	username := flag.String("username", "", "Specify the Moji username")
	password := flag.String("password", "", "Specify the Moji password")
	showVersion := flag.Bool("version", false, "Speficify if you want to check version")
	port := flag.Int("port", 9285, "Specify port number to listen")
	// Parse command-line arguments
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s\n", Version)
		return
	}

	fmt.Printf("moji-proxy-server version: %s\n", Version)

	getApplicationIDAndVersion()
	getSessionToken(*username, *password)

	// Define the HTTP server route for '/search'
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/details", detailsHandler)
	http.HandleFunc("/healthcheck", healthcheckHandler)

	// Start the HTTP server
	log.Printf("Starting server on :%d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func doPost(url string, payload map[string]interface{}) (map[string]interface{}, error) {
	payload["_ApplicationId"] = ApplicationID
	payload["_ClientVersion"] = "js3.4.1"
	payload["g_os"] = "PCWeb"
	payload["g_ver"] = MojiVersion
	if SessionToken != "" {
		payload["_SessionToken"] = SessionToken
	}
	// Convert the request body to JSON
	apiRequestBodyJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding JSON\n")
		return nil, err
	}

	//log.Printf("request json: %s\n", apiRequestBodyJSON)

	buf := bytes.NewBuffer(apiRequestBodyJSON)
	req, err := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprint(len(apiRequestBodyJSON)))
	client := &http.Client{}
	log.Printf("request: %+v\n", req)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making API request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("response status: %+v\n", resp.Status)
	if resp.StatusCode != 200 {
		log.Printf("Status code is not ok")
		return nil, fmt.Errorf("Status code is not ok")
	}
	// Read the response body
	apiResponseBodyJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response")
		return nil, err
	}

	//log.Printf("%+s\n", apiResponseBodyJSON);

	var apiResponseObject map[string]any

	// Unmarshal the JSON string into the map
	err = json.Unmarshal([]byte(apiResponseBodyJSON), &apiResponseObject)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return apiResponseObject, nil
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("version=" + Version))
}

type SearchRequest struct {
	Query  string `json:"query"`
	Expand bool   `json:"expand"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("processing the request %v\n", r.URL)

	// Decode the JSON request body
	var searchReq SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Construct the request body for the union-api
	apiRequestBody := map[string]interface{}{
		"functions": []map[string]interface{}{
			{
				"name": "search-all",
				"params": map[string]interface{}{
					"text":  searchReq.Query,
					"types": []int{102, 106},
				},
			},
		},
	}

	apiResponseObject, err := doPost("https://api.mojidict.com/parse/functions/union-api", apiRequestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	serverResponse, err := getNestedField(apiResponseObject, "result", "results", "search-all")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating server response: %v", err), http.StatusInternalServerError)
		return
	}
	serverResponseJSON, err := json.Marshal(serverResponse)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	//log.Printf("get response from moji: %s\n", apiResponseBody)

	// Send the API response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(serverResponseJSON)
}

type DetailsRequest struct {
	ObjectIds []string `json:"objectIds"`
}

type ObjectWrapper struct {
	ObjectID string `json:"objectId"`
}

func detailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("processing the request %v\n", r.URL)

	var detailsRequest DetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&detailsRequest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Wrap each objectId into an object
	var itemsJson []ObjectWrapper
	for _, objectId := range detailsRequest.ObjectIds {
		wrappedObject := ObjectWrapper{ObjectID: objectId}
		itemsJson = append(itemsJson, wrappedObject)
	}

	// // Construct the request body for the union-api
	apiRequestBody := map[string]interface{}{
		"itemsJson":       itemsJson,
		"skipAccessories": false,
	}

	apiResponseBody, err := doPost("https://api.mojidict.com/parse/functions/nlt-fetchManyLatestWords", apiRequestBody)

	apiResponseBody, err = transformResultObject(apiResponseBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	serverResponseJSON, err := json.Marshal(apiResponseBody)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	// Send the API response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(serverResponseJSON)

}

func transformResultObject(in map[string]any) (map[string]any, error) {
	nested, err := getNestedField(in, "result", "result")
	if err != nil {
		return nil, fmt.Errorf("Cannot find 'result.result' in the response object")
	}
	results := nested.([]any)
	var words []Word
	for _, result := range results {
		w := processSingleWord(result.(map[string]any))
		if w.Id != "" {
			words = append(words, w)
		}
	}

	res := map[string]any{
		"words": words,
	}
	return res, nil
}

type Word struct {
	Id         string          `json:"id"`
	Spell      string          `json:"spell"`
	Pron       string          `json:"pron"`
	Accent     string          `json:"accent"`
	Excerpt    string          `json:"excerpt"`
	Details    []WordDetail    `json:"details"`
	SubDetails []WordSubDetail `json:"subDetails"`
}
type WordDetail struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
type WordSubDetail struct {
	Id       string        `json:"id"`
	Title    string        `json:"title"`
	DetailId string        `json:"detailId"`
	Examples []WordExample `json:"examples"`
}
type WordExample struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Trans string `json:"trans"`
}

func processSingleWord(word map[string]interface{}) Word {
	var w Word
	w.Spell = getNestedFieldWithDefault("", word, "word", "spell").(string)
	w.Pron = getNestedFieldWithDefault("", word, "word", "pron").(string)
	w.Accent = getNestedFieldWithDefault("", word, "word", "accent").(string)
	w.Excerpt = getNestedFieldWithDefault("", word, "word", "excerpt").(string)
	emptyArray := []any{}
	for _, detail := range getNestedFieldWithDefault(emptyArray, word, "details").([]any) {
		var d WordDetail
		d.Id = getNestedFieldWithDefault("", detail.(map[string]any), "objectId").(string)
		d.Title = getNestedFieldWithDefault("", detail.(map[string]any), "title").(string)
		w.Details = append(w.Details, d)
	}

	idToSubdetailIdx := make(map[string]int)
	for _, subdetail := range getNestedFieldWithDefault(emptyArray, word, "subdetails").([]any) {
		var d WordSubDetail
		d.Id = getNestedFieldWithDefault("", subdetail.(map[string]any), "objectId").(string)
		d.Title = getNestedFieldWithDefault("", subdetail.(map[string]any), "title").(string)
		d.DetailId = getNestedFieldWithDefault("", subdetail.(map[string]any), "detailsId").(string)
		d.Examples = make([]WordExample, 0)
		w.SubDetails = append(w.SubDetails, d)
		idToSubdetailIdx[d.Id] = len(w.SubDetails) - 1
	}

	for _, example := range getNestedFieldWithDefault(emptyArray, word, "examples").([]any) {
		var e WordExample
		e.Id = getNestedFieldWithDefault("", example.(map[string]any), "objectId").(string)
		e.Title = getNestedFieldWithDefault("", example.(map[string]any), "title").(string)
		e.Trans = getNestedFieldWithDefault("", example.(map[string]any), "trans").(string)
		subid := getNestedFieldWithDefault("", example.(map[string]any), "subdetailsId").(string)
		if idx, ok := idToSubdetailIdx[subid]; ok {
			w.SubDetails[idx].Examples = append(w.SubDetails[idx].Examples, e)
		} else {
			log.Printf("example %v no match subdetail\n", subid)
		}
	}

	w.Id = getNestedFieldWithDefault("", word, "word", "objectId").(string)
	return w
}

func getApplicationIDAndVersion() {
	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		c      = colly.NewCollector()
		re     = regexp.MustCompile(`_ApplicationId\s*=\s*"([A-Za-z0-9]+)"`)
		ver_re = regexp.MustCompile(`(v\d+\.\d+\.\d+\.\d+)/[A-Za-z0-9]+\.js`)
		url    = "https://www.mojidict.com/"
	)
	log.Println("Getting application ID and version...")

	// Setup the callbacks for HTML and Response
	c.OnHTML("link[rel=preload][as=script]", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		if MojiVersion == "" {
			fmt.Printf("check href = %v\n", href)
			matches := ver_re.FindStringSubmatch(href)
			if len(matches) > 1 {
				MojiVersion = matches[1]
				fmt.Printf("MojiVersion = %v\n", MojiVersion)
			}
		}
		

		// Increment the WaitGroup counter
		wg.Add(1)

		// Start a new goroutine to visit the href concurrently
		go func(href string) {
			defer wg.Done()

			// Visit the href
			//fmt.Printf("visit %v\n", href)
			err := c.Visit(href)
			if err != nil {
				fmt.Printf("Error visiting %s: %v\n", href, err)
				return
			}
		}(href)
	})

	// Callback to handle the response from the visited page
	c.OnResponse(func(r *colly.Response) {
		mu.Lock()
		defer mu.Unlock()

		// Find matches using the regular expression pattern
		matches := re.FindAllStringSubmatch(string(r.Body), -1)

		// Update the result map
		for _, match := range matches {
			ApplicationID = match[1]
		}
	})

	// Start the crawling process to get the Application Id
	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	if ApplicationID == "" {
		log.Fatal("Application Id not found. Please ensure the crawling process is successful.")
	}

	if MojiVersion == "" {
		log.Fatal("Moji version not found. Please ensure the crawling process is successful")
	}

	log.Printf("applicationID = %v\n", ApplicationID)
}

func getSessionToken(cmdUsername string, cmdPassword string) {
	var username, password string
	if cmdUsername != "" && cmdPassword != "" {
		log.Println("Using command line login info")
		username = cmdUsername
		password = cmdPassword
	} else {
		config := loadConfig("config.json")
		if config.Username != "" && config.Password != "" {
			log.Println("Using config file login info")
			username = config.Username
			password = config.Password
		}
	}
	if username != "" && password != "" {
		apiRequestBody := map[string]interface{}{
			"username": username,
			"password": password,
		}
		apiResponseBody, err := doPost("https://api.mojidict.com/parse/login", apiRequestBody)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		SessionToken = getNestedFieldWithDefault("", apiResponseBody, "sessionToken").(string)
		if SessionToken != "" {
			log.Printf("Get session token = %v\n", SessionToken)
		} else {
			log.Printf("Failed to get the session token (you can still use the app)\n")
		}
	} else {
		log.Printf("No login info provided, skip getting session token (you can still use the app)\n")
	}
}
