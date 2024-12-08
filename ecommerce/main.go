package main

type User struct {
	Id       string
	Name     string
	Email    string
	Password string
	Address  string
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	Update(user *User) error
}

type UserService interface {
	Register(user *User) error
	Login(email, password string) (*User, error)
	UpdateProfile(user *User) error
}

type UserServiceImpl struct {
	repo UserRepository
}

func (u *UserServiceImpl) Register(user *User) error {
	return u.repo.Create(user)
}

func (u *UserServiceImpl) Login(email, password string) (*User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	// Add password validation logic here
	return user, nil
}

func (u *UserServiceImpl) UpdateProfile(user *User) error {
	return u.repo.Update(user)
}

type Product struct {
	ID       string
	Name     string
	Category string
	Brand    string
	SKU      string
	Stock    int
	Price    float64
}

type ProductRepository interface {
	GetByID(id string) (*Product, error)
	Search(query string) ([]Product, error)
	UpdateStock(productID string, newStock int) error
}
