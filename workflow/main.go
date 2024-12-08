package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type WorkflowStatus string
type StepStatus string
type TaskType string

const (
	WorkflowStatusRunning   WorkflowStatus = "RUNNING"
	WorkflowStatusFailed    WorkflowStatus = "FAILED"
	WorkflowStatusCompleted WorkflowStatus = "COMPLETED"
	WorkflowStatusPending   WorkflowStatus = "PENDING"

	StepStatusRunning   StepStatus = "RUNNING"
	StepStatusFailed    StepStatus = "FAILED"
	StepStatusCompleted StepStatus = "COMPLETED"
	StepStatusPending   StepStatus = "PENDING"

	TaskTypeHttp   TaskType = "HTTP"
	TaskTypeLambda TaskType = "LAMBDA"
	TaskTypeDelay  TaskType = "DELAY"
)

type Workflow struct {
	ID             string
	Name           string
	WorkflowStatus WorkflowStatus
	Steps          []*Step
	CurrentStep    int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Step struct {
	ID         string
	Name       string
	TaskType   TaskType
	StepStatus StepStatus
	Parameters map[string]interface{}
	RetryCount int
	MaxRetries int
}

type StateStore struct {
	workflows map[string]*Workflow
	mu        sync.Mutex
}

func NewStateStore() *StateStore {
	return &StateStore{
		workflows: make(map[string]*Workflow),
	}
}

func (s *StateStore) SaveWorkflow(wf *Workflow) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workflows[wf.ID] = wf
}

type TaskExecutor struct {
}

func (ts *TaskExecutor) Execute(task *Step) error {
	switch task.TaskType {
	case TaskTypeHttp:
		fmt.Println("making api call")
		time.Sleep(1 * time.Second)
	case TaskTypeDelay:
		if delay, ok := task.Parameters["duration"].(int); ok {
			time.Sleep(time.Second * time.Duration(delay))
		}
	default:
		return errors.New("fat gaya")
	}
	return nil
}

type WorkflowManager struct {
	stateStore   *StateStore
	taskExecutor *TaskExecutor
}

func NewWorkflowManager(stateStore *StateStore, taskExecutor *TaskExecutor) *WorkflowManager {
	return &WorkflowManager{
		stateStore:   stateStore,
		taskExecutor: taskExecutor,
	}
}

func (w *WorkflowManager) StartWorkflow(wf *Workflow) error {
	wf.WorkflowStatus = WorkflowStatusRunning
	w.stateStore.SaveWorkflow(wf)

	for wf.CurrentStep < len(wf.Steps) {
		step := wf.Steps[wf.CurrentStep]
		step.StepStatus = StepStatusRunning
		err := w.taskExecutor.Execute(step)
		if err != nil {
			step.StepStatus = StepStatusFailed
			wf.WorkflowStatus = WorkflowStatusFailed
			w.stateStore.SaveWorkflow(wf)
			return err
		}
		step.StepStatus = StepStatusCompleted
		wf.CurrentStep++
		w.stateStore.SaveWorkflow(wf)
	}
	wf.WorkflowStatus = WorkflowStatusCompleted
	w.stateStore.SaveWorkflow(wf)
	return nil
}

func main() {
	store := NewStateStore()
	executor := &TaskExecutor{}

	manager := NewWorkflowManager(store, executor)

	wf := &Workflow{
		ID:             "id",
		Name:           "name",
		WorkflowStatus: WorkflowStatusPending,
		Steps: []*Step{
			{
				ID:         "1",
				Name:       "n",
				TaskType:   TaskTypeHttp,
				StepStatus: StepStatusPending,
				Parameters: map[string]interface{}{},
				RetryCount: 0,
				MaxRetries: 3,
			},
			{
				ID:         "2",
				Name:       "nn",
				TaskType:   TaskTypeDelay,
				StepStatus: StepStatusPending,
				Parameters: map[string]interface{}{
					"duration": 2,
				},
				RetryCount: 0,
				MaxRetries: 3,
			},
		},
		CurrentStep: 0,
	}
	store.SaveWorkflow(wf)
	manager.StartWorkflow(wf)
}