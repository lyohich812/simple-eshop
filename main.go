package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	id          uint64
	productList []string
	status      uint8
	mu          sync.Mutex
}

const (
	created uint8 = iota
	paid
	delivered
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			o := Order{id: (uint64(i))}
			err := o.performOrder([]string{"bread", "eggs", "butter"})
			if err != nil {
				slog.Error("Error performing order", "error", err)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			nilOrder := Order{}
			err := nilOrder.performOrder([]string{"laptop", "smartphone"})
			if err != nil {
				slog.Error("Error performing nilorder", "error", err)
			}
		}()
	}
	wg.Wait()
}

func (o *Order) performOrder(productList []string) error {
	for _, product := range productList {
		if err := o.addProductToOrder(product); err != nil {
			return fmt.Errorf("Error adding product: %w", err)
		}
	}
	if err := o.performPayment(); err != nil {
		return fmt.Errorf("Error performing payment: %w", err)
	}
	if err := o.performDeliver(); err != nil {
		return fmt.Errorf("Error performing deliver: %w", err)
	}
	return nil
}

func (o *Order) addProductToOrder(product string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.status != created {
		return errors.New("invalid status")
	}

	if o.id == 0 {
		randomId := time.Now().Unix()*1000 + (rand.Int63n(100) + 1)
		o.id = uint64(randomId)
		slog.Info("Created order", "order_id", randomId)
	}

	o.productList = append(o.productList, product)
	slog.Info("Product added", "order_id", o.id, "product", product)
	return nil
}

func (o *Order) performPayment() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.status != created {
		return errors.New("invalid status")
	}

	slog.Info("Dumb payment performed", "order_id", o.id)
	o.status = paid

	return nil
}

func (o *Order) performDeliver() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.status != paid {
		return errors.New("invalid status")
	}

	slog.Info("Dumb deliver performed", "order_id", o.id)
	o.status = paid

	return nil
}

// Задача 2: Обработка заказов в интернет-магазине
// Описание задачи:
// Вы разрабатываете систему для обработки заказов в интернет-магазине.
// Каждый заказ состоит из нескольких шагов:
// добавление товаров в корзину,
// оплата,
// подтверждение доставки.
// Несколько пользователей могут одновременно добавлять товары в свои корзины,
// и необходимо гарантировать, чтобы данные о заказах не были повреждены.

// Требования:
// Реализуйте структуру Order, которая будет хранить информацию о заказе
// (например, ID заказа, список продуктов и статус заказа).
// Напишите функцию addProductToOrder, которая добавляет товар в корзину для указанного заказа.
// Если заказ ещё не создан, он должен быть создан с этим товаром.
// Используйте блокировку с помощью sync.Mutex
// для защиты доступа к данным о заказах, чтобы избежать гонки данных.
