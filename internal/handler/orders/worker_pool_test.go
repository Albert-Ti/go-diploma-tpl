package orders

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Albert-Ti/go-diploma-tpl/internal/models"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
	"github.com/stretchr/testify/require"
)

type testAccrualService struct {
	called atomic.Int32
}

func (t *testAccrualService) CheckAccrualOrder(ctx context.Context, task models.TaskOrder) error {
	t.called.Add(1)
	return nil
}
func TestWorkerPool(t *testing.T) {

	t.Run("Start workers", func(t *testing.T) {
		var (
			numWorkers = 3
			queueSize  = 3
		)

		mockSvc := &testAccrualService{}
		wp := NewWorkerPool(mockSvc, numWorkers, queueSize)
		defer wp.Stop()

		require.Eventually(t,
			func() bool {
				return wp.active.Load() == int32(numWorkers)
			},
			time.Second,
			10*time.Millisecond,
		)
	})

	t.Run("Queue overflow", func(t *testing.T) {
		var (
			numWorkers = 0
			queueSize  = 3
			totalTasks = 20
		)

		mockSvc := &testAccrualService{}
		wp := NewWorkerPool(mockSvc, numWorkers, queueSize)
		defer wp.Stop()

		for range totalTasks {
			hash, err := utils.RandomHash(5)
			require.NoError(t, err)

			wp.AddTask(hash, 1)
		}

		require.Equal(t, int32(totalTasks-queueSize), wp.dropped.Load())
		require.Equal(t, queueSize, len(wp.tasks))
	})

	t.Run("Worker pool stop", func(t *testing.T) {
		mockSvc := &testAccrualService{}

		wp := NewWorkerPool(mockSvc, 3, 10)

		done := make(chan struct{})

		go func() {
			wp.Stop()
			close(done)
		}()

		select {
		case <-done:
		// Stop() не зависает и завершается за 2 секунды
		case <-time.After(2 * time.Second):
			t.Fatal("worker pool did not stop")
		}
	})
}
