package jobertask

import (
	"context"
	"github.com/Jackalgit/BuildShortURL/internal/database"
	"github.com/google/uuid"
	"log"
)

const (
	numWorkers       = 3
	numBatchDataBase = 1000
)

var JobDict = make(map[uuid.UUID]*Job)

type Job struct {
	Ctx      context.Context
	JobID    uuid.UUID
	UserID   uuid.UUID
	TaskList *[]string
}

func NewJober(ctx context.Context, jobID uuid.UUID, userID uuid.UUID, taskList *[]string) *Job {
	return &Job{
		Ctx:      ctx,
		JobID:    jobID,
		UserID:   userID,
		TaskList: taskList,
	}
}

func (j *Job) DeleteURL() *Job {

	ctx, cancelFunc := context.WithCancel(j.Ctx)

	// сигнальный канал для завершения горутин
	doneCh := make(chan struct{})
	// закрываем его при завершении программы
	defer close(doneCh)
	// закрываем функцию отмены контекста
	defer cancelFunc()
	// запускаем принятую работу в отдельной горутине для возвращения в хендлер и и отдачи ответа клиенту
	go func() {
		// пишем входные данные в канал
		inputCh := Generator(doneCh, j.TaskList)
		// запускаем толпу рабочих дербанить канал
		fanOut(ctx, doneCh, j.UserID, inputCh)
	}()

	return &Job{}

}

func Generator(doneCh chan struct{}, input *[]string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, data := range *input {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

func fanOut(ctx context.Context, doneCh chan struct{}, userID uuid.UUID, inputCh chan string) {

	for i := 0; i < numWorkers; i++ {
		Worker(ctx, doneCh, userID, inputCh)
	}
}

func Worker(ctx context.Context, doneCh chan struct{}, userID uuid.UUID, inputCh chan string) {

	var deleteList []string

	go func() {

		for data := range inputCh {

			deleteList = append(deleteList, data)
			// пишем в базуданных при достижении определенного числа вставок
			if len(deleteList) == numBatchDataBase {
				err := database.DeleteURLUser(ctx, userID, deleteList)
				if err != nil {
					log.Println("[DeleteURLUser]", err)
					return
				}
			}

			_, ok := <-doneCh
			if !ok {
				return
			}
		}
		// дописываем остатки которые не вошли в "10"
		err := database.DeleteURLUser(ctx, userID, deleteList)
		if err != nil {
			log.Println("[DeleteURLUser]", err)
			return
		}
	}()
}
