package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	//sqlite3 has to be imported for gorm to work
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

func (t *todo) String() string {
	return fmt.Sprintf("%d - %s", t.ID, t.Text)
}

func dbMiddleware() martini.Handler {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalln(err.Error())
	}
	db.CreateTable(&todo{})

	return func(c martini.Context) {
		c.Map(&db)
	}
}

func setupMartini() *martini.ClassicMartini {
	//Note that martini.classic gives you static as "public" by default
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(dbMiddleware())

	m.Get("/", func(r *http.Request, w http.ResponseWriter) {
		http.ServeFile(w, r, "index.html")
	})

	m.Group("/todos", func(router martini.Router) {

		//Get all todos
		router.Get("/", func(rendr render.Render, db *gorm.DB) {
			log.Println("Getting all todos")
			todos := make([]todo, 0, 0)
			db.Find(&todos)
			rendr.JSON(http.StatusOK, todos)
		})

		//Get a specific todo based on id
		router.Get("/:id", func(params martini.Params, rendr render.Render, db *gorm.DB, r *http.Request) {
			id, err := strconv.Atoi(params["id"])
			if err != nil {
				log.Fatalln(err.Error())
			}

			var t todo
			db.Where(&todo{ID: id}).First(&t)

			if t.Text == "" {
				rendr.Text(http.StatusNotFound, "Resource not found")
			} else {
				rendr.JSON(http.StatusOK, t)
			}
		})

		//Add a new todo
		router.Post("/", func(r *http.Request, rendr render.Render, db *gorm.DB) {
			todoText := r.FormValue("text")
			completed := r.FormValue("completed")

			if len(todoText) == 0 || len(completed) == 0 {
				rendr.Text(http.StatusBadRequest, "'text' and 'completed' need to be provided")
				return
			}
			completedBool, err := strconv.ParseBool(completed)
			if err != nil {
				log.Fatalln(err.Error())
			}

			t := todo{Text: todoText, Completed: completedBool}
			db.Create(&t)
			rendr.JSON(http.StatusOK, t)
		})

		//Delete todo with id
		router.Delete("/:id", func(params martini.Params, rendr render.Render, db *gorm.DB) {
			id, err := strconv.Atoi(params["id"])
			if err != nil {
				log.Fatalln(err.Error())
			}
			var t todo
			db.Where(&todo{ID: id}).First(&t)
			if t.Text == "" {
				rendr.Text(http.StatusNotFound, "Resource not found")
				return
			}

			db.Delete(t)
			rendr.JSON(http.StatusOK, t)
		})

		//Update todo with id
		router.Put("/:id", func(r *http.Request, params martini.Params, rendr render.Render, db *gorm.DB) {
			id, err := strconv.Atoi(params["id"])
			if err != nil {
				log.Fatalln(err.Error())
			}

			var t todo
			db.Where(&todo{ID: id}).First(&t)
			if t.Text == "" {
				rendr.Text(http.StatusNotFound, "Resource not found")
				return
			}

			text := r.FormValue("text")
			if text != "" {
				t.Text = text
			}

			completed := r.FormValue("completed")
			if completed != "" {
				completedBool, err := strconv.ParseBool(completed)
				if err != nil {
					rendr.Text(http.StatusBadRequest, err.Error())
					return
				}
				t.Completed = completedBool
			}

			db.Save(t)
			rendr.JSON(http.StatusOK, t)
		})
	})

	return m
}

func main() {
	m := setupMartini()
	m.Run()
}
