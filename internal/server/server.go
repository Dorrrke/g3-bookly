package server

import (
	"context"
	"errors"
	"time"

	"net/http"

	"github.com/Dorrrke/g3-bookly/internal/config"
	"github.com/Dorrrke/g3-bookly/internal/domain/models"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = "VerySecurKey2000Cat"

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

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

func (s *Server) Run() error {
	log := logger.Get()
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) { ctx.String(200, "Hello") })
	users := router.Group("/users")
	{
		users.GET("/info", s.userInfo)
		users.POST("/register", s.register)
		users.POST("/login", s.login)
	}
	books := router.Group("/books")
	{
		books.GET("/:id", s.bookInfo)
		books.GET("/:id/remove", s.removeBook)
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

func (s *Server) createJWTToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		},
		UserID: uid,
	})
	key := []byte(SecretKey)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func validToken(tokenStr string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", ErrInvalidToken
	}
	return claims.UserID, nil
}
