package service

import (
	"encoding/json"
	"errors"
	"github.com/DenrianWeiss/catballCoin/model"
	"github.com/dgraph-io/badger/v3"
)

const (
	taskPrefix = "task_"
)

var (
	database *badger.DB
)

func initDatabase() {
	db, err := badger.Open(badger.DefaultOptions(GlobalConfig.DBPath))
	if err != nil {
		panic("Cannot open database")
	}
	database = db
}

func GetTaskList() []model.CoinTask {
	tasks := make([]model.CoinTask, 0)
	_ = database.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek([]byte(taskPrefix)); it.ValidForPrefix([]byte(taskPrefix)); it.Next() {
			item := it.Item()
			_ = item.Value(func(val []byte) error {
				newTask := model.CoinTask{}
				err := json.Unmarshal(val, &newTask)
				if err != nil {
					tasks = append(tasks, newTask)
				}
				return nil
			})
		}
		return nil
	})
	return tasks
}

func AddTaskToDatabase(task *model.CoinTask) error {
	return database.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(taskPrefix + task.SourceKey[:16]))
		if err == badger.ErrKeyNotFound {
			j, err := json.Marshal(task)
			if err != nil {
				return err
			}
			return txn.Set([]byte(taskPrefix + task.SourceKey[:16]), j)
		}
		return errors.New("already exists")
	})
}
