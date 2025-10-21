package db

import (
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"filmfolk/internals/models"
)

var DB *gorm.DB


func InitDB() {

	dsn := "postgres://ravflyin@localhost:5432/filmfolk?sslmode=disable"
	/* dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"),
	os.Getenv("DB_PORT"),
) */

var err error
DB, err = gorm.Open(postgres.Open(dsn),&gorm.Config{})
if err != nil {
	log.Fatalf("Database Connection Failed ; %v",err)
}

	log.Println("Database connection successfully established.")

	log.Println("Creating ENUMs")

	CreateENUM(DB)

//DB Migration
		log.Println("Running database migrations...")
		err = DB.AutoMigrate(&models.Movie{},&models.Review{},&models.User{})
		if err != nil {
			log.Fatalf("DataBase Migration Failed :%v",err)
		}

		log.Printf("Database Migration Complete")
}



func CreateENUM(db *gorm.DB) {
	enumQueries := []string{
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN CREATE TYPE user_role AS ENUM ('user', 'moderator', 'admin'); END IF; END $$;",
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'movie_status') THEN CREATE TYPE movie_status AS ENUM ('pending_approval', 'approved', 'rejected'); END IF; END $$;",
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'review_status') THEN CREATE TYPE review_status AS ENUM ('pending_moderation', 'published', 'rejected'); END IF; END $$;",
	}

	for _, query := range enumQueries {
		if err := db.Exec(query).Error; err != nil {
			log.Fatalf("Failed to create ENUM type: %v", err)
		}
	}

}