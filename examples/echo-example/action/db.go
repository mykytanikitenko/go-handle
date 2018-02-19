package action

import "sync"

// Fake db
var DB = articlesDB{}

type Article struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type articlesDB struct {
	articles []Article
	sync.RWMutex
}

func (db *articlesDB) Insert(a Article) {
	db.Lock()
	defer db.Unlock()

	db.articles = append(db.articles, a)
}

func (db *articlesDB) All() []Article {
	db.RLock()
	defer db.RUnlock()

	return db.articles
}
