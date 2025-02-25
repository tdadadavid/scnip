package main

import (
	"fmt"
	"sort"
	"time"
)

// All the sorting methods names
const (
	BY_PRICE_HIGH_TO_LOW = "Price (Low to High)"
	BY_POPULARITY        = "Popularity (Sales per View)"
	BY_NEWSET            = "Newest First"
	BY_APLHABETS         = "Alphabetical (Z to A)"
)

// Product represents a product in the catalog
type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Created    string  `json:"created"`
	SalesCount int     `json:"sales_count"`
	ViewsCount int     `json:"views_count"`
}

// ProductSorter defines the interface for sorting products
type ProductSorter interface {
	Sort(products []Product) []Product
	GetName() string
}

// BaseSorter implements common functionality for all sorters
// this implements the ProductSorter interface.
type BaseSorter struct {
	name      string
	sortLogic func(p1, p2 *Product) bool
}

// Sort sorts the products based on the sorter's logic
func (s *BaseSorter) Sort(products []Product) []Product {
	result := make([]Product, len(products))
	copy(result, products)

	sort.Slice(result, func(i, j int) bool {
		return s.sortLogic(&result[i], &result[j])
	})

	return result
}

// GetName returns the name of the sorter
func (s *BaseSorter) GetName() string {
	return s.name
}

// PriceSorter sorts products by price (ascending)
// by embedding the BaseSorter, I have made this PriceSorter implement the ProductSorter interface.
type PriceSorter struct {
	BaseSorter
}

// NewPriceSorter creates a new price sorter
func NewPriceSorter() *PriceSorter {
	return &PriceSorter{
		BaseSorter{
			name: BY_PRICE_HIGH_TO_LOW,
			sortLogic: func(p1, p2 *Product) bool {
				return p1.Price < p2.Price
			},
		},
	}
}

// SalesPerViewSorter sorts products by sales per view ratio (descending)
type SalesPerViewSorter struct {
	BaseSorter
}

// NewSalesPerViewSorter creates a new sales per view sorter
func NewSalesPerViewSorter() *SalesPerViewSorter {
	return &SalesPerViewSorter{
		BaseSorter{
			name: BY_POPULARITY,
			sortLogic: func(p1, p2 *Product) bool {
				ratio1 := float64(p1.SalesCount) / float64(p1.ViewsCount)
				ratio2 := float64(p2.SalesCount) / float64(p2.ViewsCount)
				return ratio1 > ratio2 // Descending order for popularity
			},
		},
	}
}

// NewestFirstSorter sorts products by creation date (newest first)
type NewestFirstSorter struct {
	BaseSorter
}

// NewNewestFirstSorter creates a new sorter for newest products first
func NewNewestFirstSorter() *NewestFirstSorter {
	return &NewestFirstSorter{
		BaseSorter{
			name: BY_NEWSET,
			sortLogic: func(p1, p2 *Product) bool {
				date1, _ := time.Parse("2006-01-02", p1.Created) // ignored error 🫠 for this test
				date2, _ := time.Parse("2006-01-02", p2.Created)
				return date1.After(date2)
			},
		},
	}
}

// SorterRegistry manages available sorters
type SorterRegistry struct {
	sorters map[string]ProductSorter
}

// NewSorterRegistry creates a new sorter registry with default sorters
func NewSorterRegistry() *SorterRegistry {
	registry := &SorterRegistry{
		sorters: make(map[string]ProductSorter),
	}

	// Register default sorters
	registry.RegisterSorter(NewPriceSorter())
	registry.RegisterSorter(NewSalesPerViewSorter())
	registry.RegisterSorter(NewNewestFirstSorter())

	return registry
}

// RegisterSorter adds a new sorter to the registry
func (r *SorterRegistry) RegisterSorter(sorter ProductSorter) {
	r.sorters[sorter.GetName()] = sorter
}

// GetSorter retrieves a sorter by name
func (r *SorterRegistry) GetSorter(name string) (ProductSorter, bool) {
	sorter, exists := r.sorters[name]
	return sorter, exists
}

