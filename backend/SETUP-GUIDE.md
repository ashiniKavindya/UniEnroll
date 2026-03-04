# Firebase Database Automation - Quick Reference

## Overview

You now have an automated system to generate and manage Firebase (Firestore) collections and data without manual intervention. Here are your options:

---

## Method 1: Command-Line Tools (Easiest)

### On macOS/Linux:
```bash
cd backend
export GOOGLE_APPLICATION_CREDENTIALS=$(pwd)/firebas.json
export FIREBASE_PROJECT_ID=unienroll-85a57

# Initialize and seed
bash db-manage.sh init

# Or without sample data
bash db-manage.sh init-empty

# Or just seed
bash db-manage.sh seed
```

### On Windows (PowerShell or Command Prompt):
```bash
cd backend
# Set environment variables first if needed
# $env:GOOGLE_APPLICATION_CREDENTIALS = "$PWD\firebas.json"
# $env:FIREBASE_PROJECT_ID = "unienroll-85a57"

# Using batch script
db-manage.bat init
db-manage.bat init-empty
db-manage.bat seed
```

### Using Go directly:
```bash
cd backend
go run cmd/init-db/main.go
go run cmd/init-db/main.go -seed=false
go run cmd/init-db/main.go -clear=true
```

---

## Method 2: Programmatically in Your App

Add this to your `main.go` or a startup function:

```go
import "backend/db"

// In your main() function after creating the Firestore client:
ctx := context.Background()
client, err := config.NewClient(ctx)
if err != nil {
    log.Fatalf("Failed to create Firestore client: %v", err)
}
defer client.Close()

// Initialize collections
if err := db.InitializeCollections(ctx, client); err != nil {
    log.Fatalf("Failed to initialize collections: %v", err)
}

// Optionally seed data
if os.Getenv("SEED_DB") == "true" {
    if err := db.SeedInitialData(ctx, client); err != nil {
        log.Fatalf("Failed to seed data: %v", err)
    }
}
```

Then run with:
```bash
SEED_DB=true go run main.go
```

---

## Method 3: Database Initialization During App Startup

Modify your `handlers/init.go` or create a startup endpoint:

```go
package handlers

import (
    "github.com/gin-gonic/gin"
    "backend/db"
)

// InitDB initializes the database (call this once on startup)
func InitDB(c *gin.Context) {
    ctx := c.Request.Context()
    client := c.MustGet("firestore").(*firestore.Client)
    
    if err := db.InitializeCollections(ctx, client); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "Database initialized successfully"})
}
```

Then call:
```bash
curl -X POST http://localhost:8080/api/init-db
```

---

## Commands Summary

| Command | Effect |
|---------|--------|
| `bash db-manage.sh init` | Create all collections + seed sample data |
| `bash db-manage.sh init-empty` | Create all collections (no data) |
| `bash db-manage.sh seed` | Add sample data to existing collections |
| `bash db-manage.sh clear` | Delete all data (⚠️ destructive) |
| `go run cmd/init-db/main.go` | Same as `init` |
| `go run cmd/init-db/main.go -seed=false` | Same as `init-empty` |
| `go run cmd/init-db/main.go -clear=true` | Same as `clear` |

---

## Collections Created

✓ users  
✓ students  
✓ lecturers  
✓ faculties  
✓ departments  
✓ courses  
✓ modules  
✓ enrollments  
✓ admins  

---

## Sample Data Included

When seeding, you get:

**Users:**
- Admin user (`admin@uni.edu`)
- Lecturer (`john@uni.edu`)
- Student (`jane@uni.edu`)

**Organizational:**
- 2 Faculties (Science, Arts)
- 2 Departments (CS, Math)
- 2 Courses
- 2 Modules

---

## Customization

Edit `db/init.go` to:
- Change admin credentials
- Add more users, departments, courses
- Modify module structures
- Add additional collection types

Each `seed*()` function handles one collection type.

---

## Firestore Security Rules

The `firestore.rules` file provides:
- Admin full access
- Users can see their own profiles
- Students can see their enrollments
- Lecturers/students can access courses and modules
- Public read for faculties/departments

To deploy:
```bash
firebase deploy --only firestore:rules
```

---

## Common Workflows

### Development Setup
```bash
# Day 1: Fresh start
bash db-manage.sh init

# During testing: Reset database
bash db-manage.sh clear
bash db-manage.sh init
```

### Production Deployment
```bash
# Initialize empty collections
bash db-manage.sh init-empty

# Manually add data OR use admin interface
# Deploy rules
firebase deploy --only firestore:rules
```

### Update Collections Later
```bash
# Edit db/init.go with new structure
# Backup data first!
bash db-manage.sh clear
bash db-manage.sh init
```

---

## Troubleshooting

**Problem:** "Failed to create Firestore client"
- ✓ Check `GOOGLE_APPLICATION_CREDENTIALS` path
- ✓ Verify `FIREBASE_PROJECT_ID` matches your project

**Problem:** "Permission denied"  
- ✓ Service account needs "Cloud Datastore Admin" role in Google Cloud Console
- ✓ Check IAM settings in Firebase project

**Problem:** Data not appearing
- ✓ Wait a moment (Firestore eventual consistency)
- ✓ Check Firebase Console > Firestore to verify collections exist
- ✓ Check Go application logs for errors

---

## Files Reference

| File | Purpose |
|------|---------|
| `db/init.go` | Core initialization module |
| `cmd/init-db/main.go` | CLI entry point |
| `db-manage.sh` | Linux/macOS automated script |
| `db-manage.bat` | Windows automated script |
| `firestore.rules` | Security rules (deploy to Firebase) |
| `DATABASE.md` | Detailed documentation |

---

## Next Steps

1. ✅ Create collections: `bash db-manage.sh init`
2. ✅ Test your app
3. ✅ Update `db/init.go` with real data
4. ✅ Deploy Firestore rules: `firebase deploy --only firestore:rules`
5. ✅ Use in production!
