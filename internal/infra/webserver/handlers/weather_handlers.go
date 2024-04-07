package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/go-chi/chi"
)

type WeatherHandler struct {
	WeatherAPIKey string
}

type Localizacao struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

type CurrentWeather struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
}

type WeatherApiResponse struct {
	Current CurrentWeather `json:"current"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewWeatherHandler(weatherApiKey string) *WeatherHandler {
	return &WeatherHandler{weatherApiKey}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	invalidZipcodeMessage := "Invalid zipcode"
	zipcodeNotFound := "can not find zipcode"
	cep := chi.URLParam(r, "cep")
	isCepValid := isCepValid(cep)

	if !isCepValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(&invalidZipcodeMessage)
		if err != nil {
			panic(err)
		}
		return
	}

	cepResponse := buscaCEP(cep)
	if cepResponse.Cep == "" {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(&zipcodeNotFound)
		if err != nil {
			panic(err)
		}
		return
	}
	log.Printf("CEP encontrado: %v", cepResponse)
	temperatura := buscaTemperatura(h.WeatherAPIKey, cepResponse)
	log.Printf("WeatherApiResponse: %v", temperatura)
	kelvin := temperatura.Current.TempC + 273
	response := WeatherResponse{
		TempC: temperatura.Current.TempC,
		TempF: temperatura.Current.TempF,
		TempK: kelvin,
	}
	json.NewEncoder(w).Encode(&response)
	return
}

func isCepValid(cep string) bool {
	var re = regexp.MustCompile(`^[0-9]+$`)
	if len(cep) != 8 {
		return false
	}
	if !re.MatchString(cep) {
		return false
	}
	return true
}

func buscaCEP(cep string) Localizacao {
	address := "http://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.Get(address)
	if err != nil {
		panic(err)
	}
	if req.StatusCode != http.StatusOK {
		panic("Erro ao fazer requisição para ViaCEP: status code diferente de 200: " + strconv.Itoa(req.StatusCode))
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var data Localizacao
	err = json.Unmarshal(res, &data)
	if err != nil {
		panic(err)
	}
	return data
}

func buscaTemperatura(apiKey string, localizacao Localizacao) WeatherApiResponse {
	location := url.QueryEscape(localizacao.Localidade)
	address := "http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + location + "&aqi=no"
	log.Printf("URL WEATHER API: %s", address)
	req, err := http.Get(address)
	if err != nil {
		panic(err)
	}
	if req.StatusCode != http.StatusOK {
		panic("Erro ao fazer requisição para ViaCEP: status code diferente de 200: " + strconv.Itoa(req.StatusCode))
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var data WeatherApiResponse
	err = json.Unmarshal(res, &data)
	if err != nil {
		panic(err)
	}
	return data
}
