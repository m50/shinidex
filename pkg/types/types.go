package types

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/m50/shinidex/pkg/math"
)

type User struct {
	ID       string
	Email    string
	Password string
	Created  time.Time
	Updated  time.Time
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

func (p PokemonList) IDs() []string {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return ids
}

func (p PokemonList) GMax() PokemonList {
	return slices.DeleteFunc(slices.Clone(p), func(pkmn Pokemon) bool {
		return !pkmn.IsGMax()
	})
}

func (p PokemonList) Female() PokemonList {
	return slices.DeleteFunc(slices.Clone(p), func(pkmn Pokemon) bool {
		return !pkmn.IsFemale()
	})
}

func (p PokemonList) Regional() PokemonList {
	return slices.DeleteFunc(slices.Clone(p), func(pkmn Pokemon) bool {
		return !pkmn.IsRegional()
	})
}

func (p PokemonList) StandardForms() PokemonList {
	return slices.DeleteFunc(slices.Clone(p), func(pkmn Pokemon) bool {
		return !pkmn.IsStandardForm()
	})
}

type Pokemon struct {
	ID                string
	NationalDexNumber int `db:"national_dex_number"`
	Name              string
	ShinyLocked       bool `db:"shiny_locked"`
	Form              bool
	Caught            bool
}

func (p Pokemon) IDParts() (pokemonID, formID string) {
	parts := strings.Split(p.ID, "+")
	pokemonID = parts[0]
	if len(parts) > 1 {
		formID = parts[1]
	}

	return
}

func (p Pokemon) IsGMax() bool {
	_, formID := p.IDParts()
	return strings.Contains(formID, "gigantamax")
}

func (p Pokemon) IsFemale() bool {
	_, formID := p.IDParts()
	return formID == "f" || strings.Contains(formID, "female")
}

func (p Pokemon) IsRegional() bool {
	_, formID := p.IDParts()
	return strings.Contains(formID, "alolan") ||
		strings.Contains(formID, "galarian") ||
		strings.Contains(formID, "paldean") ||
		strings.Contains(formID, "hisuian")
}

func (p Pokemon) IsStandardForm() bool {
	if p.IsFemale() || p.IsGMax() || p.IsRegional() {
		return false
	}
	return strings.Contains(p.ID, "+")
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
	config  PokedexConfig
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

func (f FormLocation) Selected(opt FormLocation) string {
	if f == opt {
		return "selected=\"true\""
	}
	return ""
}

func (f FormLocation) Value() string {
	return fmt.Sprintf("%d", f)
}
func (f FormLocation) String() string {
	if f == Off {
		return "Off"
	} else if f == After {
		return "With base form"
	} else {
		return "After Pokedex"
	}
}

func (f FormLocation) Off() bool {
	return f == Off
}

func (f FormLocation) AfterBaseForm() bool {
	return f == After
}

func (f FormLocation) Separate() bool {
	return f == Separate
}

type PokedexConfig struct {
	Shiny         bool
	Forms         FormLocation
	GenderForms   FormLocation
	RegionalForms FormLocation
	GMaxForms     FormLocation
	set           bool
}

func (p Pokedex) GetConfig() (PokedexConfig, error) {
	var err error
	if !p.config.set {
		err = json.Unmarshal([]byte(p.Config), &p.config)
		if err == nil {
			p.config.set = true
		}
	}
	return p.config, err
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
