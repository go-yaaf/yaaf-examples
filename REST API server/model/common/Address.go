package model

// Address model represents an address
// @Data
type Address struct {
	Street  string `json:"street"`  // Street address
	City    string `json:"city"`    // City
	State   string `json:"state"`   // State (if applicable)
	ZipCode string `json:"zipCode"` // Local zip code (postal cod)
	Country string `json:"country"` // Country name
}
