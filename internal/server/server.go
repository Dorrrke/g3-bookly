package server

import (
	"context"
	"net/http"

	"github.com/Dorrrke/g3-bookly/internal/config"
	"github.com/Dorrrke/g3-bookly/internal/domain/models"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Storage interface {
	SaveUser(models.User) (string, error)
	ValidUser(models.User) (string, error)
	SaveBook(models.Book) error
}

type Server struct {
	serv    *http.Server
	valid   *validator.Validate
	storage Storage
}

func New(cfg config.Config) *Server {
	server := http.Server{
		Addr: cfg.Addr,
	}
	valid := validator.New()
	return &Server{serv: &server, valid: valid}
}

func (s *Server) Run() error {
	log := logger.Get()
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) { ctx.String(200, "Hello") })
	users := router.Group("/users")
	{
		users.POST("/register", s.register)
		users.POST("/login", s.login)
		users.GET("/:id", s.userInfo)
	}
	books := router.Group("/books")
	{
		books.GET("/", s.allBooks)
		books.GET("/:id", s.bookInfo)
	}
	router.POST("/get-book", s.getBook)
	router.POST("/add-book", s.addBook)
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
		log.Error().Err(err).Msg("validate user failed")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.String(200, "user %s are logined", uuid)
}

func (s *Server) userInfo(ctx *gin.Context) {}

func (s *Server) allBooks(ctx *gin.Context) {}

func (s *Server) bookInfo(ctx *gin.Context) {}

func (s *Server) getBook(ctx *gin.Context) {}

func (s *Server) addBook(ctx *gin.Context) {
	log := logger.Get()
	var book models.Book
	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		log.Error().Err(err).Msg("unmarshal body failed")
		ctx.String(http.StatusBadRequest, "incorrectly entered data")
		return
	}
	if err := s.storage.SaveBook(book); err != nil {
		log.Error().Err(err).Msg("save user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, "book %s %s was added", book.Author, book.Lable)
}

func (s *Server) bookReturn(ctx *gin.Context) {}
