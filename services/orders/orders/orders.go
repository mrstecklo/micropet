package orders

type Order struct {
	ID    int
	Title string
}

type Database interface {
	CreateOrder(title string) (int, error)
	GetOrder(id int) (Order, error)
}

type MessagingSystem interface {
	PublishOrderCreated(order Order) error
}

type Engine struct {
	database  Database
	messaging MessagingSystem
}

func (e Engine) CreateOrder(title string) (int, error) {
	id, err := e.database.CreateOrder(title)
	e.messaging.PublishOrderCreated(Order{
		ID:    id,
		Title: title,
	})
	return id, err
}

type Config struct {
	Database  Database
	Messaging MessagingSystem
}

func NewEngine(config Config) Engine {
	return Engine{
		database:  config.Database,
		messaging: config.Messaging,
	}
}
