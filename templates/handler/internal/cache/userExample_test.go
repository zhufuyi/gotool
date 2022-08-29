package cache

import (
	"context"
	"testing"
	"time"

	"github.com/zhufuyi/goctl/templates/handler/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/zhufuyi/pkg/mysql"
)

var (
	redisServer *miniredis.Miniredis
	redisClient *redis.Client
	testData    = &model.UserExample{
		Model: mysql.Model{ID: 1},
		Name:  "foo",
	}
	uc UserExampleCache
)

func setup() {
	redisServer = mockRedis()
	redisClient = redis.NewClient(&redis.Options{Addr: redisServer.Addr()})
	uc = NewUserExampleCache(redisClient)
}

func teardown() {
	redisServer.Close()
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	return s
}

func Test_userExampleCache_Set(t *testing.T) {
	setup()
	defer teardown()

	var id uint64
	ctx := context.Background()
	id = 1
	err := uc.Set(ctx, id, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleCache_Get(t *testing.T) {
	setup()
	defer teardown()

	var id uint64
	ctx := context.Background()
	id = 1
	err := uc.Set(ctx, id, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	act, err := uc.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testData, act)
}

func Test_userExampleCache_MultiGet(t *testing.T) {
	setup()
	defer teardown()

	ctx := context.Background()
	testData := []*model.UserExample{
		{Model: mysql.Model{ID: 1}},
		{Model: mysql.Model{ID: 2}},
	}
	err := uc.MultiSet(ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	expected := make(map[string]*model.UserExample)
	expected["1"] = &model.UserExample{Model: mysql.Model{ID: 1}}
	expected["2"] = &model.UserExample{Model: mysql.Model{ID: 2}}

	act, err := uc.MultiGet(ctx, []uint64{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected, act)
}

func Test_userExampleCache_MultiSet(t *testing.T) {
	setup()
	defer teardown()

	ctx := context.Background()
	testData := []*model.UserExample{
		{Model: mysql.Model{ID: 1}},
		{Model: mysql.Model{ID: 2}},
	}
	err := uc.MultiSet(ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userExampleCache_DelUser(t *testing.T) {
	setup()
	defer teardown()

	var id uint64
	ctx := context.Background()
	id = 1
	err := uc.Del(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
}
