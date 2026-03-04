package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"backend/config"
	"backend/db"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Define flags
	clearFlag := flag.Bool("clear", false, "Clear all data from collections")
	seedFlag := flag.Bool("seed", true, "Seed initial data (default: true)")
	initFlag := flag.Bool("init", true, "Initialize collections (default: true)")
	helpFlag := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *helpFlag {
		fmt.Println(`
Firebase Database Initialization Tool

Usage:
  go run cmd/init-db/main.go [options]

Options:
  -init     Initialize collections (default: true)
  -seed     Seed initial data (default: true)
  -clear    Clear all data from collections (warning: destructive!)
  -help     Show this help message

Examples:
  # Initialize collections and seed data
  go run cmd/init-db/main.go

  # Only initialize collections without seeding
  go run cmd/init-db/main.go -seed=false

  # Clear all data (be careful!)
  go run cmd/init-db/main.go -clear=true

  # Only seed data (assumes collections exist)
  go run cmd/init-db/main.go -init=false
		`)
		os.Exit(0)
	}

	ctx := context.Background()

	// Create Firestore client
	client, err := config.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ Firestore connected successfully!")

	// Clear data if requested
	if *clearFlag {
		fmt.Println("\n⚠️  WARNING: About to clear all data from all collections!")
		fmt.Println("This action cannot be undone. Continue? (type 'yes' to confirm)")
		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			fmt.Println("Aborted.")
			os.Exit(0)
		}

		if err := db.ClearAllData(ctx, client); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		fmt.Println("\n✅ All data cleared!")
		return
	}

	// Initialize collections
	if *initFlag {
		if err := db.InitializeCollections(ctx, client); err != nil {
			log.Fatalf("Failed to initialize collections: %v", err)
		}
	}

	// Seed data
	if *seedFlag {
		if err := db.SeedInitialData(ctx, client); err != nil {
			log.Fatalf("Failed to seed data: %v", err)
		}
	}

	fmt.Println("\n✅ Database initialization complete!")
}
