package zincsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ProductsServer struct {
	Server
	Index string
}

type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var DefaultProductsServer ProductsServer = ProductsServer{
	Server: DefaultServer,
	Index:  "products_index",
}

func (p ProductsServer) InsertRandomProduct() ([]byte, *Product, error) {

	// Create a new product to add to Zincsearch
	product := Product{Name: generateRandomProductName(), Price: 19}

	// Convert the product to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		return nil, nil, err
	}
	// Send a POST request to the Zincsearch API to add the new product
	request, err := http.NewRequest("POST", p.GetOrigin()+"/api/"+p.Index+"/_doc", bytes.NewBuffer(productJSON))
	if err != nil {
		return nil, nil, err
	}
	request.SetBasicAuth(p.Credentials.Username, p.Credentials.Password)

	bytes, err := performRequestAndReadResponse(request)
	if err != nil {
		return nil, nil, err
	}
	return bytes, &product, nil
}

func (p ProductsServer) SearchProductsByName(name string) ([]byte, error) {
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
	request, err := http.NewRequest("POST", p.GetOrigin()+"/api/"+p.Index+"/_search", strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(p.Credentials.Username, p.Credentials.Password)

	performRequestAndReadResponse(request)
	bytes, err := performRequestAndReadResponse(request)
	if err != nil {
		return nil, err
	}
	return bytes, nil

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
