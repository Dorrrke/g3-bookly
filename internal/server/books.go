package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/Dorrrke/g3-bookly/internal/domain/models"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	storerrros "github.com/Dorrrke/g3-bookly/internal/storage/errros"
	"github.com/gin-gonic/gin"
)

func (s *Server) allBooks(ctx *gin.Context) {
	log := logger.Get()
	_, exist := ctx.Get("uid")
	if !exist {
		log.Error().Msg("user ID not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	books, err := s.storage.GetBooks()
	if err != nil {
		if errors.Is(err, storerrros.ErrEmptyBooksList) {
			ctx.String(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (s *Server) bookInfo(ctx *gin.Context) {
	log := logger.Get()
	_, exist := ctx.Get("uid")
	if !exist {
		log.Error().Msg("user ID not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	id := ctx.Param("id")
	book, err := s.storage.GetBook(id)
	if err != nil {
		if errors.Is(err, storerrros.ErrBookNoExist) {
			ctx.String(http.StatusNotFound, err.Error())
			return
		}
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusFound, book)
}

func (s *Server) addBook(ctx *gin.Context) {
	log := logger.Get()
	_, exist := ctx.Get("uid")
	if !exist {
		log.Error().Msg("user ID not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	var book models.Book
	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		log.Error().Err(err).Msg("unmarshal body failed")
		ctx.String(http.StatusBadRequest, "incorrectly entered data")
		return
	}
	book.Count = 1
	if err := s.storage.SaveBook(book); err != nil {
		log.Error().Err(err).Msg("save user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, "book %s %s was added", book.Author, book.Lable)
}

func (s *Server) addBooks(ctx *gin.Context) {
	log := logger.Get()
	_, exist := ctx.Get("uid")
	if !exist {
		log.Error().Msg("user ID not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	var books []models.Book
	if err := ctx.ShouldBindBodyWithJSON(&books); err != nil {
		log.Error().Err(err).Msg("unmarshal body failed")
		ctx.String(http.StatusBadRequest, "incorrectly entered data")
		return
	}
	if err := s.storage.SaveBooks(books); err != nil {
		log.Error().Err(err).Msg("save user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, "%s books was added", len(books))
}

func (s *Server) removeBook(ctx *gin.Context) {
	log := logger.Get()
	_, exist := ctx.Get("uid")
	if !exist {
		log.Error().Msg("user ID not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	id := ctx.Param("id")
	if err := s.storage.SetDeleteStatus(id); err != nil {
		log.Error().Msg("user ID not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	s.delChan <- struct{}{}
	log.Debug().Int("chan len", len(s.delChan)).Msg("book send into delChan")
	ctx.String(http.StatusOK, "book "+id+" was deleted")
}

func (s *Server) deleter(ctx context.Context) {
	log := logger.Get()
	defer log.Debug().Msg("deleter was ended")
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("deleter context done")
			return
		default:
			if len(s.delChan) == cap(s.delChan) {
				log.Debug().Int("cap", cap(s.delChan)).Int("len", cap(s.delChan)).Msg("start deleting")
				for i := 0; i < cap(s.delChan); i++ {
					<-s.delChan
				}
				if err := s.storage.DeleteBooks(); err != nil {
					log.Error().Err(err).Msg("deleting books failed")
					s.ErrChan <- err
					return
				}
			}
		}
	}
}

func (s *Server) bookReturn(ctx *gin.Context) {}
