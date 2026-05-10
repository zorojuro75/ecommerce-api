package usecase

import (
    "fmt"

    "ecommerce-api/internal/domain/entity"
    domainrepo "ecommerce-api/internal/domain/repository"
    "ecommerce-api/pkg/apperror"
    "ecommerce-api/pkg/hash"
    "ecommerce-api/pkg/jwt"
)

type userUsecase struct {
    repo      domainrepo.UserRepository
    jwtSecret string
}

func NewUserUsecase(repo domainrepo.UserRepository, jwtSecret string) entity.UserUsecase {
    return &userUsecase{repo: repo, jwtSecret: jwtSecret}
}

func (uc *userUsecase) Register(req entity.RegisterRequest) (*entity.User, error) {
    existing, _ := uc.repo.FindByEmail(req.Email)
    if existing != nil {
        return nil, fmt.Errorf("Register: %w", apperror.ErrConflict)
    }

    hashed, err := hash.BcryptHash(req.Password)
    if err != nil {
        return nil, fmt.Errorf("Register: hashing: %w", err)
    }

    u := &entity.User{
        Name:         req.Name,
        Email:        req.Email,
        PasswordHash: hashed,
        Role:         entity.RoleCustomer,
    }

    if err := u.Validate(); err != nil {
        return nil, fmt.Errorf("Register: %w", err)
    }
    if err := uc.repo.Create(u); err != nil {
        return nil, fmt.Errorf("Register: %w", err)
    }
    return u, nil
}

func (uc *userUsecase) Login(req entity.LoginRequest) (string, error) {
    u, err := uc.repo.FindByEmail(req.Email)
    if err != nil {
        return "", fmt.Errorf("Login: %w", apperror.ErrUnauthorized)
    }

    if !hash.BcryptCheck(req.Password, u.PasswordHash) {
        return "", fmt.Errorf("Login: %w", apperror.ErrUnauthorized)
    }

    token, err := jwt.Generate(u.ID, string(u.Role), uc.jwtSecret)
    if err != nil {
        return "", fmt.Errorf("Login: token: %w", err)
    }
    return token, nil
}

func (uc *userUsecase) GetUser(id uint) (*entity.User, error) {
    u, err := uc.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("GetUser id=%d: %w", id, err)
    }
    return u, nil
}