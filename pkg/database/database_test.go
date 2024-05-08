package database

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/m50/shinidex/pkg/types"
	"github.com/stretchr/testify/assert"
	l "github.com/m50/shinidex/pkg/logger"
)

func SetupDB(t *testing.T) *Database {
	t.Helper()
	d, _ := os.Getwd()
	db, err := NewLocal(":memory:")
	assert.Nil(t, err, "There is an error creating in memory database ", err)
	err = db.Migrate(d + "/../../migrations")
	assert.Nil(t, err, "There is an error migrating ", err)

	return db
}

func SetupDBWithLogger(t *testing.T, logger *log.Logger) *Database {
	t.Helper()
	l.SetDefaultLogger(logger)
	return SetupDB(t)
}

func TestGenerateId(t *testing.T) {
	t.Parallel()
	id := generateId()
	assert.Equal(t, 32, len(id))
}

func TestGenerateIdSequential(t *testing.T) {
	t.Parallel()
	ids := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		ids[i] = generateId()
		<-time.After(time.Millisecond)
	}
	assert.True(t, sort.StringsAreSorted(ids), "Strings are not sorted")
}

func TestMigrate(t *testing.T) {
	t.Parallel()
	db := SetupDB(t)
	defer db.Close()
	c := 0
	err := db.conn.Get(&c, "SELECT count(*) FROM migrations")
	assert.Nil(t, err, "There is an error counting migrations ", err)
	assert.Greater(t, c, 0)
}

func TestUser(t *testing.T) {
	t.Parallel()
	db := SetupDB(t)
	defer db.Close()
	_, err := db.Users().Insert(types.User{
		Email:    "test@test.com",
		Password: "my-password",
	})
	assert.Nil(t, err, "Unable to insert user")
	user, err := db.Users().FindByEmail("test@test.com")
	assert.Nil(t, err, "Unable to get user")
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

func TestPokemon(t *testing.T) {
	t.Parallel()
	db := SetupDB(t)
	defer db.Close()

	pkmn, err := db.Pokemon().GetAll()
	assert.Nil(t, err, "Unable to fetch all Pokemon")
	assert.Greater(t, len(pkmn), 1020)

	blastoise, err := db.Pokemon().FindByID("blastoise")
	assert.Nil(t, err, "Unable to lookup Blastoise")
	assert.Equal(t, 9, blastoise.NationalDexNumber)
	assert.Equal(t, types.Kanto, blastoise.Generation())

	meowscarada, err := db.Pokemon().FindByID("meowscarada")
	assert.Nil(t, err, "Unable to lookup Meowscarada")
	assert.Equal(t, 908, meowscarada.NationalDexNumber)
	assert.Equal(t, types.Paldea, meowscarada.Generation())

	venaforms, err := db.Pokemon().Forms().FindByPokemonID("venusaur")
	assert.Nil(t, err, "Unable to lookup Venasaur forms")
	assert.Len(t, venaforms, 2)

	first30, err := db.Pokemon().Get(30, 1)
	assert.Nil(t, err, "Failed to get first 30 pokemon")
	assert.Equal(t, "bulbasaur", first30[0].ID)
}
