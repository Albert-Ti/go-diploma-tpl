package orders

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type WorkerPool struct {
	workers   int
	queueSize int
	tasks     chan string
	results   chan int
	done      chan struct{}
	wg        sync.WaitGroup
}

func NewWorkerPool(numWorkers int, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		workers:   numWorkers,
		queueSize: queueSize,
		tasks:     make(chan string, queueSize),
		results:   make(chan int, queueSize),
	}

	go func() {
		for i := 0; i < wp.workers; i++ {
			wp.wg.Add(1)
			go wp.worker(i)
		}
	}()

	return wp
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case task := <-wp.tasks:
			fmt.Println("рабочий", id, "запущена задача", task)

			time.Sleep(time.Second)

			wp.results <- checkAccrualOrder(task)

			fmt.Println("рабочий", id, "закончил задача", task)
		case <-wp.done:
			return
		}
	}
}

func (wp *WorkerPool) AddTask(task string) {
	select {
	case wp.tasks <- task:
		slog.Info("Список", "задач", len(wp.tasks))
	default:
		slog.Info("очередь заполнена, задача", task, "не добавлена")
	}
}

func (wp *WorkerPool) Stop() {
	wp.wg.Wait()      // ждём завершения всех воркеров
	close(wp.done)    // сигнал воркерам завершиться
	close(wp.tasks)   // закрываем канал задач
	close(wp.results) // закрываем канал результатов
	slog.Info("Worker Pool остановлен")
}
