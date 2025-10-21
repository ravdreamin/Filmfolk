package main 

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"filmfolk/internals/db"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No ENV file found using os")
	}

	//DB function
	db.InitDB();
}

func main() {
	router := gin.Default()
	log.Println("Server running at 8080")


	if err := router.Run(":8080");err != nil{
		log.Fatalf("Server Failed to Start :%v",err)
	}
}