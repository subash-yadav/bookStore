package models

type BookISBN struct {
	Isbn string `json:"bookISBN"`
}

type BookDetails struct {
	Title           string   `json:"title"`
	Author          []string `json:"author"`
	Language        string   `json:"language"`
	NumberOfPages   int      `json:"number_of_pages"`
	PublicationDate string   `json:"publication_date"`
	Publication     string   `json:"publication"`
}

type Author struct {
	AuthorName string `json:"author_name"`
	Sortby     string `json:"sort_by"`
	Sorttype   string `json:"sort_type"`
}

type BookListByAuthor struct {
	ISBN            string `json:"book_isbn"`
	Title           string `json:"title"`
	NumberOfPages   int    `json:"number_of_pages"`
	Language        string `json:"language"`
	Publication     string `json:"publication"`
	PublicationDate string `json:"publication_date"`
}

type BookScore struct {
	ISBN            string   `json:"book_isbn"`
	Title           string   `json:"title"`
	Language        string   `json:"language"`
	NumberOfPages   int      `json:"number_of_pages"`
	PublicationDate string   `json:"publication_date"`
	Category        string   `json:"category"`
	Author          []string `json:"author"`
	Publication     string   `json:"publication"`
	NoOfVisitors    int      `json:"no_of_visitors"`
	Score           float64  `json:"book_score"`
}

type BookCategory struct {
	Category string `json:"book_category"`
}

type AuthorDetails struct {
	AuthorName string `json:"author_name"`
}
