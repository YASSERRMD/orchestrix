package worker

import (
	"orchestrix/internal/repository"
	"orchestrix/internal/services"
	"log"
	"strconv"
	"time"
)

type Worker struct {
	workflowService *services.WorkflowService
	lockRepo        *repository.LockRepository
	stopCh          chan bool
}

func NewWorker() *Worker {
	return &Worker{
		workflowService: services.NewWorkflowService(),
		lockRepo:        repository.NewLockRepository(),
		stopCh:          make(chan bool),
	}
}

func (w *Worker) Start() {
	log.Println("Starting background worker...")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.checkTimeouts()
		case <-w.stopCh:
			log.Println("Worker stopped")
			return
		}
	}
}

func (w *Worker) Stop() {
	w.stopCh <- true
}

func (w *Worker) checkTimeouts() {
	w.lockRepo.Cleanup()

	instances, err := w.workflowService.GetTimedOutWorkflows(time.Now())
	if err != nil {
		log.Printf("Error getting timed out workflows: %v", err)
		return
	}

	for _, instance := range instances {
		lockKey := "workflow_timeout_" + strconv.FormatInt(instance.ID, 10)
		owner := "worker_" + time.Now().Format("20060102150405")

		acquired, err := w.lockRepo.Acquire(lockKey, owner, time.Now().Add(5*time.Minute))
		if err != nil {
			log.Printf("Error acquiring lock: %v", err)
			continue
		}

		if !acquired {
			continue
		}

		_, err = w.workflowService.Escalate(instance.ID)
		if err != nil {
			log.Printf("Error escalating workflow %d: %v", instance.ID, err)
		}

		w.lockRepo.Release(lockKey, owner)
	}
}
