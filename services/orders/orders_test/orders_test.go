package orders_test

import (
	"errors"
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
				Database: databaseMock,
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

func TestOrderEngine_ReturnsDatabaseCreateOrderError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	databaseMock := orders_mock.NewMockDatabase(mockCtrl)
	engine := orders.NewEngine(orders.Config{
		Database: databaseMock,
	})
	expectedError := errors.New("oh, no!")
	databaseMock.EXPECT().
		CreateOrder(gomock.Any()).
		Return(0, expectedError)

	_, err := engine.CreateOrder("someting")

	assert.Equal(t, expectedError, err)
	assert.True(t, err == expectedError)
}
