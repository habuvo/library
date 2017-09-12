package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"encoding/json"
	"library/model"
	"context"
	"strconv"
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
	r.Route("/api/book", func(r chi.Router) {
		r.Get("/", env.listAll)
		r.Post("/add", env.addBook)
		r.Get("/{bookID:[0-9]{3}}", env.showBookByID)
		r.Post("/delete", env.deleteBook)
	})
	//начинаем слушать порт 8090
	http.ListenAndServe(":8090", r)
}

//добавить книгу
func (env *Env) addBook(w http.ResponseWriter, r *http.Request) {

	//парсим запрос

	author := r.FormValue("author")
	publisher := r.FormValue("publisher")
	year := r.FormValue("pubyear")
	name := r.FormValue("name")
	genre := r.FormValue("genre")
	short_desc := r.FormValue("short_desc")

	if len(author) == 0 || len(name) == 0 {
		w.Write([]byte("Не заполнены обязательные поля"))
		return
	}

	book := model.Book{}
	book.Author = author
	book.Genre = genre
	book.Name = name
	pyear, err := strconv.Atoi(year);
	if err != nil {
		w.Write([]byte("Год издания не число"))
		return
	}
	book.Pub_Year = pyear
	book.Publisher = publisher
	book.Short_Desc = short_desc

	//write entity
	err = env.BookStore.Add(&book)
	if err != nil {
		// произошла ошибка
		print("произошла ошибка")
		return
	}

	context, err := json.Marshal(book)
	w.Write(context)

}

//удалить книгу
func (env *Env) deleteBook(w http.ResponseWriter, r *http.Request) {

	bookID := r.FormValue("ID");
	if bookID == "" {
		w.Write([]byte("Не задан ID"))
		return
	}

	ibookID, err := strconv.Atoi(bookID);
	if err != nil {
		w.Write([]byte("ошибка приведения"))
		return
	}
	//проверяем есть ли такая книга
	book, er := env.BookStore.GetByID(ibookID);
	if book == nil {
		if er != nil {
			w.Write([]byte("Нет книги с номером " + bookID))
			return
		} else {
			w.Write([]byte("ошибка БД"))
		}
	}

	er = env.BookStore.Remove(ibookID);
	if er == nil {
		message := "удалена книга под номером " + bookID
		w.Write([]byte(message))
		return
	} else {
		http.Error(w, http.StatusText(422), 422)
	}
}

//показать книгу по ID
func (env *Env) showBookByID(w http.ResponseWriter, r *http.Request) {

	book := &model.Book{}

	if bookID := chi.URLParam(r, "bookID"); bookID != "" {
		ibookID, err := strconv.Atoi(bookID);
		if err != nil {
			w.Write([]byte("ошибка приведения"))
			return
		}
		book, err = env.BookStore.GetByID(ibookID);
		if err != nil {
			http.Error(w, http.StatusText(422), 422)
			return
		}
	} else {
		w.Write([]byte("не задан ID"))
		return
	}

	context, err := json.Marshal(book);
	if err != nil {
		w.Write([]byte("Ошибка маршаллинга"))
		return
	}
	w.Write(context)

}

//показать все книги
func (env *Env) listAll(w http.ResponseWriter, r *http.Request) {
	//books is []model.Book{}
	books, err := env.BookStore.GetAll();
	if err != nil {
		w.Write([]byte("Ошибка получения списка"))
		return
	}

	for _, value := range *books {
		context, err := json.Marshal(value);
		if err != nil {
			w.Write([]byte("Ошибка маршаллинга"))
			break
		}
		w.Write(context)
	}
}

/* 		---- context mode ---
		r.Route("/{bookID:[0-9]{3}}", func(r chi.Router) {
			r.Use(env.BookCtx)
			r.Get("/", env.showBookByID)
			r.Post("/delete", env.deleteBook)
		})
*/

//контекст для книги по ID
func (env *Env) BookCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")
		ibookID, err := strconv.Atoi(bookID)
		if err != nil {
			w.Write([]byte("Ошибка приведения к инт"))
			return
		}
		book, err := env.BookStore.GetByID(ibookID)
		if err != nil {
			w.Write([]byte("Нет такой книги"))
			return
		}

		ctx := context.WithValue(r.Context(), "book", book)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//вывести книгу по ID (с контекстом)
func (env *Env) getBookCtx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	book, ok := ctx.Value("book").(*model.Book)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	context, err := json.Marshal(book);
	if err == nil {
		w.Write(context)
	}
}

//удалить книгу по ID (с контекстом)
func (env *Env) deleteBookCtx(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	book, ok := ctx.Value("book").(*model.Book)

	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	err := env.BookStore.Remove(book.ID);
	if err == nil {
		message := "удалена книга под номером " + strconv.Itoa(book.ID)
		w.Write([]byte(message))
		return
	} else {
		http.Error(w, http.StatusText(422), 422)
	}
}
