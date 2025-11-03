package orders_test

import (
	"fmt"
	"testing"

	"github.com/mrstecklo/micropet/services/orders/orders"
	"github.com/mrstecklo/micropet/services/orders/orders_mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderEngine_ForwardsCreateOrderToDatabase(t *testing.T) {
	data := []struct {
		id    int
		title string
	}{
		{
			1, "some title",
		},
		{
			1421, "duckling",
		},
	}
	for _, d := range data {
		t.Run(fmt.Sprint(d), func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			databaseMock := orders_mock.NewMockDatabase(mockCtrl)
			engine := orders.NewEngine(orders.Config{
				DB: databaseMock,
			})

			databaseMock.EXPECT().
				CreateOrder(d.title).
				Return(d.id, nil)

			id, err := engine.CreateOrder(d.title)

			assert.Equal(t, d.id, id)
			assert.Nil(t, err)
		})
	}
}
