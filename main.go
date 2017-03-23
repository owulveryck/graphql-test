package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
	"net/http"
	"os"
)

var query string
var dbFile string

func init() {
	flag.StringVar(&query, "query", "{}", "query to ask server for data")
	flag.StringVar(&dbFile, "db", "", "DB Path")
}

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"location": &graphql.Field{
			Type: graphql.String,
		},
		"sku": &graphql.Field{
			Type: graphql.String,
		},
		"operatingSystem": &graphql.Field{
			Type: graphql.String,
		},
		"instanceType": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func main() {
	flag.Parse()
	var objects aws
	if dbFile == "" {
		log.Println("Please provide a dbFile")
		return
	}
	file, err := os.Open(dbFile) // For read access.
	if err != nil {
		// Cannot open file... Let's create it
		resp, err := http.Get("https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/20170224022054/index.json")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&objects)
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(dbFile)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		enc := gob.NewEncoder(file)
		err = enc.Encode(objects)
		if err != nil {
			log.Fatal("encode error:", err)
		}
	} else {
		defer file.Close()
		dec := gob.NewDecoder(file)
		err = dec.Decode(&objects)
		if err != nil {
			log.Fatal("decode error:", err)
		}
	}
	// Let's graphql
	fields := graphql.Fields{
		"products": &graphql.Field{
			Type: graphql.NewList(productType),
			Args: graphql.FieldConfigArgument{
				"sku": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				type myproduct struct {
					Sku             string `json:"sku"`
					Location        string `json:"location"`
					InstanceType    string `json:"instanceType"`
					OperatingSystem string `json:"operatingSystem"`
				}

				var prds []*myproduct
				if sku, skuok := p.Args["sku"].(string); skuok {
					prds = append(prds, &myproduct{
						Sku:             objects.Products[sku].Sku,
						Location:        objects.Products[sku].Attributes.Location,
						InstanceType:    objects.Products[sku].Attributes.InstanceType,
						OperatingSystem: objects.Products[sku].Attributes.OperatingSystem,
					})

				} else {

					for _, prd := range objects.Products {
						prds = append(prds, &myproduct{
							Sku:             prd.Sku,
							Location:        prd.Attributes.Location,
							InstanceType:    prd.Attributes.InstanceType,
							OperatingSystem: prd.Attributes.OperatingSystem,
						})
					}
				}
				return prds, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if r.HasErrors() {
		log.Fatalf("Failed due to errors: %v\n", r.Errors)
	}

	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s", rJSON)

}