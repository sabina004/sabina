package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	db *sql.DB
}

type OrderInfo struct {
	CustomerName  string
	CustomerEmail string
}

func dbConnect() server {
	db, err := sql.Open("sqlite3", "shop.db")
	if err != nil {
		log.Fatal(err)
	}

	s := server{db: db}

	return s
}

func (s *server) orderHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customerName := r.FormValue("customer_name")
	customerEmail := r.FormValue("customer_email")
	totalPrice := r.FormValue("total_price")
	totalPriceInt := 0
	_, err := fmt.Sscanf(totalPrice, "%d", &totalPriceInt)
	if err != nil {
		http.Error(w, "Invalid total price", http.StatusBadRequest)
		return
	}



	_, err = s.db.Exec("INSERT INTO orders (customer_name, customer_email) VALUES (?, ?)", customerName, customerEmail)
	if err != nil {
		http.Error(w, "Failed to insert order", http.StatusInternalServerError)
		return
	}

	orderInfo := OrderInfo{
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
	}

	outputHTML(w, "./static/postmethod.html", orderInfo)
}

func outputHTML(w http.ResponseWriter, filename string, orderInfo OrderInfo) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Fatal(err)
	}

	errExecute := t.Execute(w, orderInfo)
	if errExecute != nil {
		log.Fatal(errExecute)
	}
}

func main() {
	s := dbConnect()
	defer s.db.Close()

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	http.HandleFunc("/order", s.orderHandle)

	fmt.Println("Server running...")
	http.ListenAndServe(":8080", nil)
}
