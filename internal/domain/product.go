package domain

import (
	"errors"
	"strings"
	"sync"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func (p *Product) ValidateForCreate() error {
	if p.Name == "" {
		return errors.New("Нет имени")
	}
	if p.Price <= 0 {
		return errors.New("Нет стоимости")
	}
	if p.Quantity < 0 {
		return errors.New("Нет количества")
	}
	return nil
}

type Store struct {
	mu       sync.RWMutex
	products map[int]Product
	nextID   int
}

func NewStore() *Store {
	return &Store{
		products: make(map[int]Product),
		nextID:   1,
	}
}

func (s *Store) List(search string) []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if search == "" {

		temparr := make([]Product, 0, len(s.products))

		for _, v := range s.products {
			temparr = append(temparr, v)
		}

		return temparr
	}
	temparr := make([]Product, 0)

	for _, v := range s.products {
		if strings.Contains(v.Name, search) {
			temparr = append(temparr, v)
		}
	}

	return temparr

}

func (s *Store) Get(id int) (Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.products[id]
	return c, ok

}

func (s *Store) Create(p Product) Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	temp := Product{
		ID:       s.nextID,
		Name:     p.Name,
		Price:    p.Price,
		Quantity: p.Quantity,
	}
	s.products[temp.ID] = temp

	s.nextID++
	return temp
}

func (s *Store) Update(id int, p Product) (Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, ok := s.products[id]; ok == false {
		return c, ok
	}
	tempPr := s.products[id]
	if p.Name != "" {
		tempPr.Name = p.Name
	}
	if p.Price > 0 {
		tempPr.Price = p.Price
	}
	if p.Quantity >= 0 {
		tempPr.Quantity = p.Quantity
	}

	s.products[id] = tempPr

	return tempPr, true

}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.products[id]; ok == false {
		return false
	}
	delete(s.products, id)
	return true
}
