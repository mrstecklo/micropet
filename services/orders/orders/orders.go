package orders

type Order struct {
	ID    int
	Title string
}

type Database interface {
	CreateOrder(title string) (int, error)
	GetOrder(id int) (Order, error)
}

type Engine struct {
	db Database
}

func (e Engine) CreateOrder(title string) (int, error) {
	return e.db.CreateOrder(title)
}

type Config struct {
	DB Database
}

func NewEngine(config Config) Engine {
	return Engine{
		db: config.DB,
	}
}
