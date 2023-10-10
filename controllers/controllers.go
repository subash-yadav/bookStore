package controllers

import (
	"assignment/database"
	"assignment/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBookDetailsByISBN(c *gin.Context) {

	c.Header("Content-Type", "application.json")

	var bookDetails models.BookScore
	var recommendedBookList []models.BookScore

	params := c.Request.URL.Query()
	bookISBN := params.Get("bookISBN")

	fmt.Println("isbn number is ", bookISBN)

	if bookISBN == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide isbn number"})
		return
	}

	bookDetails, recommendedBookList, err := getBookByISBN(bookISBN)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"searched_book": bookDetails, "suggested_book": recommendedBookList})
}

func getBookByISBN(isbnNo string) (models.BookScore, []models.BookScore, error) {

	db := database.CreateConnection()

	defer db.Close()

	getBookDetailsQuery :=
		`select book.isbn13,book.title, book_language.language_name, book.num_pages,
		book.publication_date, category.category_name,book.no_of_visitors, json_agg(author.author_name) as author, publisher.publisher_name
		from book
		join book_category on book.book_id = book_category.book_id
		join category on category.category_id = book_category.category_id
		join book_language on book.language_id = book_language.language_id
		join publisher on book.publisher_id = publisher.publisher_id
		join book_author on book.book_id = book_author.book_id
		join author on book_author.author_id = author.author_id
		where book.isbn13 = $1
		group by book.isbn13,book.title, book_language.language_name, book.num_pages,
		book.publication_date, category.category_name, book.no_of_visitors, publisher.publisher_name`

	data, err := db.Query(getBookDetailsQuery, isbnNo)

	if err != nil {
		panic(err)
	}

	defer data.Close()
	var bookListStruct models.BookScore

	for data.Next() {
		var authorList []byte
		err := data.Scan(&bookListStruct.ISBN, &bookListStruct.Title, &bookListStruct.Language, &bookListStruct.NumberOfPages,
			&bookListStruct.PublicationDate, &bookListStruct.Category, &bookListStruct.NoOfVisitors, &authorList, &bookListStruct.Publication)
		if err != nil {
			panic(err)
		}
		publishDate, _ := time.Parse(time.RFC3339, bookListStruct.PublicationDate)

		bookListStruct.PublicationDate = publishDate.Format("02-01-2006")

		err = json.Unmarshal(authorList, &bookListStruct.Author)
		if err != nil {
			panic(err)
		}

	}
	UpdateNoOfVisitorsForBooks(isbnNo)

	recommendedBooks := GetRecommendedBook(isbnNo)

	return bookListStruct, recommendedBooks, err

}

func GetBookByAuthor(c *gin.Context) {
	c.Header("Content-Type", "application.json")

	params := c.Request.URL.Query()

	if params.Get("author_name") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Message": " Please provide filter criteria"})
		return
	}

	bookListByAuthor, err := getBookList(params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if bookListByAuthor == nil {
		c.JSON(http.StatusOK, gin.H{"Message": "No relavant book Found"})
		return
	}
	c.JSON(http.StatusOK, bookListByAuthor)

}

func getQuery(orderBy, sortType string) string {
	getBookListByAuthorQuery := fmt.Sprintf(`select  b.isbn13 as isbn, b.title, b.num_pages, l.language_name, p.publisher_name, b.publication_date
	from book b, publisher p,book_language l
	where book_id in( select book_id from book_author
		where author_id = ( select author_id from author
			where author_name = $1)) and b.publisher_id = p.publisher_id and
	b.language_id = l.language_id
	ORDER BY %s %s`, orderBy, sortType)

	return getBookListByAuthorQuery
}

func getBookList(params url.Values) ([]models.BookListByAuthor, error) {
	db := database.CreateConnection()

	defer db.Close()

	authorName := params.Get("author_name")
	sortBy := params.Get("sort_by")
	sortType := params.Get("sort_type")

	if sortBy == "" {
		sortBy = "publication_date"
	}
	if sortType == "" {
		sortType = "ASC"
	}

	getBookListByAuthorQuery := getQuery(sortBy, sortType)

	data, err := db.Query(getBookListByAuthorQuery, authorName)

	if err != nil {
		panic(err)
	}

	defer data.Close()

	var bookList []models.BookListByAuthor

	for data.Next() {
		var bookListStruct models.BookListByAuthor

		err = data.Scan(&bookListStruct.ISBN, &bookListStruct.Title, &bookListStruct.NumberOfPages, &bookListStruct.Language, &bookListStruct.Publication, &bookListStruct.PublicationDate)
		if err != nil {
			panic(err)
		}

		publishDate, _ := time.Parse(time.RFC3339, bookListStruct.PublicationDate)
		bookListStruct.PublicationDate = publishDate.Format("02-01-2006")

		bookList = append(bookList, bookListStruct)
	}
	return bookList, err
}

func PostAuthor(c *gin.Context) {
	c.Header("Content-Type", "application.json")

	var authorDetails models.AuthorDetails

	body, _ := io.ReadAll(c.Request.Body)
	err := json.Unmarshal([]byte(body), &authorDetails)

	if err != nil {
		panic(err)
	}
	if authorDetails.AuthorName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Please provide author name"})
		return
	}
	err = PostAuthorData(authorDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding author"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Author added successfully"})
}
