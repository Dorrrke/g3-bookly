package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/Dorrrke/g3-bookly/internal/domain/models"
	authgrpc "github.com/Dorrrke/g3-bookly/internal/grpc"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	storerrros "github.com/Dorrrke/g3-bookly/internal/storage/errros"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	resp, err := s.grpcServer.Register(context.Background(), &authgrpc.UserRegister{
		Login:    user.Email,
		Password: user.Pass,
		Age:      int32(user.Age),
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Msg("register user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.Header("Authorization", resp.GetToken())
	ctx.String(http.StatusCreated, resp.GetMessage())
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
	resp, err := s.grpcServer.Login(context.Background(), &authgrpc.User{
		Login:    user.Email,
		Password: user.Pass,
	})
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if status.Code(err) == codes.InvalidArgument {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Msg("login user failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.Header("Authorization", resp.GetToken())
	ctx.String(200, resp.GetMessage())
}

func (s *Server) userInfo(ctx *gin.Context) {
	log := logger.Get()
	uid := ctx.GetString("uid")
	user, err := s.storage.GetUser(uid)
	if err != nil {
		log.Error().Err(err).Msg("failed get user from db")
		if errors.Is(err, storerrros.ErrUserNotFound) {
			ctx.String(http.StatusNotFound, err.Error())
			return
		}
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusFound, user)
}
