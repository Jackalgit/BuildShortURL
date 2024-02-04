package jobertask

import (
	"context"
	"github.com/Jackalgit/BuildShortURL/internal/database"
	"github.com/google/uuid"
	"log"
	"sync"
)

const (
	numWorkers       = 10
	numBatchDataBase = 100
)

var JobDict = make(map[uuid.UUID]*Job)

type Job struct {
	Ctx      context.Context
	JobID    uuid.UUID
	UserID   uuid.UUID
	TaskList []string
}

func NewJober(ctx context.Context, jobID uuid.UUID, userID uuid.UUID, taskList []string) *Job {
	return &Job{
		Ctx:      ctx,
		JobID:    jobID,
		UserID:   userID,
		TaskList: taskList,
	}
}

func (j *Job) DeleteURL() *Job {

	// запускаем принятую работу в отдельной горутине для возвращения в хендлер и и отдачи ответа клиенту
	go func() {
		var wg sync.WaitGroup
		// сигнальный канал для завершения горутин
		doneCh := make(chan struct{})
		// закрываем его при завершении программы
		defer close(doneCh)
		// пишем входные данные в канал
		inputCh := Generator(doneCh, j.TaskList)
		// запускаем толпу рабочих дербанить канал
		fanOut(j.Ctx, &wg, doneCh, j.UserID, inputCh)

		wg.Wait()

	}()

	return &Job{}

}

func Generator(doneCh chan struct{}, input []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()
	return inputCh
}

func fanOut(ctx context.Context, wg *sync.WaitGroup, doneCh chan struct{}, userID uuid.UUID, inputCh chan string) {

	for i := 0; i < numWorkers; i++ {
		Worker(ctx, wg, doneCh, userID, inputCh)
	}
}

func Worker(ctx context.Context, wg *sync.WaitGroup, doneCh chan struct{}, userID uuid.UUID, inputCh chan string) {

	var deleteList []string

	wg.Add(1)

	go func() {
		defer wg.Done()

		for data := range inputCh {

			select {
			case <-doneCh:
				return
			default:
				deleteList = append(deleteList, data)
				if len(deleteList) == numBatchDataBase {
					err := database.DeleteURLUser(ctx, userID, deleteList)
					if err != nil {
						log.Println("[DeleteURLUser]", err)
						return
					}
				}
			}
		}
		// дописываем остатки которые не вошли в numBatchDataBase
		err := database.DeleteURLUser(ctx, userID, deleteList)
		if err != nil {
			log.Println("[DeleteURLUser]", err)
			return
		}
	}()
}
