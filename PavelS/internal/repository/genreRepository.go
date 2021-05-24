package repository

import (
	"fmt"
	"github.com/SavinskiPavel/bookstore/internal/database"
	"github.com/SavinskiPavel/bookstore/internal/model"
	"log"
)

func FindAllGenresByValue() []int {
	db := database.InitDB()
	rows, err := db.Query("SELECT * FROM genres")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var genres []int
	for rows.Next() {
		g := model.Genre{}
		err := rows.Scan(&g.Id, &g.Value, &g.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		genres = append(genres, g.Value)
	}
	return genres
}
