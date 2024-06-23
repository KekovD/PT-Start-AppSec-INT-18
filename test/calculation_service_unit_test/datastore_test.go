package calculation_service_unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Add(key string, start, value int64) (int64, error) {
	args := m.Called(key, start, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) Get(key string, start int64) (int64, error) {
	args := m.Called(key, start)
	return args.Get(0).(int64), args.Error(1)
}

func TestDatastore(t *testing.T) {
	mockClient := new(MockRedisClient)
	datastore := mockClient

	key := "testKey"
	start := int64(1)
	value := int64(2)

	mockClient.On("Add", key, start, value).Return(int64(1), nil)
	res, err := datastore.Add(key, start, value)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res)

	mockClient.On("Get", key, start).Return(value, nil)
	res, err = datastore.Get(key, start)
	assert.NoError(t, err)
	assert.Equal(t, value, res)

	mockClient.AssertExpectations(t)
}
