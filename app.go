package main

import (
  "encoding/json"
  "log"
  "net/http"
  "database/sql"
  _ "github.com/lib/pq"
  "fmt"
)

var db *sql.DB

type Product struct {
  id int
  name string
  description string
  balance int
  price int
  category string
}

type Products struct {
  Products []Product
}

func main() {
  var err error

  db, err = sql.Open("postgres", "host=127.0.0.1 user=devops password=7777777 dbname=api sslmode=disable")
  if err != nil {
    panic(err)
  }

  defer db.Close()

  fmt.Println("Starting server...")

  http.HandleFunc("/v1/products/add", addProduct)
  http.HandleFunc("/v1/products/", getProducts)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func addProduct(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    http.Error(w, "Method Not Allowed", 405)
  } else {
    decoder := json.NewDecoder(r.Body)
    var g_product Product

    err := decoder.Decode(&g_product)
    if err != nil {
        panic(err)
    }
    query := fmt.Sprintf("INSERT INTO products(name, description, balance, discount, category) VALUES('%s', '%s', %d, %d, %d) RETURNING id", g_product.name, g_product.description, g_product.balance, g_product.price, g_product.category)

    fmt.Println("# INSERT QUERY: %s", query)

    rows, err := db.Query(query)
    if err != nil {
        panic(err)
    }

    for rows.Next() {
      var id int
      err = rows.Scan(&id)
      if err != nil {
        panic(err)
      }
      fmt.Fprintf(w, "{\"id\":%d}", id)
    }

  }
}

func getProducts(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "Method Not Allowed", 405)
  } else {
    w_array := Products{}

    fmt.Println("# Querying")
    rows, err := db.Query("SELECT id,name,description,category from products")
    if err != nil {
        panic(err)
    }

    for rows.Next() {
      w_product := Product{}

      err = rows.Scan(&w_product.id,&w_product.name,&w_product.description,&w_product.category)
      if err != nil {
        panic(err)
      }
      w_array.Products = append(w_array.Products, w_product)

    }

    json.NewEncoder(w).Encode(w_array)
  }
}
