package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Import PostgreSQL driver
	"github.com/rs/cors"
)

// Todo represents a single todo item
type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// Global variable for the database connection
var db *sql.DB

// Initialize the PostgreSQL database connection
func init() {
	var err error
	db, err = sql.Open("postgres", "postgresql://demo_du6c_user:1f2qZEeQmeYIpBR5EOI5G2m4uoPDIdC6@dpg-ctls88aj1k6c73d2gti0-a.oregon-postgres.render.com/demo_du6c")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Ensure the connection is valid
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}
}

// Get all todos from the database
func getTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, text, completed FROM todos")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying todos: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.Completed); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning todos: %v", err), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding todos: %v", err), http.StatusInternalServerError)
		return
	}
}

// Add a new todo to the database
func addTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding todo: %v", err), http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow("INSERT INTO todos(text, completed) VALUES($1, $2) RETURNING id", newTodo.Text, newTodo.Completed).Scan(&id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting new todo: %v", err), http.StatusInternalServerError)
		return
	}

	newTodo.ID = id

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newTodo); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding new todo: %v", err), http.StatusInternalServerError)
		return
	}
}

// Update a todo's completion status in the database
func updateTodo(w http.ResponseWriter, r *http.Request) {
	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding todo: %v", err), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE todos SET completed=$1 WHERE id=$2", updatedTodo.Completed, updatedTodo.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating todo: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTodo); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding updated todo: %v", err), http.StatusInternalServerError)
		return
	}
}

// Delete a todo from the database
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	var todoToDelete Todo
	if err := json.NewDecoder(r.Body).Decode(&todoToDelete); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding todo: %v", err), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM todos WHERE id=$1", todoToDelete.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting todo: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoToDelete); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding deleted todo: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Set up CORS options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                            // Allow all origins (you can restrict this to a specific domain)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, // Allowed HTTP methods
		AllowedHeaders: []string{"Content-Type"},                 // Allowed headers
	})

	// Set up routes
	http.HandleFunc("/todos", getTodos)          // GET /todos - fetch all todos
	http.HandleFunc("/todos/add", addTodo)       // POST /todos/add - add a new todo
	http.HandleFunc("/todos/update", updateTodo) // PUT /todos/update - update a todo
	http.HandleFunc("/todos/delete", deleteTodo) // DELETE /todos/delete - delete a todo

	// Wrap the handler with CORS support
	handler := corsHandler.Handler(http.DefaultServeMux)

	// Start the server
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
