#!/bin/bash

# Firebase Database Management Script
# Automates database initialization and management tasks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GOOGLE_APPLICATION_CREDENTIALS="${GOOGLE_APPLICATION_CREDENTIALS:-$BACKEND_DIR/firebas.json}"
FIREBASE_PROJECT_ID="${FIREBASE_PROJECT_ID:-unienroll-85a57}"

# Export for Go commands
export GOOGLE_APPLICATION_CREDENTIALS
export FIREBASE_PROJECT_ID

# Functions
show_help() {
    cat << EOF
Firebase Database Management Tool

Usage: bash db-manage.sh [command] [options]

Commands:
  init              Initialize collections and seed data (default)
  seed              Only seed sample data (assumes collections exist)
  init-empty        Initialize collections without seeding
  clear             Clear all data (destructive!)
  help              Show this help message

Environment Variables:
  GOOGLE_APPLICATION_CREDENTIALS  Path to Firebase service account JSON
  FIREBASE_PROJECT_ID             Firebase project ID

Examples:
  # Initialize everything
  bash db-manage.sh init

  # Only initialize empty collections
  bash db-manage.sh init-empty

  # Seed data only
  bash db-manage.sh seed

  # Clear all data
  bash db-manage.sh clear

  # With custom service account
  GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json bash db-manage.sh init

EOF
}

check_credentials() {
    if [ ! -f "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
        echo -e "${RED}✗ Error: Firebase credentials file not found${NC}"
        echo "  Expected: $GOOGLE_APPLICATION_CREDENTIALS"
        echo ""
        echo "Set GOOGLE_APPLICATION_CREDENTIALS or ensure firebas.json is in backend directory"
        exit 1
    fi

    if [ -z "$FIREBASE_PROJECT_ID" ]; then
        echo -e "${RED}✗ Error: FIREBASE_PROJECT_ID not set${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ Credentials verified${NC}"
}

run_init() {
    echo -e "${YELLOW}Initializing Firebase database...${NC}"
    echo "  Project: $FIREBASE_PROJECT_ID"
    echo ""
    
    cd "$BACKEND_DIR"
    go run cmd/init-db/main.go -init=true -seed=true
    
    echo ""
    echo -e "${GREEN}✅ Database initialization complete!${NC}"
}

run_seed() {
    echo -e "${YELLOW}Seeding sample data...${NC}"
    echo ""
    
    cd "$BACKEND_DIR"
    go run cmd/init-db/main.go -init=false -seed=true
    
    echo ""
    echo -e "${GREEN}✅ Sample data seeded!${NC}"
}

run_init_empty() {
    echo -e "${YELLOW}Initializing empty collections...${NC}"
    echo ""
    
    cd "$BACKEND_DIR"
    go run cmd/init-db/main.go -init=true -seed=false
    
    echo ""
    echo -e "${GREEN}✅ Collections initialized!${NC}"
}

run_clear() {
    echo -e "${RED}⚠️  WARNING: About to clear ALL data from database!${NC}"
    echo "This action CANNOT be undone!"
    echo ""
    echo "Type 'clear all' to confirm, or press Ctrl+C to cancel:"
    read -r response
    
    if [ "$response" != "clear all" ]; then
        echo "Aborted."
        exit 0
    fi
    
    echo ""
    echo -e "${YELLOW}Clearing all data...${NC}"
    
    cd "$BACKEND_DIR"
    go run cmd/init-db/main.go -clear=true
    
    echo ""
    echo -e "${GREEN}✅ All data cleared!${NC}"
}

# Main script
main() {
    check_credentials
    
    case "${1:-init}" in
        init)
            run_init
            ;;
        seed)
            run_seed
            ;;
        init-empty)
            run_init_empty
            ;;
        clear)
            run_clear
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            echo -e "${RED}Unknown command: $1${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

main "$@"
