package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestGetWeather_notFound(t *testing.T) {
	weatherHandler := NewWeatherHandler()
	r := chi.NewRouter()
	r.Route("/temperatura", func(r chi.Router) {
		r.Get("/{cep}", weatherHandler.GetWeather)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, respBody, err := doRequest("12345678", ts)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
	assert.Equal(t, respBody, "\"can not find zipcode\"\n")
}

func TestGetWeather_invalidByChars(t *testing.T) {
	weatherHandler := NewWeatherHandler()
	r := chi.NewRouter()
	r.Route("/temperatura", func(r chi.Router) {
		r.Get("/{cep}", weatherHandler.GetWeather)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, respBody, err := doRequest("0123-567", ts)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(t, respBody, "\"Invalid zipcode\"\n")
}

func TestGetWeather_invalidByLen(t *testing.T) {
	weatherHandler := NewWeatherHandler()
	r := chi.NewRouter()
	r.Route("/temperatura", func(r chi.Router) {
		r.Get("/{cep}", weatherHandler.GetWeather)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, respBody, err := doRequest("0123", ts)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity)
	assert.Equal(t, respBody, "\"Invalid zipcode\"\n")
}

func TestGetWeather_success(t *testing.T) {
	weatherHandler := NewWeatherHandler()
	r := chi.NewRouter()
	r.Route("/temperatura", func(r chi.Router) {
		r.Get("/{cep}", weatherHandler.GetWeather)
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, respBody, err := doRequest("22460900", ts)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	var data WeatherResponse
	err = json.Unmarshal([]byte(respBody), &data)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	expectedKelvin := data.TempC + 273
	assert.Equal(t, data.TempK, expectedKelvin)
}

func doRequest(value string, svr *httptest.Server) (*http.Response, string, error) {
	req, err := http.NewRequest("GET", svr.URL+"/temperatura/"+value, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	return resp, string(respBody), nil
}
