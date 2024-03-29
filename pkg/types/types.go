package types

import "encoding/json"

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
	OwnerID string `db:"owner_id"`
	config  string
	Created int64
	Updated int64
}

func NewPokedex(ownerID string, config PokedexConfig) (Pokedex, error) {
	c, err := json.Marshal(config)
	if err != nil {
		return Pokedex{}, err
	}
	return Pokedex{
		OwnerID: ownerID,
		config:  string(c),
	}, nil
}

type FormLocation int
const (
	Off FormLocation = iota
	Inline
	After
)

type PokedexConfig struct {
	Shiny         bool
	GenderForms   FormLocation
	RegionalForms FormLocation
	GMaxForms     FormLocation
}

func (p Pokedex) Config() (PokedexConfig, error) {
	var conf PokedexConfig
	err := json.Unmarshal([]byte(p.config), &conf)
	return conf, err
}
func (p *Pokedex) UpdateConfig(config PokedexConfig) error {
	c, err := json.Marshal(config)
	if err != nil {
		return err
	}
	p.config = string(c)
	return nil
}

type PokedexEntry struct {
	PokedexID string `db:"pokedex_id"`
	PokemonID string `db:"pokemon_id"`
	FormID    string `db:"form_id"`
	Created   int64
	Updated   int64
}
