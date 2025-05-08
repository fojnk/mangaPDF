package db

import (
	"log/slog"
	"os"

	"github.com/fojnk/Task-Test-devBack/configs"
	"github.com/fojnk/Task-Test-devBack/pkg/tools"
	bolt "go.etcd.io/bbolt"
)

var DBconn *bolt.DB

func init() {
	var err error

	dbFileName := "grabber_data.db"

	checkDBFile := tools.IsFileExist(dbFileName)

	DBconn, err = bolt.Open(dbFileName, 0664, nil)
	if err != nil {
		slog.Error(
			"Ошибка при открытии файла БД",
			slog.String("Message", err.Error()),
		)
		os.Exit(1)
	}

	err = DBconn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("AppSettings"))
		return err
	})
	if err != nil {
		slog.Error(
			"Ошибка при создании бакета (AppSettings) в БД",
			slog.String("Message", err.Error()),
		)
		os.Exit(1)
	}

	err = DBconn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MangaFavs"))
		return err
	})
	if err != nil {
		slog.Error(
			"Ошибка при создании бакета (MangaFavs) в БД",
			slog.String("Message", err.Error()),
		)
		os.Exit(1)
	}

	err = DBconn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("History"))
		return err
	})
	if err != nil {
		slog.Error(
			"Ошибка при создании бакета (History) в БД",
			slog.String("Message", err.Error()),
		)
		os.Exit(1)
	}

	if !checkDBFile {
		err = DBconn.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AppSettings"))
			err := b.Put([]byte("dbver"), []byte(configs.DBver))
			return err
		})
		if err != nil {
			slog.Error(
				"Ошибка при создании стартовой БД",
				slog.String("Message", err.Error()),
			)
			os.Exit(1)
		}
	}
}
