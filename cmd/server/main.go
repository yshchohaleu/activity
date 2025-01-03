package main

import (
	"context"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"activity/internal/auth"
	"activity/internal/graph"
	"activity/internal/migrations"
	"activity/internal/models"
)

func main() {
	log.Println("Starting server...")

	// Log the credentials file path
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	log.Printf("Loading Firebase credentials from: %s", credPath)

	// Check if file exists
	if _, err := os.Stat(credPath); os.IsNotExist(err) {
		log.Fatalf("Firebase credentials file does not exist at path: %s", credPath)
	}

	// Initialize Firebase with credentials
	opt := option.WithCredentialsFile(credPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("Firebase initialization error details: %+v", err)
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize Firebase auth client
	_, err = app.Auth(context.Background())
	if err != nil {
		log.Printf("Firebase Auth client error details: %+v", err)
		log.Fatalf("Failed to create Firebase Auth client: %v", err)
	}
	log.Println("Firebase initialized successfully")

	// Initialize database
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{}, &models.Activity{})

	// Run data migrations
	migrations.EnsureDefaultUser(db)

	// Initialize GraphQL server
	resolver := graph.NewResolver(db, app)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Setup routes
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.FirebaseAuthMiddleware(app, srv))

	log.Printf("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
