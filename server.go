package main

//sqlite3 has to be imported for gorm to work
import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
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
	err = db.CreateTable(&todo{}).Error
	if err != nil {
		panic(err)
	}
	//Enable logging; prints queries
	db.LogMode(true)

	return func(c martini.Context) {
		c.Map(&db)
	}
}

func findTodo(db *gorm.DB, predicate *todo) (*todo, error) {
	t := new(todo)
	err := db.First(t, predicate).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func allTodosHandler(rendr render.Render, db *gorm.DB) {
	log.Println("Getting all todos")
	todos := make([]todo, 0, 0)
	db.Find(&todos)
	rendr.JSON(http.StatusOK, todos)
}

func oneTodoHandler(params martini.Params, rendr render.Render, db *gorm.DB, r *http.Request) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}

	t, err := findTodo(db, &todo{ID: id})
	if err == gorm.RecordNotFound {
		rendr.Text(http.StatusNotFound, "Resource not found")
		return
	}
	if err != nil {
		rendr.Text(http.StatusInternalServerError, err.Error())
		return
	}

	rendr.JSON(http.StatusOK, t)
}

func newTodoHandler(r *http.Request, rendr render.Render, db *gorm.DB) {
	todoText := r.FormValue("text")
	if todoText == "" {
		rendr.Text(http.StatusBadRequest, "Provide a `text` parameter")
		return
	}
	completedBool, err := strconv.ParseBool(r.FormValue("completed"))
	if err != nil {
		rendr.Text(http.StatusBadRequest, "Provide a `completed` paramter as a boolean")
		return
	}

	t := todo{Text: todoText, Completed: completedBool}
	db.Create(&t)

	rendr.JSON(http.StatusOK, t)
}

func deleteTodoHandler(params martini.Params, rendr render.Render, db *gorm.DB) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err)
	}

	t, err := findTodo(db, &todo{ID: id})
	if err != nil {
		rendr.Text(http.StatusBadRequest, err.Error())
		return
	}

	db.Delete(&t)
	rendr.JSON(http.StatusOK, t)
}

func updateTodoHandler(r *http.Request, params martini.Params, rendr render.Render, db *gorm.DB) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		rendr.Text(http.StatusBadRequest, "Provide a `completed` paramter as a boolean")
		return
	}

	t, err := findTodo(db, &todo{ID: id})
	if err != nil {
		rendr.Text(http.StatusBadRequest, err.Error())
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
}

func setupMartini() *martini.ClassicMartini {
	//Note that martini.classic gives you static as "public" by default
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(dbMiddleware())

	//Get index.html
	m.Get("/", func(r *http.Request, w http.ResponseWriter) {
		http.ServeFile(w, r, "index.html")
	})

	m.Group("/apiv0/todos", func(router martini.Router) {
		//Get all todos
		router.Get("/", allTodosHandler)
		//Get a specific todo based on id
		router.Get("/(?P<id>[0-9]+)", oneTodoHandler)
		//Add a new todo
		router.Post("/", newTodoHandler)
		//Delete todo with id
		router.Delete("/(?P<id>[0-9]+)", deleteTodoHandler)
		//Update todo with id
		router.Put("/(?P<id>[0-9]+)", updateTodoHandler)
	})

	return m
}

func main() {
	m := setupMartini()
	m.Run()
}
