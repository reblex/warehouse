package main

import (
	"errors"
	"sync"
)

type Storer interface {
	Store([]Article) error // Add sum of count to keys
	GetAll() ([]Article, error)
	// Reserve([]Reservation) error
}

type MemoryStorer struct {
	mu      sync.RWMutex
	storage map[ArticleId]Article
}

func NewMemoryStorer() *MemoryStorer {
	return &MemoryStorer{
		storage: make(map[ArticleId]Article),
	}
}

func (ms *MemoryStorer) Store(articles []Article) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, article := range articles {
		if article.Stock < 0 {
			return errors.New("article stock addition cannot be negative")
		}

		if entry, ok := ms.storage[article.Id]; ok {
			entry.Stock += article.Stock
			ms.storage[entry.Id] = entry
		} else {
			ms.storage[article.Id] = Article{
				Id:    article.Id,
				Name:  article.Name,
				Stock: article.Stock,
			}
		}
	}

	return nil
}

func (ms *MemoryStorer) GetAll() ([]Article, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	articles := make([]Article, 0, len(ms.storage))
	for _, value := range ms.storage {
		articles = append(articles, value)
	}

	return articles, nil
}
