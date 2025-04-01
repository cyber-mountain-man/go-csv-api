package main

// These are packages from the Go standard library and a third-party library.
// They provide built-in functions for tasks like file access, HTTP, time, and working with CSV data.
import (
	"encoding/csv" // For reading data from CSV files
	"fmt"          // For creating formatted strings
	"net/http"     // For HTTP status codes like 200 OK
	"os"           // For opening files
	"strconv"      // For converting strings to integers or floats
	"time"         // For getting the current time (used in API responses)

	"github.com/gin-gonic/gin" // Gin is a popular third-party package for building web APIs
)

// Item is a custom data structure (a "struct") that represents one row of data from the CSV file.
// Each field corresponds to a column in the CSV, and the tags (e.g. `json:"year"`) tell Gin how to name the fields in the JSON output.
type Item struct {
	Year            int     `json:"year"`
	Month           int     `json:"month"`
	Supplier        string  `json:"supplier"`
	ItemCode        string  `json:"item_code"`
	ItemDescription string  `json:"item_description"`
	ItemType        string  `json:"item_type"`
	RetailSales     float64 `json:"retail_sales"`
	RetailTransfers float64 `json:"retail_transfers"`
	WarehouseSales  float64 `json:"warehouse_sales"`
}

// items is a global slice (a dynamic list) that will hold all Item records from the CSV file.
var items []Item

// paginate is a helper function that returns a "page" of results from a slice of Items.
// It uses the `offset` (where to start) and `limit` (how many to return) parameters.
func paginate(data []Item, offset, limit int) []Item {
	start := offset
	if start > len(data) {
		start = len(data) // avoid out-of-range errors
	}
	end := start + limit
	if end > len(data) {
		end = len(data) // make sure we don't go past the end
	}
	return data[start:end]
}

// loadCSV loads data from the CSV file and parses each row into an Item struct.
// It returns an error if anything goes wrong (e.g., file not found).
func loadCSV(filepath string) error {
	file, err := os.Open(filepath) // Open the CSV file
	if err != nil {
		return err // Return the error to the caller if we can't open it
	}
	defer file.Close() // Make sure the file gets closed after this function is done

	reader := csv.NewReader(file)

	// The first row of the CSV usually contains headers, so we skip it
	_, _ = reader.Read()

	// Read the remaining rows from the file
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Loop through each record (row) and convert it into an Item
	for _, r := range records {
		// Convert strings to appropriate types
		year, _ := strconv.Atoi(r[0])
		month, _ := strconv.Atoi(r[1])
		retailSales, _ := strconv.ParseFloat(r[6], 64)
		retailTransfers, _ := strconv.ParseFloat(r[7], 64)
		warehouseSales, _ := strconv.ParseFloat(r[8], 64)

		// Create an Item from the parsed values
		item := Item{
			Year:            year,
			Month:           month,
			Supplier:        r[2],
			ItemCode:        r[3],
			ItemDescription: r[4],
			ItemType:        r[5],
			RetailSales:     retailSales,
			RetailTransfers: retailTransfers,
			WarehouseSales:  warehouseSales,
		}

		// Add this item to our global list
		items = append(items, item)
	}
	return nil
}

// getAllItems handles GET requests to /items
// It returns a paginated list of all items.
func getAllItems(c *gin.Context) {
	// Read limit and offset from the query string (e.g., ?limit=10&offset=0)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))   // default to 10
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))  // default to 0

	// Use the paginate helper function to get a slice of items
	pagedItems := paginate(items, offset, limit)

	// Send a JSON response with metadata and the paged data
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339), // current time in a standard format
		"count":     len(pagedItems),                 // number of items returned
		"total":     len(items),                      // total number of items
		"offset":    offset,
		"limit":     limit,
		"data":      pagedItems,
	})
}

// getItemsByType handles GET requests to /items/type?type=WINE
// It filters items by item type and returns a paginated response
func getItemsByType(c *gin.Context) {
	itemType := c.Query("type") // read the `type` query parameter
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var filtered []Item
	for _, item := range items {
		if item.ItemType == itemType {
			filtered = append(filtered, item) // keep matching items
		}
	}

	paged := paginate(filtered, offset, limit)

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"count":     len(paged),
		"total":     len(filtered),
		"offset":    offset,
		"limit":     limit,
		"message":   fmt.Sprintf("Found %d items of type %s", len(filtered), itemType),
		"data":      paged,
	})
}

// getItemsBySupplier handles GET requests to /supplier/:supplier
// It filters items by supplier name and returns a paginated response
func getItemsBySupplier(c *gin.Context) {
	supplier := c.Param("supplier") // Get the path parameter from the URL
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var filtered []Item
	for _, item := range items {
		if item.Supplier == supplier {
			filtered = append(filtered, item)
		}
	}

	paged := paginate(filtered, offset, limit)

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"timestamp": time.Now().Format(time.RFC3339),
		"count":     len(paged),
		"total":     len(filtered),
		"offset":    offset,
		"limit":     limit,
		"message":   fmt.Sprintf("Found %d items from supplier %s", len(filtered), supplier),
		"data":      paged,
	})
}

// main is the entry point of the program.
// This is where everything is put together: load data, set up routes, and start the server.
func main() {
	// Step 1: Load the data from the CSV file
	err := loadCSV("data/Warehouse_and_Retail_Sales.csv")
	if err != nil {
		panic(err) // Stop the program if we couldn't load the file
	}

	// Step 2: Create a new Gin router
	r := gin.Default()

	// Step 3: Define our API routes and connect them to handler functions
	r.GET("/items", getAllItems)                     // List all items with pagination
	r.GET("/items/type", getItemsByType)             // Filter by item type
	r.GET("/supplier/:supplier", getItemsBySupplier) // Filter by supplier

	// Step 4: Start the server on port 8080
	r.Run(":8080")
}
