package server

import (
	"context"
	"errors"

	"net/http"

	"github.com/Dorrrke/g3-bookly/internal/config"
	"github.com/Dorrrke/g3-bookly/internal/domain/models"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	storerrros "github.com/Dorrrke/g3-bookly/internal/storage/errros"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Storage interface {
	SaveUser(models.User) (string, error)
	ValidUser(models.User) (string, error)
	SaveBook(models.Book) error
	SaveBooks([]models.Book) error
	GetUser(string) (models.User, error)
	GetBooks() ([]models.Book, error)
	GetBook(string) (models.Book, error)
}

type Server struct {
	serv    *http.Server
	valid   *validator.Validate
	storage Storage
}

func New(cfg config.Config, stor Storage) *Server {
	server := http.Server{
		Addr: cfg.Addr,
	}
	valid := validator.New()
	return &Server{serv: &server, valid: valid, storage: stor}
}

func (s *Server) ShutdownServer() error {
	return s.serv.Shutdown(context.Background())
}

func (s *Server) Run() error {
	log := logger.Get()
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) { ctx.String(200, "Hello") })
	users := router.Group("/users")
	{
		users.GET("/:id/info", s.userInfo)
		users.POST("/register", s.register)
		users.POST("/login", s.login)
	}
	books := router.Group("/books")
	{
		books.GET("/:id", s.bookInfo)
		books.GET("/", s.allBooks)
	}
	router.POST("/add-book", s.addBook)
	router.POST("/add-books", s.addBooks)
	router.POST("/book-return", s.bookReturn)

	s.serv.Handler = router

	log.Info().Str("host", s.serv.Addr).Msg("server started")
	if err := s.serv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() {
	s.serv.Shutdown(context.TODO())
}

func (s *Server) register(ctx *gin.Context) {
	log := logger.Get()
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Error().Err(err).Msg("unmarshal body failed")
		ctx.String(http.StatusBadRequest, "incorrectly entered data")
		return
	}
	if err := s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("validate user failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuid, err := s.storage.SaveUser(user)
	if err != nil {
		if errors.Is(err, storerrros.ErrUserExists) {
			log.Error().Msg(err.Error())
			ctx.String(http.StatusConflict, err.Error())
			return
		}
		log.Error().Err(err).Msg("save user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Str("uuid", uuid).Send()
	ctx.String(http.StatusCreated, uuid)
}

func (s *Server) login(ctx *gin.Context) {
	log := logger.Get()
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		log.Error().Err(err).Msg("unmarshal body failed")
		ctx.String(http.StatusBadRequest, "incorrectly entered data")
		return
	}
	if err := s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("validate user failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uuid, err := s.storage.ValidUser(user)
	if err != nil {
		if errors.Is(err, storerrros.ErrUserNoExist) {
			log.Error().Err(err).Msg("user not found")
			ctx.String(http.StatusNotFound, "invalid login or password: %w", err)
			return
		}
		if errors.Is(err, storerrros.ErrInvalidPassword) {
			log.Error().Err(err).Msg("invalid pass")
			ctx.String(http.StatusUnauthorized, err.Error())
			return
		}
		log.Error().Err(err).Msg("validate user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.String(200, "user %s are logined", uuid)
}

func (s *Server) userInfo(ctx *gin.Context) {
	// TODO: Должно быть возможно, только
	// TODO:       при наличии токена ( когда юзер вошел в систему )
	id := ctx.Param("id")
	user, err := s.storage.GetUser(id)
	if err != nil {
		if errors.Is(err, storerrros.ErrUserNotFound) {
			ctx.String(http.StatusNotFound, err.Error())
			return
		}
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusFound, user)
}

func (s *Server) allBooks(ctx *gin.Context) {
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
	// TODO: Должно быть возможно, только
	// TODO:       при наличии токена ( когда юзер вошел в систему )
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

func (s *Server) bookReturn(ctx *gin.Context) {}