// GetAvailableSorters returns a list of available sorter names
func (r *SorterRegistry) GetAvailableSorters() []string {
	var names []string
	for name := range r.sorters {
		names = append(names, name)
	}
	return names
}

// ProductCatalog manages the product inventory and sorting
type ProductCatalog struct {
	products []Product
	registry *SorterRegistry
}

// NewProductCatalog creates a new product catalog
func NewProductCatalog(products []Product) *ProductCatalog {
	return &ProductCatalog{
		products: products,
		registry: NewSorterRegistry(),
	}
}

// GetSortedProducts returns products sorted according to the specified method
func (c *ProductCatalog) GetSortedProducts(sorterName string) ([]Product, error) {
	sorter, exists := c.registry.GetSorter(sorterName)
	if !exists {
		return nil, fmt.Errorf("sorter '%s' not found", sorterName)
	}

	return sorter.Sort(c.products), nil
}

// AddSortingLogic allows adding a new sorter to the catalog
func (c *ProductCatalog) AddSortingLogic(sorter ProductSorter) {
	c.registry.RegisterSorter(sorter)
}

func main() {
	// Testcases.
	products := []Product{
		{
			ID:         1,
			Name:       "Alabaster Table",
			Price:      12.99,
			Created:    "2019-01-04",
			SalesCount: 32,
			ViewsCount: 730,
		},
		{
			ID:         2,
			Name:       "Zebra Table",
			Price:      44.49,
			Created:    "2012-01-04",
			SalesCount: 301,
			ViewsCount: 3279,
		},
		{
			ID:         3,
			Name:       "Coffee Table",
			Price:      10.00,
			Created:    "2014-05-28",
			SalesCount: 1048,
			ViewsCount: 20123,
		},
	}

	catalog := NewProductCatalog(products)

	// we use existing sorting methods
	fmt.Println("Available sorting methods:")
	for _, name := range catalog.registry.GetAvailableSorters() {
		fmt.Printf("- %s\n", name)
	}
	fmt.Println()

	// Demonstrate each sorting method
	sorterNames := []string{BY_PRICE_HIGH_TO_LOW, BY_POPULARITY, BY_NEWSET}
	for _, name := range sorterNames {
		// this sorts the product
		sortedProducts, err := catalog.GetSortedProducts(name)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		fmt.Printf("Products sorted by %s:\n", name)
		// ranging over the sorted product to print it out
		// One way I can extend this is for each product to have its own Display Method (thinking about it in API level)
		// Yeah but this is just simulating each display for each sorting method.
		for _, p := range sortedProducts {
			switch name {
			case BY_PRICE_HIGH_TO_LOW:
				fmt.Printf("- %s: $%.2f\n", p.Name, p.Price) // so this is the display for Price from high - low.
			case BY_POPULARITY:
				ratio := float64(p.SalesCount) / float64(p.ViewsCount)
				fmt.Printf("- %s: %.5f (Sales: %d, Views: %d)\n", p.Name, ratio, p.SalesCount, p.ViewsCount)
			case BY_NEWSET:
				fmt.Printf("- %s: %s\n", p.Name, p.Created)
			}
		}
		fmt.Println()
	}

	// Creating and registering a custom sorter (by a different team) allowing the Product guy to extend it.
	fmt.Println("Adding a custom sorter: Alphabetical...")
	alphabeticalSorter := &BaseSorter{
		name: BY_APLHABETS,
		sortLogic: func(p1, p2 *Product) bool {
			return p2.Name < p1.Name // here I am doing reverse alphabets Z-A
		},
	}
	catalog.AddSortingLogic(alphabeticalSorter)

	// Use the new custom sorter
	sortedAlphabetically, _ := catalog.GetSortedProducts(BY_APLHABETS)
	fmt.Printf("Products sorted using %s method:\n", BY_APLHABETS)
	for _, p := range sortedAlphabetically {
		fmt.Printf("- %s\n", p.Name)
	}
}

//NOTE: I didn't not handle all the error cases, in real world scenarios they will be taken care of.
