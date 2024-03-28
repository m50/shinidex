package database

import (
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	id := generateId()
	assert.Equal(t, 32, len(id))
}

func TestGenerateIdSequential(t *testing.T) {
	ids := make([]string, 5)
	for i := 0; i < 5; i++ {
		ids[i] = generateId()
	}
	assert.True(t, sort.StringsAreSorted(ids), "Strings are not sorted")
}

func setupDB(t *testing.T) *Database {
	d, _ := os.Getwd()
	db, err := NewLocal(":memory:")
	assert.Nil(t, err, "There is an error creating in memory database ", err)
	err = db.Migrate(d + "/../../migrations")
	assert.Nil(t, err, "There is an error migrating ", err)

	return db
}

func TestMigrate(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	c := 0
	err := db.conn.Get(&c, "SELECT count(*) FROM migrations")
	assert.Nil(t, err, "There is an error counting migrations ", err)
	assert.Greater(t, c, 0)
}

func TestUser(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	err := db.Users().Insert(User{
		Email:    "test@test.com",
		Password: "my-password",
	})
	assert.Nil(t, err, "Unable to insert user")
	user, err := db.Users().FindByEmail("test@test.com")
	assert.Nil(t, err, "Unable to get user")
	fmt.Println(user)
	assert.Equal(t, "my-password", user.Password)
	user.Password = "test"
	err = db.Users().Update(user)
	assert.Nil(t, err, "Unable to update user")
	user, err = db.Users().FindByID(user.ID)
	assert.Nil(t, err, "Unable to get user")
	assert.Equal(t, "test", user.Password)
	err = db.Users().Delete(user.ID)
	assert.Nil(t, err, "Unable to delete user")
	_, err = db.Users().FindByID(user.ID)
	assert.Error(t, err)
}
