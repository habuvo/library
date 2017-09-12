
#Book API

Simple API on Chi and GORM

methods:
+ **/api/book/{ID}**
get book by ID (HTTP GET), correct ID is three digits
*example:* http://myhost.me/api/book/001 

+ **/api/book/{ID}/delete**
delete book by ID (HTTP POST), correct ID is three digits
*example:* http://myhost.me/api/book/001/delete

+ **/api/book/**
show list of all books (HTTP GET)
*example:* http://myhost.me/api/book

+ **/api/book/add**
add book and show result (HTTP POST), required fields: author, name, year
*example:* http://myhost.me/api/add 