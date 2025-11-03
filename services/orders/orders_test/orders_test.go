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

type ordersEngineFixture struct {
	engine        orders.Engine
	mockCtrl      *gomock.Controller
	databaseMock  *orders_mock.MockDatabase
	messagingMock *orders_mock.MockMessagingSystem
}

func setUpOrdersEngineTest(t *testing.T) ordersEngineFixture {
	mockCtrl := gomock.NewController(t, gomock.WithOverridableExpectations())
	databaseMock := orders_mock.NewMockDatabase(mockCtrl)
	messagingMock := orders_mock.NewMockMessagingSystem(mockCtrl)
	messagingMock.EXPECT().
		PublishOrderCreated(gomock.Any()).
		Return(nil).
		AnyTimes()
	engine := orders.NewEngine(orders.Config{
		Database:  databaseMock,
		Messaging: messagingMock,
	})
	return ordersEngineFixture{
		engine:        engine,
		mockCtrl:      mockCtrl,
		databaseMock:  databaseMock,
		messagingMock: messagingMock,
	}
}

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
			f := setUpOrdersEngineTest(t)

			f.databaseMock.EXPECT().
				CreateOrder(d.title).
				Return(d.id, nil)

			id, err := f.engine.CreateOrder(d.title)

			assert.Equal(t, d.id, id)
			assert.Nil(t, err)
		})
	}
}

func TestOrderEngine_ReturnsDatabaseCreateOrderError(t *testing.T) {
	f := setUpOrdersEngineTest(t)
	expectedError := errors.New("oh, no!")
	f.databaseMock.EXPECT().
		CreateOrder(gomock.Any()).
		Return(0, expectedError)

	_, err := f.engine.CreateOrder("someting")

	assert.Equal(t, expectedError, err)
	assert.True(t, err == expectedError)
}

func TestOrderEngine_PublishesCreatedOrder(t *testing.T) {
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
			f := setUpOrdersEngineTest(t)
			f.databaseMock.EXPECT().
				CreateOrder(gomock.Any()).
				Return(d.id, nil)

			f.messagingMock.EXPECT().
				PublishOrderCreated(orders.Order{
					ID:    d.id,
					Title: d.title,
				}).
				Return(nil)

			_, _ = f.engine.CreateOrder(d.title)
		})
	}
}
