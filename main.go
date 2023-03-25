package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type ZingSearchServerAccess struct {
	username string
	password string
}

type ZingSearchServer struct {
	host   string
	access ZingSearchServerAccess
}

var zingSearchServer ZingSearchServer = ZingSearchServer{
	host: "localhost:4080",
	access: ZingSearchServerAccess{
		username: "admin",
		password: "Complexpass#123",
	},
}

func main() {

	GatherNameFromStdinAndSearchProductsByName()

	/*responseBodyAsString, err := insertProduct()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("OK, response: ", responseBodyAsString)
	}*/
}

func listIndexes() (string, error) {
	request, err := http.NewRequest("GET", "http://"+zingSearchServer.host+"/api/index", nil)
	request.SetBasicAuth(zingSearchServer.access.username, zingSearchServer.access.password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		errMsg := fmt.Sprintf("Status code is not OK: %d", response.StatusCode)
		return "", errors.New(errMsg)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func insertProduct() (string, error) {

	const index_to_use = "my_index" // ChatGPT says: When inserting data into Zincsearch using its REST API, you need to specify an index for each insertion. An index is a logical namespace that groups related data together, and it is used to organize and search the data in Zincsearch.

	type Product struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}

	// Create a new product to add to Zincsearch
	product := Product{Name: generateRandomProductName(), Price: 19}

	// Convert the product to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		panic(err)
	}

	// Send a POST request to the Zincsearch API to add the new product
	request, err := http.NewRequest("POST", "http://"+zingSearchServer.host+"/api/"+index_to_use+"/_doc", bytes.NewBuffer(productJSON))
	request.SetBasicAuth(zingSearchServer.access.username, zingSearchServer.access.password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Println("Error closing body on defer: ", err)
		}
	}()

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		errMsg := fmt.Sprintf("Status code is not OK: %d", response.StatusCode)
		return "", errors.New(errMsg)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func generateRandomProductName() string {
	var prefixes []string = []string{
		"Cuaderno",
		"Lapicera",
		"Diccionario",
		"Varita",
		"Jornadas",
	}
	var posfixes []string = []string{
		"De Luxe",
		"Super reforzado/a",
		"Lumpurus",
		"De Harray potus",
		"De Disciplina",
	}
	prefixIndex := rand.Intn(len(prefixes))
	posfixIndex := rand.Intn(len(posfixes))
	return prefixes[prefixIndex] + " " + posfixes[posfixIndex]
}

func GatherNameFromStdinAndSearchProductsByName() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1] // removes the enter
	searchProductsByName(name)
}

// TODO : change implementation for returning results (I mean: "promote to" a pure function)
func searchProductsByName(name string) {
	const index_to_use = "my_index" // see comment on @insertProduct#index_to_use

	// Define the query string for the search
	query := `{
        "search_type": "match",
        "query":
        {
            "term": "{{term}}"
        },
        "from": 0,
        "max_results": 20,
        "_source": []
    }` // stolen from https://zincsearch-docs.zinc.dev/api/search/search/#golang-example

	query = strings.Replace(query, "{{term}}", name, 1)
	fmt.Println(query)

	// Send a GET request to the Zincsearch search API
	request, err := http.NewRequest("POST", "http://"+zingSearchServer.host+"/api/"+index_to_use+"/_search", strings.NewReader(query))
	request.SetBasicAuth(zingSearchServer.access.username, zingSearchServer.access.password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Println("Error closing body on defer: ", err)
		}
	}()
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		errMsg := fmt.Sprintf("Status code is not OK: %d", response.StatusCode)
		fmt.Println(errMsg)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	fmt.Println(string(data))

	// Parse the JSON response
	/*var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the search results
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		fmt.Println(source)
	}*/
}
