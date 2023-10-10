# 1. Online Book Store

This is a simple online book store backend design.

# 2. Functionality
    - It can get book details using ISBN number
    - It can get Book details based on Auther name and extensible to other filteration criteria.
    - It is able to sort the data based on user given scenario
    - Getting book using ISBN number will also give a suggestion list of five books with all book details
    - User Authentication also implemented based jwt token
    - Authorised user can add details

# Useful commands
    - go run main.go // to start server
    - go get <package name> // to load specific package
    - clear // clear the terminal

# ER Diagram 
final db design : 
https://dbdiagram.io/d/6524a563ffbf5169f05b4afa

# JWT 
JWT payload structure
```ruby
{
  "email" : "abc@gmail.com"
}
```

# Formulae used to calculate recommendation score
```ruby
    inCategory := 1.0
	other := 0.5
	for idx, value := range bookList {
		multiplier := 0.0
		if recommendedbookCategory.Category == value.Category {
			multiplier = inCategory
		} else {
			multiplier = other  
		}
		bookList[idx].Score = float64(value.NoOfVisitors) * multiplier
	}
```

# API ENDPOINT with response screenshot of postman

    1. Get book details using isbn number
        . EndPoint : http://localhost:8080/books/getBookByISBN?bookISBN=785342314526
        . Method : GET

![AScreenshot1.png](/screenShots/Screenshot1.png)

    2. Get book details using author name
        . EndPoint : http://localhost:8080/books/getBookByAuthor?author_name=Plato&sort_by=publication_date&sort_type=DESC
        . Method : GET
![screenshot2.png](/screenShots/screenshot2.png)

    3. Post author 
        . EndPoint : http://localhost:8080/user/postAuthorDetails
        . Method : POST
        . Payload : {
                        "author_name" : "abc_test"
                    }
        . Header : Authorization
![screenshot3.png](/screenShots/screenshot3.png)
