package types

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/m50/shinidex/pkg/math"
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

type PokemonList []Pokemon

func (p PokemonList) Box(i int) PokemonList {
	start := i * 30
	end := math.Min((i*30)+30, len(p))
	return p[start:end]
}

func (p PokemonList) GMax() PokemonList {
	r := PokemonList{}
	for _, pkmn := range p {
		if strings.Contains(pkmn.ID, "gigantamax") {
			r = append(r, pkmn)
		}
	}
	return r
}

func (p PokemonList) Female() PokemonList {
	r := PokemonList{}
	for _, pkmn := range p {
		if pkmn.ID[len(pkmn.ID)-1:] == "f" {
			r = append(r, pkmn)
		} else if strings.Contains(pkmn.ID, "female") {
			r = append(r, pkmn)
		}
	}
	return r
}

func (p PokemonList) Regional() PokemonList {
	r := PokemonList{}
	for _, pkmn := range p {
		if strings.Contains(pkmn.ID, "alolan") {
			r = append(r, pkmn)
		} else if strings.Contains(pkmn.ID, "galarian") {
			r = append(r, pkmn)
		} else if strings.Contains(pkmn.ID, "paldean") {
			r = append(r, pkmn)
		} else if strings.Contains(pkmn.ID, "hisuian") {
			r = append(r, pkmn)
		}
	}
	return r
}

func (p PokemonList) Forms() PokemonList {
	r := PokemonList{}
	for _, pkmn := range p {
		if strings.Contains(pkmn.ID, "-") {
			if pkmn.ID == "porygon-z" || pkmn.ID == "mime-jr" || pkmn.ID == "mr-mime" {
				continue
			}
			if strings.Contains(pkmn.ID, "tapu-") {
				continue
			}
			if pkmn.ID[len(pkmn.ID)-2:] == "-o" {
				continue
			}
			if pkmn.NationalDexNumber >= 984 && pkmn.NationalDexNumber <= 995 {
				continue
			}
			if pkmn.NationalDexNumber >= 1001 && pkmn.NationalDexNumber <= 1010 {
				continue
			}
			if pkmn.NationalDexNumber >= 1020 && pkmn.NationalDexNumber <= 1023 {
				continue
			}
			r = append(r, pkmn)
		}
	}
	return r
}

type Pokemon struct {
	ID                string
	NationalDexNumber int `db:"national_dex_number"`
	Name              string
	ShinyLocked       bool `db:"shiny_locked"`
	Form              bool
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

func (p Pokemon) GetLocalImagePath(shiny bool) string {
	id := strings.Replace(p.ID, "+", "-", 1)
	id, _ = strings.CutSuffix(id, "-antique")
	id, _ = strings.CutSuffix(id, "-masterpiece")
	id, _ = strings.CutSuffix(id, "-artisan")
	shinyStr := "normal"
	if shiny {
		shinyStr = "shiny"
	}

	workingDir, _ := os.Getwd()
	localFile := fmt.Sprintf("/assets/imgs/%s/%s.png", shinyStr, id)

	return workingDir + localFile
}

func (p Pokemon) GetImageURL(shiny bool) string {
	id := strings.Replace(p.ID, "+", "-", 1)
	id, _ = strings.CutSuffix(id, "-antique")
	id, _ = strings.CutSuffix(id, "-masterpiece")
	id, _ = strings.CutSuffix(id, "-artisan")
	shinyStr := "normal"
	if shiny {
		shinyStr = "shiny"
	}

	workingDir, _ := os.Getwd()
	localFile := fmt.Sprintf("/assets/imgs/%s/%s.png", shinyStr, id)
	if _, err := os.Stat(workingDir + localFile); err != nil {
		return fmt.Sprintf("https://img.pokemondb.net/sprites/home/%s/%s.png", shinyStr, id)
	}

	return localFile
}

type PokemonForm struct {
	ID          string
	Name        string
	PokemonID   string `db:"pokemon_id"`
	ShinyLocked bool   `db:"shiny_locked"`
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
		Name:    name,
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
	Forms         FormLocation
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
