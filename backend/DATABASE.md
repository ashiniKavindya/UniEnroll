# Firebase Database Management Guide

This directory contains tools to automatically initialize and manage your Firebase (Firestore) database without manual intervention.

## Files

- **`db/init.go`** - Core database initialization module
- **`cmd/init-db/main.go`** - Command-line tool to run initialization
- **`firestore.rules`** - Firestore security rules (deploy to Firebase Console)

## Quick Start

### 1. Initialize Database with Sample Data

Run this command to create all collections and seed initial data:

```bash
cd backend
export GOOGLE_APPLICATION_CREDENTIALS=$(pwd)/firebas.json
export FIREBASE_PROJECT_ID=unienroll-85a57
go run cmd/init-db/main.go
```

### 2. Initialize Only (No Sample Data)

If you only want to create the collection structure without seeding:

```bash
go run cmd/init-db/main.go -seed=false
```

### 3. Clear All Data (Development/Testing Only)

⚠️ **WARNING: This is destructive and cannot be undone!**

```bash
go run cmd/init-db/main.go -clear=true
```

You'll be prompted to confirm before deletion.

## Collections Created

The initialization script automatically creates these Firestore collections:

- `users` - User accounts (admin, lecturer, student)
- `students` - Student profiles
- `lecturers` - Lecturer profiles
- `faculties` - Academic faculties
- `departments` - Academic departments
- `courses` - Courses
- `modules` - Course modules
- `enrollments` - Student enrollments
- `admins` - Admin profiles

## Sample Data Included

When you run with `-seed=true`, the script populates:

- **3 Sample Users**: Admin, Lecturer, and Student accounts
- **2 Faculties**: Science and Arts
- **2 Departments**: Computer Science and Mathematics
- **2 Courses**: Introduction to Programming and Data Structures
- **2 Modules**: Python Basics and Web Development

### Login Credentials

After seeding, you can log in with:

```
Email: admin@uni.edu
Email: john@uni.edu
Email: jane@uni.edu
```

(Update the password hashes in `db/init.go` with real hashes)

## Customizing Sample Data

Edit `db/init.go` to modify the sample data:

1. **seedUsers()** - Change user accounts
2. **seedFaculties()** - Change faculties
3. **seedDepartments()** - Change departments
4. **seedCourses()** - Change courses
5. **seedModules()** - Change modules

## Deploying Firestore Rules

After testing, deploy the security rules to Firebase:

```bash
firebase deploy --only firestore:rules
```

Or using Firebase Console:
1. Go to Cloud Firestore > Rules
2. Copy the content of `firestore.rules`
3. Click "Publish"

## Adding to Your Application

To use the initialization in your Go application:

```go
import "backend/db"

// In your main.go or init function:
if err := db.InitializeCollections(ctx, client); err != nil {
    log.Fatalf("Failed to initialize collections: %v", err)
}

if err := db.SeedInitialData(ctx, client); err != nil {
    log.Fatalf("Failed to seed data: %v", err)
}
```

## Workflow

### Development
1. Run initialization with sample data
2. Test your application
3. Use `-clear=true` to reset and start fresh
4. Update sample data as needed

### Production
1. Initialize collections (without seeding)
2. Deploy Firestore rules
3. Manually add data or use admin panel
4. **Never use `-clear=true` in production!**

## Troubleshooting

### "Failed to create Firestore client"
- Check that `GOOGLE_APPLICATION_CREDENTIALS` environment variable points to your Firebase service account JSON
- Verify `FIREBASE_PROJECT_ID` is set correctly

### "Permission denied" errors when deploying rules
- Ensure your service account has the "Cloud Datastore Admin" role in Google Cloud Console
- Check that you're deploying to the correct project

### Some data didn't seed
- Check for duplicate document IDs
- Verify collection names match exactly (case-sensitive)
- Check Firebase Console for any validation errors

## Next Steps

1. **Update password hashes** in `db/init.go` with actual bcrypt hashes
2. **Customize sample data** to match your institution's structure
3. **Configure Firestore rules** based on your security requirements
4. **Add migration system** for future schema changes (optional)
