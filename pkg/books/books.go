package books

type Book struct {
	Id    string
	Title string
}

type BookService interface {
	ListBooks() []Book
}
