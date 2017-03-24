package main

import (
	"time"
)

type product struct {
	Sku           string `json:"sku"`
	ProductFamily string `json:"productFamily"`
	Attributes    struct {
		Servicecode                 string `json:"servicecode"`
		Location                    string `json:"location"`
		LocationType                string `json:"locationType"`
		InstanceType                string `json:"instanceType"`
		CurrentGeneration           string `json:"currentGeneration"`
		InstanceFamily              string `json:"instanceFamily"`
		Vcpu                        string `json:"vcpu"`
		PhysicalProcessor           string `json:"physicalProcessor"`
		ClockSpeed                  string `json:"clockSpeed"`
		Memory                      string `json:"memory"`
		Storage                     string `json:"storage"`
		NetworkPerformance          string `json:"networkPerformance"`
		ProcessorArchitecture       string `json:"processorArchitecture"`
		Tenancy                     string `json:"tenancy"`
		OperatingSystem             string `json:"operatingSystem"`
		LicenseModel                string `json:"licenseModel"`
		Usagetype                   string `json:"usagetype"`
		Operation                   string `json:"operation"`
		EnhancedNetworkingSupported string `json:"enhancedNetworkingSupported"`
		PreInstalledSw              string `json:"preInstalledSw"`
		ProcessorFeatures           string `json:"processorFeatures"`
	} `json:"attributes"`
}
type priceDimension struct {
	RateCode     string `json:"rateCode"`
	Description  string `json:"description"`
	BeginRange   string `json:"beginRange"`
	EndRange     string `json:"endRange"`
	Unit         string `json:"unit"`
	PricePerUnit struct {
		USD string `json:"USD"`
	} `json:"pricePerUnit"`
	AppliesTo []interface{} `json:"appliesTo"`
}

type onDemand struct {
	OfferTermCode   string                    `json:"offerTermCode"`
	Sku             string                    `json:"sku"`
	EffectiveDate   time.Time                 `json:"effectiveDate"`
	PriceDimensions map[string]priceDimension `json:"priceDimensions"`
	TermAttributes  struct {
		LeaseContractLength string `json:"LeaseContractLength"`
		OfferingClass       string `json:"OfferingClass"`
		PurchaseOption      string `json:"PurchaseOption"`
	} `json:"termAttributes"`
}

type reserved map[string]struct {
	OfferTermCode   string                    `json:"offerTermCode"`
	Sku             string                    `json:"sku"`
	EffectiveDate   time.Time                 `json:"effectiveDate"`
	PriceDimensions map[string]priceDimension `json:"priceDimensions"`
	TermAttributes  struct {
		LeaseContractLength string `json:"LeaseContractLength"`
		OfferingClass       string `json:"OfferingClass"`
		PurchaseOption      string `json:"PurchaseOption"`
	} `json:"termAttributes"`
}

type aws struct {
	FormatVersion   string              `json:"formatVersion"`
	Disclaimer      string              `json:"disclaimer"`
	OfferCode       string              `json:"offerCode"`
	Version         string              `json:"version"`
	PublicationDate time.Time           `json:"publicationDate"`
	Products        map[string]*product `json:"products"`
	Terms           struct {
		OnDemand map[string]map[string]*onDemand `json:"OnDemand"`
		Reserved map[string]*reserved            `json:"Reserved"`
	}
}
