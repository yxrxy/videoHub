package service

import "github.com/yxrxy/videoHub/app/interaction/domain/repository"

type InteractionService struct {
	db repository.InteractionRepository
}

func NewInteractionService(db repository.InteractionRepository) *InteractionService {
	return &InteractionService{db: db}
}
