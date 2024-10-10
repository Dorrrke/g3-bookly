package storage

import (
	"errors"

	"github.com/Dorrrke/g3-bookly/internal/domain/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MemStorage struct {
	usersStor map[string]models.User
	bookStor  map[string]models.Book
}

func New() *MemStorage {
	return &MemStorage{
		usersStor: make(map[string]models.User),
		bookStor:  make(map[string]models.Book),
	}
}

func (ms *MemStorage) SaveUser(user models.User) (string, error) {
	//log := //logger.Get()
	uuid := uuid.New().String()
	if _, err := ms.findUser(user.Email); err == nil {
		return "", errors.New("user alredy exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	if err != nil {
		//log.Error().Err(err).Msg("save user failed")
		return "", err
	}
	//log.Debug().Str("hash", string(hash)).Send()
	user.Pass = string(hash)
	user.UID = uuid
	ms.usersStor[uuid] = user
	//log.Debug().Any("storage", ms.usersStor).Send()
	return uuid, nil
}
func (ms *MemStorage) ValidUser(user models.User) (string, error) {
	//log := //logger.Get()
	//log.Debug().Any("storage", ms.usersStor).Send()
	memUser, err := ms.findUser(user.Email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(memUser.Pass), []byte(user.Pass)); err != nil {
		return "", errors.New("invalid password")
	}
	return memUser.UID, nil
}

func (ms *MemStorage) SaveBook(book models.Book) error {
	memBook, err := ms.findBook(book)
	if err == nil {
		memBook.Count++
		ms.bookStor[memBook.BID] = memBook
		return nil
	}
	book.Count = 1
	bid := uuid.New().String()
	ms.bookStor[bid] = book
	return nil
}

func (ms *MemStorage) findUser(login string) (models.User, error) {
	for _, user := range ms.usersStor {
		if user.Email == login {
			return user, nil
		}
	}
	return models.User{}, errors.New("user does not exists")
}

func (ms *MemStorage) findBook(value models.Book) (models.Book, error) {
	for _, book := range ms.bookStor {
		if book.Lable == value.Lable && book.Author == value.Author {
			return book, nil
		}
	}
	return models.Book{}, errors.New("book does not exists")
}
