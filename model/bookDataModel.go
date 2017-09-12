package model

type Book struct {
	ID         int    `json:"bookid" gorm:"AUTO_INCREMENT"`
	Author     string `json:"bookauthor"`
	Publisher  string `json:"bookpublisher"`
	Pub_Year   int    `json:"pubyear"`
	Name       string `json:"name"`
	Genre      string `json:"genre"`
	Short_Desc string `json:"shortdesc"`
}

type BookDataStore struct {
	*Connection
}

//добавить книгу
func (d *BookDataStore) Add(up *Book) (error) {
	err := d.Connection.DB.Create(&up).Error
	if err != nil {
		return err
	}
	return nil
}

//удаление книги по ID
func (d *BookDataStore) Remove(bookID int) (error) {
	book := Book{ID: bookID}
	err := d.Connection.DB.Delete(&book).Error
	return err
}

//поиск книги по ID
func (d *BookDataStore) GetByID(bookID int) (*Book, error) {
	book := &Book{ID: bookID}
	err := d.Connection.DB.First(&book, bookID).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}

//список книг
func (d *BookDataStore) GetAll() (*[]Book, error) {
	books := &[]Book{}
	err := d.Connection.DB.Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}
