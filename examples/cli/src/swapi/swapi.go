package swapi

import (
	"net/http"
	"net/url"
	"time"

	"github.com/68696c6c/goat"
	"github.com/pkg/errors"
)

type Person struct {
	Name      string    `json:"name"`
	Height    string    `json:"height"`
	Mass      string    `json:"mass"`
	HairColor string    `json:"hair_color"`
	SkinColor string    `json:"skin_color"`
	EyeColor  string    `json:"eye_color"`
	BirthYear string    `json:"birth_year"`
	Gender    string    `json:"gender"`
	Homeworld string    `json:"homeworld"`
	Films     []string  `json:"films"`
	Species   []string  `json:"species"`
	Vehicles  []string  `json:"vehicles"`
	Starships []string  `json:"starships"`
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	Url       string    `json:"url"`
}

type Planet struct {
	Name           string    `json:"name"`
	RotationPeriod string    `json:"rotation_period"`
	OrbitalPeriod  string    `json:"orbital_period"`
	Diameter       string    `json:"diameter"`
	Climate        string    `json:"climate"`
	Gravity        string    `json:"gravity"`
	Terrain        string    `json:"terrain"`
	SurfaceWater   string    `json:"surface_water"`
	Population     string    `json:"population"`
	Residents      []string  `json:"residents"`
	Films          []string  `json:"films"`
	Created        time.Time `json:"created"`
	Edited         time.Time `json:"edited"`
	Url            string    `json:"url"`
}

type Starship struct {
	Name                 string    `json:"name"`
	Model                string    `json:"model"`
	Manufacturer         string    `json:"manufacturer"`
	CostInCredits        string    `json:"cost_in_credits"`
	Length               string    `json:"length"`
	MaxAtmospheringSpeed string    `json:"max_atmosphering_speed"`
	Crew                 string    `json:"crew"`
	Passengers           string    `json:"passengers"`
	CargoCapacity        string    `json:"cargo_capacity"`
	Consumables          string    `json:"consumables"`
	HyperdriveRating     string    `json:"hyperdrive_rating"`
	MGLT                 string    `json:"MGLT"`
	StarshipClass        string    `json:"starship_class"`
	Pilots               []string  `json:"pilots"`
	Films                []string  `json:"films"`
	Created              time.Time `json:"created"`
	Edited               time.Time `json:"edited"`
	Url                  string    `json:"url"`
}

type Swapi interface {
	GetPerson(id string) (Person, error)
	GetPlanet(id string) (Planet, error)
	GetStarship(id string) (Starship, error)
}

func NewSwapi() (Swapi, error) {
	baseUrl, err := url.Parse("https://swapi.dev/api/")
	if err != nil {
		return nil, err
	}
	return swapi{
		baseUrl: baseUrl,
	}, nil
}

type swapi struct {
	baseUrl *url.URL
}

func (s swapi) get(path, id string, result any) error {
	u := s.baseUrl.JoinPath(path).JoinPath(id)

	req, err := goat.NewRequest(u.String())
	if err != nil {
		return errors.Wrap(err, "failed to initialize request")
	}
	req.SetMethod(http.MethodGet)

	_, err = req.SendAndRead(result)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}

	return nil
}

func (s swapi) GetPerson(id string) (Person, error) {
	result := Person{}
	err := s.get("people", id, &result)
	if err != nil {
		return Person{}, errors.Wrap(err, "failed to get person")
	}
	return result, nil
}

func (s swapi) GetPlanet(id string) (Planet, error) {
	result := Planet{}
	err := s.get("planets", id, &result)
	if err != nil {
		return Planet{}, errors.Wrap(err, "failed to get planet")
	}
	return result, nil
}

func (s swapi) GetStarship(id string) (Starship, error) {
	result := Starship{}
	err := s.get("starships", id, &result)
	if err != nil {
		return Starship{}, errors.Wrap(err, "failed to get starship")
	}
	return result, nil
}
