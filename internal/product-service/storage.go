package main

import (
	"sync"
)

type Storer interface {
	Upsert([]Product) error // Add or update by name
	GetAll() ([]Product, error)
	// Order([]Order) error
}

type MemoryStorer struct {
	mu      sync.RWMutex
	storage map[ProductName]Product
}

func NewMemoryStorer() *MemoryStorer {
	return &MemoryStorer{
		storage: make(map[ProductName]Product),
	}
}

func (ms *MemoryStorer) Upsert(products []Product) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, product := range products {
		if entry, ok := ms.storage[product.Name]; ok {
			entry.Articles = product.Articles
			entry.Price = product.Price
			ms.storage[entry.Name] = entry
		} else {
			ms.storage[product.Name] = Product{
				Name:     product.Name,
				Articles: product.Articles,
				Price:    product.Price,
			}
		}
	}

	return nil
}

func (ms *MemoryStorer) GetAll() ([]Product, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	products := make([]Product, 0, len(ms.storage))
	for _, value := range ms.storage {
		products = append(products, value)
	}

	return products, nil
}
