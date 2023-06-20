package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/vituchon/zincsearch-playground/zincsearch"
)

func main() {

	version, err := zincsearch.DefaultProductsServer.Version()
	if err != nil {
		fmt.Println("Error getting zyncsearch version: ", err)
		return
	}
	fmt.Printf("Connected to zyncseach server (version:'%s') at %s\n", version, zincsearch.DefaultProductsServer.GetOrigin())

	// inserta un producto
	bytes, product, err := zincsearch.DefaultProductsServer.InsertRandomProduct()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Random product(='%+v') inserted, response: ", product, string(bytes))
	}
	// si se lo busca escribiendo una palabra del producto (Respetando mayusculas o minusculas) el tipo de busqudea "match" (ver https://zincsearch-docs.zinc.dev/api/search/types/) devuelve un resultado
	bytes, err = GatherNameFromStdinAndSearchProductsByName()
	fmt.Println(string(bytes), err)
}

func GatherNameFromStdinAndSearchProductsByName() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter product name: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1] // removes the enter
	return zincsearch.DefaultProductsServer.SearchProductsByName(name)
}
