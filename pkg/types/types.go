package types

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       string
	Email    string
	Password string
	Created  int64
	Updated  int64
}

type Generation uint8

const (
	Kanto Generation = iota + 1
	Johto
	Hoenn
	Sinnoh
	Unova
	Kalos
	Alola
	Galar
	Paldea
	UNKNOWN = 0
)

type Pokemon struct {
	ID                string
	NationalDexNumber int `db:"national_dex_number"`
	Name              string
	ShinyLocked       bool `db:"shiny_locked"`
}

func (p Pokemon) Generation() Generation {
	if p.NationalDexNumber <= 151 {
		return Kanto
	} else if p.NationalDexNumber <= 251 {
		return Johto
	} else if p.NationalDexNumber <= 386 {
		return Hoenn
	} else if p.NationalDexNumber <= 493 {
		return Sinnoh
	} else if p.NationalDexNumber <= 649 {
		return Unova
	} else if p.NationalDexNumber <= 721 {
		return Kalos
	} else if p.NationalDexNumber <= 809 {
		return Alola
	} else if p.NationalDexNumber <= 905 {
		return Galar
	} else if p.NationalDexNumber <= 1025 {
		return Paldea
	}

	return UNKNOWN
}

type PokemonForm struct {
	ID          string
	Name        string
	PokemonID   string `db:"pokemon_id"`
	ShinyLocked bool   `db:"shiny_locked"`
}

type PokemonWithForms struct {
	Pokemon
	Forms []PokemonForm
}

type Pokedex struct {
	ID      string
	Name    string
	OwnerID string `db:"owner_id"`
	Config  string
	Created int64
	Updated int64
}

func NewPokedex(ownerID, name string, config PokedexConfig) (Pokedex, error) {
	c, err := json.Marshal(config)
	if err != nil {
		return Pokedex{}, err
	}
	return Pokedex{
		OwnerID: ownerID,
		Name: name,
		Config:  string(c),
	}, nil
}

type FormLocation int

const (
	Off FormLocation = iota
	After
	Separate
)

func (f FormLocation) Value() string {
	return fmt.Sprint(f)
}
func (f FormLocation) ToString() string {
	if f == Off {
		return "Off"
	} else if f == After {
		return "After base form"
	} else {
		return "Separate"
	}
}

type PokedexConfig struct {
	Shiny         bool
	GenderForms   FormLocation
	RegionalForms FormLocation
	GMaxForms     FormLocation
}

func (p Pokedex) GetConfig() (PokedexConfig, error) {
	var conf PokedexConfig
	err := json.Unmarshal([]byte(p.Config), &conf)
	return conf, err
}
func (p *Pokedex) UpdateConfig(config PokedexConfig) error {
	c, err := json.Marshal(config)
	if err != nil {
		return err
	}
	p.Config = string(c)
	return nil
}

type PokedexEntry struct {
	PokedexID string `db:"pokedex_id"`
	PokemonID string `db:"pokemon_id"`
	FormID    string `db:"form_id"`
	Created   int64
	Updated   int64
}
