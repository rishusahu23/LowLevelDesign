package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type TaskState string

const (
	Idle      TaskState = "Idle"
	Scheduled TaskState = "Scheduled"
	Completed TaskState = "Completed"
	Cancelled TaskState = "Cancelled"
)

type Task struct {
	Id          string
	UserId      string
	Name        string
	StartTime   time.Time
	Interval    time.Duration
	MaxRetries  int
	RetryCount  int
	Action      func()
	State       TaskState
	StopChannel chan struct{}
}

type TaskScheduler struct {
	tasks map[string]*Task
	mutex sync.Mutex
}

func generateId() string {
	times := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("error in generating random string")
	}
	randomPart := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%d-%s", times, randomPart)
}

func (s *TaskScheduler) AddTask(userId, name string, startTime time.Time, dur time.Duration, maxRetires int, action func()) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	taskId := generateId()
	task := &Task{
		Id:          taskId,
		UserId:      userId,
		Name:        name,
		StartTime:   startTime,
		Interval:    dur,
		MaxRetries:  maxRetires,
		RetryCount:  0,
		Action:      action,
		State:       Idle,
		StopChannel: make(chan struct{}),
	}

	s.tasks[taskId] = task
	return taskId
}

func (s *TaskScheduler) startTask(task *Task) {
	for {
		select {
		case <-task.StopChannel:
			task.State = Cancelled
			fmt.Println("task has been stopped")
			return
		default:
			if time.Now().After(task.StartTime) && task.State == Idle {
				if task.RetryCount >= task.MaxRetries {
					task.State = Completed
					fmt.Println("task has been completed")
					return
				}
				task.State = Scheduled
				go task.Action()
				task.RetryCount++
				fmt.Println("working on task")
				task.StartTime = time.Now().Add(task.Interval)
			}

		}
		time.Sleep(1 * time.Second)
	}
}

func (s *TaskScheduler) StopTask(taskId string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if task, exists := s.tasks[taskId]; exists {
		close(task.StopChannel)
		task.State = Cancelled
	}
}