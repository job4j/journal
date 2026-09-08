package postgres

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository { return &AuthRepository{} }
