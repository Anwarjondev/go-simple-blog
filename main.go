package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/lib/pq"
)

var templates = template.Must(template.ParseGlob("views/**/*"))

func DBConnection() (*sql.DB, error) {
	User := "postgres"
	Password := "1234"
	Host := os.Getenv("DB_HOST")
	Port := "5432"
	Database := "go-simple-blog"

	if User == "" || Password == "" || Host == "" || Port == "" || Database == "" {
		log.Fatal("Database environment variables not set")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Host, Port, User, Password, Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Create table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL DEFAULT 'Empty article :P',
		content TEXT NOT NULL,
		author VARCHAR(50) NOT NULL DEFAULT 'Neil Amnstrong',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return db, nil
}

func main() {
	// Initialize database connection
	db, err := DBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Register routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, db)
	})
	http.HandleFunc("/show", func(w http.ResponseWriter, r *http.Request) {
		Show(w, r, db)
	})
	http.HandleFunc("/create", Create)
	http.HandleFunc("/store", func(w http.ResponseWriter, r *http.Request) {
		Store(w, r, db)
	})
	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		Edit(w, r, db)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		Update(w, r, db)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		Delete(w, r, db)
	})

	log.Println("Server starting on port 8000...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

type Post struct {
	Id        int
	Title     string
	Content   string
	Author    string
	CreatedAt string
}

func Home(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	posts, err := db.Query("SELECT * FROM posts")
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		log.Printf("Error fetching posts: %v", err)
		return
	}
	defer posts.Close()

	postArrays := []Post{}
	for posts.Next() {
		var post Post
		var created_at string
		err = posts.Scan(&post.Id, &post.Title, &post.Content, &post.Author, &created_at)
		if err != nil {
			http.Error(w, "Error scanning post", http.StatusInternalServerError)
			log.Printf("Error scanning post: %v", err)
			return
		}
		post.CreatedAt = FormatDate(created_at)
		if len(post.Content) > 450 {
			post.Content = post.Content[:450] + "..."
		}
		postArrays = append(postArrays, post)
	}

	if err = posts.Err(); err != nil {
		http.Error(w, "Error iterating posts", http.StatusInternalServerError)
		log.Printf("Error iterating posts: %v", err)
		return
	}

	if err := templates.ExecuteTemplate(w, "home", postArrays); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func Show(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	var post Post
	var created_at string
	err := db.QueryRow("SELECT * FROM posts WHERE id = $1", id).Scan(
		&post.Id, &post.Title, &post.Content, &post.Author, &created_at)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching post", http.StatusInternalServerError)
			log.Printf("Error fetching post: %v", err)
		}
		return
	}

	post.CreatedAt = FormatDate(created_at)
	if err := templates.ExecuteTemplate(w, "show", post); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "create", nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func Store(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	author := r.FormValue("author")

	if title == "" || content == "" || author == "" {
		http.Redirect(w, r, "/create", http.StatusSeeOther)
		return
	}

	_, err := db.Exec("INSERT INTO posts (title, content, author) VALUES ($1, $2, $3)",
		title, content, author)
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		log.Printf("Error creating post: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Edit(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	var post Post
	var created_at string
	err := db.QueryRow("SELECT * FROM posts WHERE id = $1", id).Scan(
		&post.Id, &post.Title, &post.Content, &post.Author, &created_at)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching post", http.StatusInternalServerError)
			log.Printf("Error fetching post: %v", err)
		}
		return
	}

	post.CreatedAt = FormatDate(created_at)
	if err := templates.ExecuteTemplate(w, "edit", post); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func Update(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	author := r.FormValue("author")

	if title == "" || content == "" || author == "" {
		http.Redirect(w, r, "/edit?id="+id, http.StatusSeeOther)
		return
	}

	_, err := db.Exec("UPDATE posts SET title = $1, content = $2, author = $3 WHERE id = $4",
		title, content, author, id)
	if err != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		log.Printf("Error updating post: %v", err)
		return
	}

	http.Redirect(w, r, "/show?id="+id, http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		log.Printf("Error deleting post: %v", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FormatDate(date string) string {
	// Try parsing ISO 8601 format first
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		// If that fails, try the old format
		t, err = time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			return date
		}
	}
	return t.Format("January 2, 2006")
}
