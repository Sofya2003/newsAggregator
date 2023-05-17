package storage

// Post Публикация, получаемая из RSS.
type Post struct {
	ID      int    // Номер записи
	Title   string // Заголовок публикации
	Content string // Содержание публикации
	PubTime int64  // Время публикации
	Link    string // Ссылка на источник
}

type Interface interface {
	GetPosts(n int) ([]Post, error)
	AddPost(p Post) error
}
