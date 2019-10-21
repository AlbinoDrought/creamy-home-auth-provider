package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type accountRepository interface {
	Add(username, password string) error
	Exists(username, password string) error
	Remove(username string) error
}

type fileBasedAccountRepository struct {
	directory string
	mapper    *mapper
}

func (repo *fileBasedAccountRepository) path(username string) (string, error) {
	hash, err := repo.mapper.Map(username)
	if err != nil {
		return "", err
	}

	return path.Join(repo.directory, hash), nil
}

func (repo *fileBasedAccountRepository) hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (repo *fileBasedAccountRepository) verify(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (repo *fileBasedAccountRepository) Add(username, password string) error {
	path, err := repo.path(username)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0400)
	if err != nil {
		return err
	}
	defer file.Close()

	hash, err := repo.hash(password)
	if err != nil {
		file.Close()
		os.Remove(path)
		return err
	}

	_, err = file.WriteString(hash)
	if err != nil {
		file.Close()
		os.Remove(path)
	}
	return err
}

func (repo *fileBasedAccountRepository) Exists(username, password string) error {
	path, err := repo.path(username)
	if err != nil {
		return err
	}

	fileBytes, fileErr := ioutil.ReadFile(path)
	var savedHash string
	if fileErr == nil {
		savedHash = string(fileBytes)
	} else {
		savedHash = ""
	}
	time.Sleep(time.Duration(rand.Int()%50) * time.Millisecond)

	err = repo.verify(password, savedHash)

	if fileErr != nil {
		return fileErr
	}

	return err
}

func (repo *fileBasedAccountRepository) Remove(username string) error {
	path, err := repo.path(username)
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func newFileBasedAccountRepository(directory string) (accountRepository, error) {
	if err := os.MkdirAll(directory, 0700); err != nil {
		return nil, err
	}
	return &fileBasedAccountRepository{
		directory: directory,
		mapper:    mustMakeHashMapper("sha2-256"),
	}, nil
}
