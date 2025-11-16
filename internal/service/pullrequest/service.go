package pullrequest

// Service структура для бизнес-логики pull requests
type Service struct {
	repo Repo
}

// NewService создает новый Service
func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}
