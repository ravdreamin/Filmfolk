package main 

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"filmfolk/internals/routes"
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

	v1 := router.Group("/api/v1")
	{
		routes.SetupAuthRoutes(v1)
	}




	if err := router.Run(":8080");err != nil{
		log.Fatalf("Server Failed to Start :%v",err)
	}
}