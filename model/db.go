package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// обьект хранящий в себе курсор для взаимодействия с бд
type Connection struct {
	*gorm.DB
}

// генерация строки подключения
// сигнатура функции - возвращает обьект бд, и ошибку. Если коннект прошел - ошибка = nil
func ConnectToDB() (*Connection, error) {

	db, err := gorm.Open("postgres", "host="+
		"127.0.0.1"+" user="+"book"+
		" dbname="+"book"+
		" sslmode=disable"+
		"password="+"book")
	if err != nil {
		print(err)
		return nil, err
	}
	print("connected!")
	//после подключения возвращаем ссылку на курсор работы с бд
	return &Connection{db}, nil
}
