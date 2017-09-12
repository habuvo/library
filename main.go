package library

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"encoding/json"
	"library/model"
	"context"
	"strconv"
	"fmt"
	"os/user"
)

type Env struct {
	BookStore model.BookDataStore
}

func main() {
	//формируем переменную окружения
	env := &Env{}
	//пилим роутер обрабатывающий подключения
	r := chi.NewRouter()
	//коннектимся к бд
	connection, err := model.ConnectToDB()
	if err != nil {
		log.Panic(err.Error())
	}
	// передаем в ORM чтоб она изменяла таблицу в бд в соответствии с структурой данных в модели
	connection.DB.AutoMigrate(
		&model.Book{},
	)

	// создаем и засовываем в переменную окружения обьект хранилища книг
	bookDataStore := model.BookDataStore{connection}
	env.BookStore = bookDataStore

	// вызываем прослойки для корректной работы сервера
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	//обьявляем пути по которым доступны данные
	r.Route("/book", func(r chi.Router) {
		r.Get("/", env.listAll)
		r.Post("/add", env.addBook)
		r.Route("/{bookID}", func(r chi.Router) {
			r.Use(env.BookCtx)
			r.Get("/", env.getBook)
			r.Post("/delete", env.deleteBook)
		})
	})
	//начинаем слушать порт 8090
	http.ListenAndServe(":8090", r)
}

//контекст для книги по ID
func (env *Env) BookCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		ibookID, err := strconv.Atoi(bookID)		//TODO Error handling
		book, err :=  env.BookStore.GetByID(ibookID)
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), "book", book)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//вывести книгу по ID
func (env *Env) getBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	book,ok := ctx.Value("book").(*model.Book)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	context, err := json.Marshal(book); if err == nil {
		w.Write(context)
	}
}


//добавить книгу   TODO fill fields
func (env *Env) addBook(w http.ResponseWriter, r *http.Request) {


/*	Author 			string `json:"bookauthor"`
	Publisher       string `json:"bookpublisher"`
	Pub_Year        int    `json:"pubyear"`
	Name	        string `json:"name"`
	Genre	        string `json:"genre"`
	Short_Desc

    r.Context() chi

*/

//парсим строку запроса
	author := r.FormValue("author")
	publisher := r.FormValue("publisher")


	if len(author) == 0 || len(publisher) == 0 { print ("oops") }


	book := model.Book {}
	book.Author = author


	err := env.BookStore.Add(&book)
	if err != nil {
		// произошла ошибка
		print("произошла ошибка")
		return
	}
	// если мы не вышли из выполнения по ошибке, идем дальше
	//w.Write([]byte("user write ok"))
	// или по другому -
	context, err := json.Marshal(book)
	w.Write(context)
	// в этом случае мы разберем на джсон конкретный обьект

}

//удалить книгу по ID
func (env *Env) deleteBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	book,ok := ctx.Value("book").(*model.Book)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	err := env.BookStore.Remove(book.ID)
	if err == nil {
		w.Write([]byte("удалена книга под номером "))
		w.Write([]byte(book.ID))
	}	else {
	http.Error(w, http.StatusText(422), 422)
	return
	}
}

