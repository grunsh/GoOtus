package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskCh := make(chan func() error)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	var workersWG sync.WaitGroup
	var retErr error

	// Запускаем n работяг
	for i := 0; i < n; i++ {
		workersWG.Add(1)
		go func() {
			defer workersWG.Done()
			for task := range taskCh {
				if err := task(); err != nil {
					errCh <- err
				}
			}
		}()
	}

	// Насяльника. Раздаёт всем задачки и считает ошибки. Завершает работу после завершения работы всех работяг.
	go func() {
		errCount := 0
	taskloop:
		for _, task := range tasks {
		chanRWLoop:
			for {
				select {
				case <-errCh:
					errCount++
					if errCount >= m {
						close(taskCh)
						break taskloop
					}
				case taskCh <- task:
					break chanRWLoop
				}
			}
		}

		// Если из цикла вышли НЕ по достижению лимита ошибок (кончились задачи)
		if !(errCount >= m) {
			close(taskCh)
		}

		// Остаётся только дождаться завершения работяг и досчитать ошибки, которые они вернут
		// Ожидание идёт в основной функции: workersWG.Wait() после которой получаем сообщение в doneCh
		for {
			select {
			case <-doneCh:
				if errCount >= m {
					retErr = ErrErrorsLimitExceeded
				} else {
					retErr = nil
				}
				return
			case <-errCh:
				errCount++
			}
		}
	}()

	// Ждём завершения всех работяг. Они либо доделают работу. Либо завершат работу по лимиту ошибок.
	workersWG.Wait()

	// Работяги закончили свою работу. Начальнику не от кого ждать ошибок. Увольняем.
	doneCh <- struct{}{}

	return retErr
}
