package orders

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
)

type WorkerPool struct {
	svc       service.AccrualChecker
	workers   int
	queueSize int
	tasks     chan models.TaskOrder
	done      chan struct{}
	wg        sync.WaitGroup
	active    atomic.Int32
	dropped   atomic.Int32

	// зашита от дублировании
	mu        sync.Mutex
	queueTask map[string]bool
}

func NewWorkerPool(svc service.AccrualChecker, numWorkers int, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		workers:   numWorkers,
		queueSize: queueSize,
		tasks:     make(chan models.TaskOrder, queueSize),
		done:      make(chan struct{}),
		svc:       svc,
		active:    atomic.Int32{},
		dropped:   atomic.Int32{},
		queueTask: make(map[string]bool),
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
	// atomic(атомарно) чтение + инкремент + запись происходят за один такт
	wp.active.Add(1)
	defer wp.active.Add(-1)

	for {
		select {
		case task := <-wp.tasks:
			slog.Info("Worker Pool", "worker", id, "started task", task)

			time.Sleep(time.Second)
			if err := wp.svc.CheckAccrualOrder(context.Background(), task); err != nil {
				slog.Error("Worker Pool", "worker", id, "error", err)
			}

			// Удаляем задачу из мапы после завершения
			wp.mu.Lock()
			delete(wp.queueTask, task.OrderID)
			wp.mu.Unlock()

			slog.Info("Worker Pool", "worker", id, "completed task", task)
		case <-wp.done:
			return
		}
	}
}

func (wp *WorkerPool) AddTask(orderID string, userID int) {
	wp.mu.Lock()

	if wp.queueTask[orderID] == true {
		wp.mu.Unlock()
		slog.Info("Duplicate task ignored", "order", orderID)
		return
	}
	wp.queueTask[orderID] = true
	wp.mu.Unlock()

	task := models.TaskOrder{
		OrderID: orderID,
		UserID:  userID,
	}
	select {
	case wp.tasks <- task:
		slog.Info("Added", "task", task.OrderID)
	default:
		wp.dropped.Add(1)
		slog.Info("Queue is full, task", task.OrderID, "not added")
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.done)

	wp.wg.Wait()
	close(wp.tasks)
	slog.Info("Worker Pool stopped")
}
