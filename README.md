# Async - Safe Goroutine Utilities for Go

Библиотека для безопасного запуска горутин с обработкой паник и удобными утилитами для асинхронного программирования в Go.

## Особенности

- **Безопасная обработка паник**: Все горутины автоматически защищены от паник
- **Контекстная поддержка**: Поддержка cancellation через context
- **Таймауты**: Встроенная поддержка таймаутов
- **Параллельное выполнение**: Утилиты для параллельного запуска нескольких задач
- **Retry механизм**: Повторные попытки с экспоненциальной задержкой
- **Пулы горутин**: Ограничение количества одновременно работающих горутин
- **Thread-safe**: Полная поддержка конкурентного использования

## Установка

```bash
go get github.com/WhiCu/async
```

## Быстрый старт

### Базовое использование

```go
package main

import (
    "fmt"
    "github.com/WhiCu/async"
)

func main() {
    // Безопасный запуск горутины
    g := async.SafeGo(func() {
        fmt.Println("Hello from goroutine!")
    })
    
    // Ожидание завершения
    err := g.Wait()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Обработка паник

```go
// Если функция паникует, паника будет преобразована в error
g := async.SafeGo(func() {
    panic("something went wrong!")
})

err := g.Wait()
if err != nil {
    fmt.Printf("Caught panic: %v\n", err)
    fmt.Printf("Original panic value: %v\n", g.Panic())
}
```

### Работа с возвращаемыми значениями

```go
// Получение результата через канал
result := <-async.Go(func() int {
    return 42
})

fmt.Printf("Value: %v, Error: %v\n", result.Value, result.Error)
```

### Таймауты и контекст

```go
// Таймаут
g := async.SafeGoWithTimeout(5*time.Second, func() {
    // долгая операция
})

err := g.Wait()
if err != nil {
    fmt.Printf("Timeout: %v\n", err)
}

// С контекстом
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

g := async.SafeGoWithContext(ctx, func() {
    // операция с возможностью отмены
})
```

### Параллельное выполнение

```go
// Запуск нескольких функций параллельно
errors := async.Parallel(
    func() { fmt.Println("Task 1") },
    func() { fmt.Println("Task 2") },
    func() { fmt.Println("Task 3") },
)

for i, err := range errors {
    if err != nil {
        fmt.Printf("Task %d failed: %v\n", i+1, err)
    }
}
```

### Пулы горутин

```go
// Создание пула с ограничением в 3 горутины
pool := async.NewPool(3)
defer pool.Close()

// Отправка задач в пул
for i := 0; i < 10; i++ {
    taskNum := i
    pool.Submit(func() {
        fmt.Printf("Task %d executing\n", taskNum)
        // обработка задачи
    })
}

// Ожидание завершения всех задач
pool.Wait()
```

### Retry с экспоненциальной задержкой

```go
value, err := async.RetryWithBackoff(3, time.Second, func() (int, error) {
    return someOperationThatMightFail()
})

if err != nil {
    fmt.Printf("Failed after all retries: %v\n", err)
}
```

## API Reference

### Goroutine API

- `SafeGo(f func()) *Goroutine` - Запуск функции в горутине с обработкой паник
- `SafeGoWithContext(ctx, f func()) *Goroutine` - Запуск с поддержкой контекста
- `SafeGoWithTimeout(timeout, f func()) *Goroutine` - Запуск с таймаутом

### Goroutine методы

- `Wait() error` - Ожидание завершения
- `WaitFor(timeout) error` - Ожидание с таймаутом
- `Cancel()` - Отмена выполнения
- `Panic() any` - Получение значения паники
- `HasPanicked() bool` - Проверка на панику
- `Done() <-chan struct{}` - Канал сигнализации завершения

### Value API

- `Go[T](f func() T) <-chan Result[T]` - Получение результата через канал
- `GoErr[T](f func() (T, error)) <-chan Result[T]` - Получение результата и ошибки
- `GoWithCallback[T](f func() T, callback)` - Вызов с callback

### Utility API

- `Parallel(...func()) []error` - Параллельное выполнение функций
- `ParallelWithResults[T](...func() T) []Result[T]` - Параллельное выполнение с результатами
- `RetryWithBackoff[T](attempts, backoff, f func() (T, error)) (T, error)` - Повтор с задержкой

### Pool API

- `NewPool(size int) *Pool` - Создание пула горутин
- `pool.Submit(f func())` - Отправка задачи в пул
- `pool.Wait()` - Ожидание завершения всех задач
- `pool.Close()` - Закрытие пула

## Архитектура

Библиотека состоит из двух основных пакетов:

- `panics` - Базовые утилиты для обработки паник
- `async` - Высокоуровневые функции для работы с горутинами

Библиотека построена на основе вашего существующего пакета `panics`, который обеспечивает thread-safe захват и обработку паник с использованием `atomic` операций.

## Производительность

Все функции оптимизированы для производительности:

- Минимальные аллокации памяти
- Использование батчинга для GC оптимизации
- Atomic операции для thread-safety без блокировок
- Эффективные каналы для communication

Запустите бенчмарки для проверки производительности:

```bash
go test -bench=. ./...
```

## Лицензия

MIT License
