package internal

type UserUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uc *UserUsecase) GetUserByID(id int) (*User, error) {
	data, err := uc.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}
