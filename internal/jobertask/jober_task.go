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

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

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

//func (j *Job) DeleteURL() *Job {
//
//	// запускаем принятую работу в отдельной горутине для возвращения в хендлер и и отдачи ответа клиенту
//	go func() {
//		var wg sync.WaitGroup
//		// сигнальный канал для завершения горутин
//		doneCh := make(chan struct{})
//		// закрываем его при завершении программы
//		defer close(doneCh)
//		// пишем входные данные в канал
//		inputCh := Generator(doneCh, j.TaskList)
//		// запускаем толпу рабочих дербанить канал
//		fanOut(&wg, doneCh, j.UserID, inputCh)
//
//		wg.Wait()
//
//	}()
//
//	return &Job{}
//
//}
//
//func Generator(doneCh chan struct{}, input []string) chan string {
//	inputCh := make(chan string)
//
//	go func() {
//		defer close(inputCh)
//
//		for _, data := range input {
//			select {
//			case <-doneCh:
//				return
//			case inputCh <- data:
//			}
//		}
//	}()
//	return inputCh
//}
//
//func fanOut(wg *sync.WaitGroup, doneCh chan struct{}, userID uuid.UUID, inputCh chan string) {
//
//	for i := 0; i < numWorkers; i++ {
//		Worker(wg, doneCh, userID, inputCh)
//	}
//}
//
//func Worker(wg *sync.WaitGroup, doneCh chan struct{}, userID uuid.UUID, inputCh chan string) {
//
//	var deleteList []string
//
//	wg.Add(1)
//
//	go func() {
//		defer wg.Done()
//
//		for data := range inputCh {
//
//			select {
//			case <-doneCh:
//				return
//			default:
//				deleteList = append(deleteList, data)
//				if len(deleteList) == numBatchDataBase {
//					err := database.DeleteURLUser(userID, deleteList)
//					if err != nil {
//						log.Println("[DeleteURLUser]", err)
//						return
//					}
//				}
//			}
//		}
//		// дописываем остатки которые не вошли в numBatchDataBase
//		err := database.DeleteURLUser(userID, deleteList)
//		if err != nil {
//			log.Println("[DeleteURLUser]", err)
//			return
//		}
//	}()
//}
