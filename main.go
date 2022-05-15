package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "example/go-web-service-gin/docs"
)

// swagger:params postAlbum
type albumRequest struct {
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type albumResponse struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func seed() {
	db, err := sql.Open("mysql", "root:my-secret-pw@(localhost:3306)/GoTest?parseTime=true")

	if err != nil {
		log.Println(err)
		return
	}

	query := `
		CREATE TABLE IF NOT EXISTS album(
			ID INT AUTO_INCREMENT PRIMARY KEY,
			Title text,
			Artist text,
			Price decimal(19,2)
		)
	`

	if _, err := db.Exec(query); err != nil {
		log.Println(err)
	}

}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} albumResponse
// @Failure 500 {object} object{message=string}
// @Router /album/ [get]
func getAlbums(c *gin.Context) {
	db, err := sql.Open("mysql", "root:my-secret-pw@(localhost:3306)/GoTest?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	query := "SELECT * FROM album"

	r, err := db.Query(query)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	defer r.Close()
	albums := []albumResponse{}

	for r.Next() {
		var a albumResponse
		if err := r.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		albums = append(albums, a)
	}

	c.IndentedJSON(http.StatusOK, albums)

	db.Close()

}

// @BasePath /api/v1
// @Summary post album
// @Schemes
// @Description Post Album
// @Tags example
// @Accept json
// @Produce json
// @Param Album body albumRequest true "album"
// @Success 201 {object} object{success=bool}
// @Failure 500 {object} object{message=string}
// @Router /album/ [post]
func postAlbum(c *gin.Context) {
	db, err := sql.Open("mysql", "root:my-secret-pw@(localhost:3306)/GoTest?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	var a albumRequest
	if err := c.BindJSON(&a); err != nil {
		log.Println("Could not bind JSON")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not add"})
		return
	}

	r, err := db.Exec("INSERT INTO album(Title, Artist, Price) VALUES(?,?,?)", a.Title, a.Artist, a.Price)
	if err != nil {
		log.Printf("could not insert; %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not add"})
	}

	id, err := r.RowsAffected()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not add"})
		return
	}
	db.Close()
	c.IndentedJSON(http.StatusCreated, gin.H{"success": id > 0})
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Param id path int true "Album ID"
// @Success 200 {object} albumResponse
// @Failure 500 {object} object{message=string}
// @Router /album/{id} [get]
func getAlbumById(c *gin.Context) {
	db, err := sql.Open("mysql", "root:my-secret-pw@(localhost:3306)/GoTest?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	var a albumResponse

	id := c.Param("id")
	if err := db.QueryRow("SELECT * FROM album WHERE ID = ?", id).Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}

	c.IndentedJSON(http.StatusOK, a)

}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	seed()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		album := v1.Group("/album")
		{
			album.GET("", getAlbums)
			album.POST("", postAlbum)
			album.GET(":id", getAlbumById)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run("localhost:8080")
}
