package main

import (
	"encoding/gob"
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"os"
)

// Configuration ...
type Configuration struct {
	DB         string `default:"db"`
	ListenAddr string `default:":8080"`
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
		"offers": &graphql.Field{
			Type: graphql.NewList(offerType),
		},
	},
})

var offerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Offer",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"LeaseContractLenght": &graphql.Field{
			Type: graphql.String,
		},
		"PurchaseOption": &graphql.Field{
			Type: graphql.String,
		},
		"OfferingClass": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func main() {
	var config Configuration
	err := envconfig.Process("graphql", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	var objects aws
	if config.DB == "" {
		log.Println("Please provide a config.DB")
		return
	}
	file, err := os.Open(config.DB) // For read access.
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
		file, err := os.Create(config.DB)
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

				type myoffer struct {
					Type string `json:"type"`
					Code string `json:"code"`
				}
				type myproduct struct {
					Sku             string    `json:"sku"`
					Location        string    `json:"location"`
					InstanceType    string    `json:"instanceType"`
					OperatingSystem string    `json:"operatingSystem"`
					Offer           []myoffer `json:"offers"`
				}

				var prds []*myproduct
				if sku, skuok := p.Args["sku"].(string); skuok {
					var odtc string
					for _, od := range objects.Terms.OnDemand[sku] {
						odtc = od.OfferTermCode
					}
					prds = append(prds, &myproduct{
						Sku:             objects.Products[sku].Sku,
						Location:        objects.Products[sku].Attributes.Location,
						InstanceType:    objects.Products[sku].Attributes.InstanceType,
						OperatingSystem: objects.Products[sku].Attributes.OperatingSystem,
						Offer: []myoffer{
							myoffer{
								Type: "OnDemand",
								Code: odtc,
							},
						},
					})

				} else {

					for _, prd := range objects.Products {
						var odtc string
						for _, od := range objects.Terms.OnDemand[sku] {
							odtc = od.OfferTermCode
						}
						prds = append(prds, &myproduct{
							Sku:             prd.Sku,
							Location:        prd.Attributes.Location,
							InstanceType:    prd.Attributes.InstanceType,
							OperatingSystem: prd.Attributes.OperatingSystem,
							Offer: []myoffer{
								myoffer{
									Type: "OnDemand",
									Code: odtc,
								},
							},
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
	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	// serve HTTP
	// static file server to serve Graphiql in-browser editor
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/graphql", h)
	http.Handle("/", fs)
	log.Println("Listening on:", config.ListenAddr)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
