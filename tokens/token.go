package tokens

import (
	"assignment/database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func ValidateToken(c *gin.Context) error {
	email, err := ExtractEmailFromToken(c)
	fmt.Println("Email is ", email)
	if err != nil {
		fmt.Println("Error is", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	var counter int = 0
	db := database.CreateConnection()

	defer db.Close()

	userQuery := `select count(*) from customer
	where email = $1`

	userData := db.QueryRow(userQuery, email)
	err = userData.Scan(&counter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return err
	}

	if counter == 0 {
		fmt.Println("This email not belongs to any customer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This email not belongs to any customer"})
		return nil
	}

	return nil

}

func ExtractEmailFromToken(c *gin.Context) (string, error) {
	err := godotenv.Load()
	// fmt.Println("I am here")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//retrieves token from header
	tokenString := c.Request.Header.Get("Authorization")
	fmt.Println("Token string is", tokenString)

	//parsing token string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println("expiry date is ", claims["exp"])
	fmt.Println("token validity is ", token.Valid)

	if ok && token.Valid {
		email := fmt.Sprintf("%v", claims["email"])
		fmt.Println("email extracted", email)
		return email, nil

	}
	return "", nil
}
