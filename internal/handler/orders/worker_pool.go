package orders

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
)

type WorkerPool struct {
	svc       *service.Service
	workers   int
	queueSize int
	tasks     chan models.TaskOrder
	done      chan struct{}
	wg        sync.WaitGroup
}

func NewWorkerPool(svc *service.Service, numWorkers int, queueSize int) *WorkerPool {
	wp := &WorkerPool{
		workers:   numWorkers,
		queueSize: queueSize,
		tasks:     make(chan models.TaskOrder, queueSize),
		svc:       svc,
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
			slog.Info("Worker Pool", "worker", id, "started task", task)

			time.Sleep(time.Second)
			if err := wp.svc.CheckAccrualOrder(context.Background(), task); err != nil {
				slog.Error("Worker Pool", "worker", id, "error", err)
			}

			slog.Info("Worker Pool", "worker", id, "completed task", task)
		case <-wp.done:
			return
		}
	}
}

func (wp *WorkerPool) AddTask(orderID string, userID int) {
	task := models.TaskOrder{
		OrderID: orderID,
		UserID:  userID,
	}
	select {

	case wp.tasks <- task:
		slog.Info("Added", "task", task.OrderID)
	default:
		slog.Info("Queue is full, task", task.OrderID, "not added")
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.done)

	wp.wg.Wait()
	close(wp.tasks)
	slog.Info("Worker Pool stopped")
}
