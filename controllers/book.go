package controllers

import (
	"assignment/database"
	"assignment/models"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

func UpdateNoOfVisitorsForBooks(isbnNo string) {
	db := database.CreateConnection()
	defer db.Close()

	updateNoOfVisitorsQuery := `UPDATE book
	SET no_of_visitors = no_of_visitors+1
	WHERE isbn13= $1`

	data, err := db.Query(updateNoOfVisitorsQuery, isbnNo)

	if err != nil {
		panic(err)
	}
	fmt.Println("Number of visitors updated successfully for book with isbn no : ", isbnNo)

	defer data.Close()
}

func GetRecommendedBook(isbnNo string) []models.BookScore {
	db := database.CreateConnection()
	defer db.Close()

	getCurrentBookByISBN := `select category.category_name from category
	inner join book_category on category.category_id = book_category.category_id
	inner join book on book.book_id = book_category.book_id
	where book.isbn13 = $1`

	bookCategory, err := db.Query(getCurrentBookByISBN, isbnNo)
	if err != nil {
		panic(err)
	}

	var recommendedbookCategory models.BookCategory

	for bookCategory.Next() {
		err := bookCategory.Scan(&recommendedbookCategory.Category)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("recommendedbookCategory is", recommendedbookCategory.Category)

	getBookQuery :=
		// `select book.isbn13,book.no_of_visitors, category.category_name
		// from book
		// inner join book_category on book.book_id = book_category.book_id
		// inner join category on category.category_id = book_category.category_id`
		`select book.isbn13,book.title, book_language.language_name, book.num_pages,
		book.publication_date, category.category_name,book.no_of_visitors, json_agg(author.author_name) as author, publisher.publisher_name
		from book
		join book_category on book.book_id = book_category.book_id
		join category on category.category_id = book_category.category_id
		join book_language on book.language_id = book_language.language_id
		join publisher on book.publisher_id = publisher.publisher_id
		join book_author on book.book_id = book_author.book_id
		join author on book_author.author_id = author.author_id
		group by book.isbn13,book.title, book_language.language_name, book.num_pages,
		book.publication_date, category.category_name, book.no_of_visitors, publisher.publisher_name`

	data, err := db.Query(getBookQuery)
	if err != nil {
		panic(err)
	}

	var bookList []models.BookScore

	for data.Next() {
		var bookListStruct models.BookScore
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
		bookList = append(bookList, bookListStruct)
	}

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

	sort.Slice(bookList, func(i, j int) bool {
		return bookList[i].Score > bookList[j].Score
	})

	recommendedBook := bookList[:5]

	for _, value := range recommendedBook {
		fmt.Println(value)
	}
	return recommendedBook
}
