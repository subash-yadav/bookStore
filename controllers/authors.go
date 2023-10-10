package controllers

import (
	"assignment/database"
	"assignment/models"
)

func PostAuthorData(authorDetail models.AuthorDetails) error {
	db := database.CreateConnection()

	defer db.Close()

	var authorId int

	getAuthorIdQuery := `select max(author_id) from author`

	id := db.QueryRow(getAuthorIdQuery)
	err := id.Scan(&authorId)
	if err != nil {
		panic(err)
	}

	postAuthorDetailsQuery := `insert into author
	values($1,$2)`

	_, err = db.Exec(postAuthorDetailsQuery, authorId+1, authorDetail.AuthorName)

	return err
}
