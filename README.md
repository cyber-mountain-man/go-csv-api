```markdown
# go-csv-api

A lightweight REST API built with **Go** and **Gin** that loads warehouse and retail sales data from a CSV file and serves it through paginated, filterable endpoints.

---

## 🚀 Features

- ✅ Load and parse a structured CSV file  
- ✅ Serve data as a REST API using Gin  
- ✅ Filter by item type or supplier  
- ✅ Supports pagination via `limit` and `offset` query parameters  
- ✅ Clean, JSON-formatted responses with metadata  

---

## 📦 API Endpoints

### 🔹 Get All Items (with pagination)
```
GET /items?limit=10&offset=0
```

### 🔹 Filter by Item Type (e.g., WINE)
```
GET /items/type?type=WINE&limit=5
```

### 🔹 Filter by Supplier (URL encoded if it contains spaces)
```
GET /supplier/PWSWN%20INC?limit=3
```

---

## 📁 Project Structure

```
go-csv-api/
├── data/
│   └── Warehouse_and_Retail_Sales.csv
├── main.go
├── go.mod
├── .gitignore
└── README.md
```

> 💡 `Warehouse_and_Retail_Sales.csv` is listed in `.gitignore` by default and is not pushed to GitHub. Add your own data file inside the `data/` folder.

---

## ⚙️ Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/go-csv-api.git
cd go-csv-api
```

### 2. Download Dependencies

```bash
go mod tidy
```

### 3. Run the App

```bash
go run main.go
```

> The server will start at `http://localhost:8080`

---

## 🛠 Requirements

- [Go](https://golang.org/dl/) 1.18 or higher  
- [Gin](https://github.com/gin-gonic/gin) web framework  

Install Gin if needed:

```bash
go get github.com/gin-gonic/gin
```

---

## 📊 Example JSON Response

```json
{
  "status": "success",
  "timestamp": "2025-04-01T19:15:00Z",
  "count": 5,
  "total": 89,
  "offset": 10,
  "limit": 5,
  "message": "Found 89 items of type WINE",
  "data": [
    {
      "year": 2020,
      "month": 1,
      "supplier": "PWSWN INC",
      "item_code": "100024",
      "item_description": "SOME PRODUCT",
      "item_type": "WINE",
      "retail_sales": 0.82,
      "retail_transfers": 0,
      "warehouse_sales": 4
    }
  ]
}
```

---

## 🧾 License

This project is open source and available under the [MIT License](LICENSE).

---

## 🙌 Acknowledgments

- Built with [Go](https://golang.org/) and [Gin](https://github.com/gin-gonic/gin)  
- CSV format inspired by warehouse/retail reporting datasets  

---

## 💡 Want to Contribute?

Pull requests are welcome! Feel free to fork the repo, submit enhancements, or suggest new features.

---