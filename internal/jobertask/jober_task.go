package jobertask

import (
	"github.com/Jackalgit/BuildShortURL/internal/database"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

const (
	numWorkers       = 3
	numBatchDataBase = 20
	timeDelete       = 3
)

var JobDict = make(map[uuid.UUID]*Job)

type Job struct {
	JobID    uuid.UUID
	UserID   uuid.UUID
	TaskList []string
}

func NewJober(jobID uuid.UUID, userID uuid.UUID, taskList []string) *Job {
	return &Job{
		JobID:    jobID,
		UserID:   userID,
		TaskList: taskList,
	}
}

// При старте сервиса создается общий канал inputChUserURL для записи пары УРЛ и Юзер отправленные пользователями
// на удаление.
// Можно настроить число numWorkers и numBatchDataBase для оптимизации скорости удаления.
// Воркер отправляет на удаление когда соберет numBatchDataBase если numBatchDataBase не удается собрать за
// время timeDelete, то запрос будет отправлен по прошествии этого времени.
// Так канал inputChUserURL открывается при запуске сервиса, то он по сути всегда открыт и время timeDelete
// позволяет избежать ситуации когда в ожидании numBatchDataBase зависнут "хвосты" на удаление.

func (j *Job) DeleteURL(inputChUserURL chan models.UserDeleteURL) *Job {

	// запускаем принятую работу в отдельной горутине для возвращения в хендлер и и отдачи ответа клиенту
	go func() {
		var wg sync.WaitGroup
		// сигнальный канал для завершения горутин
		doneCh := make(chan struct{})
		// закрываем его при завершении программы
		defer close(doneCh)
		// пишем входные данные в канал
		inputCh := Generator(doneCh, inputChUserURL, j.UserID, j.TaskList)
		// запускаем толпу рабочих дербанить канал
		fanOut(&wg, doneCh, inputCh)

		wg.Wait()

	}()

	return &Job{}

}

func Generator(doneCh chan struct{}, inputChUserURL chan models.UserDeleteURL, userID uuid.UUID, input []string) chan models.UserDeleteURL {

	go func() {

		for _, data := range input {
			select {
			case <-doneCh:
				return
			default:
				inputChUserURL <- models.UserDeleteURL{UserID: userID, ShortURL: data}
			}
		}
	}()
	return inputChUserURL
}

func fanOut(wg *sync.WaitGroup, doneCh chan struct{}, inputChUserURL chan models.UserDeleteURL) {

	for i := 0; i < numWorkers; i++ {
		Worker(wg, doneCh, inputChUserURL)
	}
}

func Worker(wg *sync.WaitGroup, doneCh chan struct{}, inputChUserURL chan models.UserDeleteURL) {

	var deleteList []models.UserDeleteURL

	ticker := time.NewTicker(timeDelete * time.Second)

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-doneCh:
				return
			case <-ticker.C:
				if len(deleteList) > 0 {
					err := database.DeleteURLUser(deleteList)
					if err != nil {
						log.Println("[DeleteURLUser]", err)
						return
					}
					deleteList = nil
				}
			}
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		for data := range inputChUserURL {

			select {
			case <-doneCh:
				return
			default:
				deleteList = append(deleteList, data)
				if len(deleteList) == numBatchDataBase {
					err := database.DeleteURLUser(deleteList)
					if err != nil {
						log.Println("[DeleteURLUser]", err)
						return
					}
					deleteList = nil
				}
			}
		}
		// дописываем остатки которые не вошли в numBatchDataBase, при закрытии канала inputChUserURL
		err := database.DeleteURLUser(deleteList)
		if err != nil {
			log.Println("[DeleteURLUser]", err)
			return
		}
	}()
}
