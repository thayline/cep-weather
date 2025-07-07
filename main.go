package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
	"golang.org/x/text/unicode/norm"
)

type PageData struct {
	Result string
	Error  string
	Cidade string
	TempC  string
	TempF  string
	TempK  string
}

type brasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type weather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		WindchillC float64 `json:"windchill_c"`
		WindchillF float64 `json:"windchill_f"`
		HeatindexC float64 `json:"heatindex_c"`
		HeatindexF float64 `json:"heatindex_f"`
		DewpointC  float64 `json:"dewpoint_c"`
		DewpointF  float64 `json:"dewpoint_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

var apiKey string

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar .env")
	}
	apiKey = os.Getenv("WEATHER_API_KEY")

	http.HandleFunc("/", inputHandler)
	http.ListenAndServe(":8080", nil)
}

func inputHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	if r.Method == http.MethodPost {
		r.ParseForm()
		cep := r.FormValue("cep")
		data.Result = cep

		body, err := brasilApiRequest(cep)
		if err != nil {
			fmt.Println(err)
			data.Error = fmt.Sprintf("Erro na requisição: %v", err)
			showTemplate(w, data)
			return
		}

		var brasilApi_result brasilApi
		err = json.Unmarshal(body, &brasilApi_result)
		if err != nil {
			fmt.Println(err)
			data.Error = fmt.Sprintf("CEP inválido: %v", err)
			showTemplate(w, data)
			return
		}

		data.Cidade = brasilApi_result.City
		if data.Cidade == "" {
			data.Error = "Cidade não encontrada para o CEP informado."
			showTemplate(w, data)
			return
		}
		cidadeSanitizada := sanitizeCidade(data.Cidade)

		body, err = weatherRequest(cidadeSanitizada)
		if err != nil {
			data.Error = fmt.Sprintf("Erro na requisição: %v", err)
			showTemplate(w, data)
			return
		}

		var temp weather
		err = json.Unmarshal(body, &temp)
		if err != nil {
			data.Error = fmt.Sprintf("Cidade inválida: %v", err)
		}

		data.TempC = fmt.Sprintf("%.1f", temp.Current.TempC)
		data.TempF = fmt.Sprintf("%.1f", temp.Current.TempF)
		data.TempK = fmt.Sprintf("%.1f", CtoK(temp.Current.TempC))

	}

	showTemplate(w, data)
}

func showTemplate(w http.ResponseWriter, data PageData) {
	tmpl := template.Must(template.New("form").Parse(`
		<!DOCTYPE html>
		<html>
		<head><title>CEP-WEATHER</title></head>
		<body>
			<h2>Digite o CEP:</h2>
			<form method="POST">
				<input type="text" name="cep" />
				<input type="submit" value="Buscar" />
			</form>

			{{if .Result}}<p>CEP pesquisado: {{.Result}}</p>{{end}}
			{{if .Cidade}}<p>Cidade: {{.Cidade}}</p>{{end}}
			{{if .TempC}}<p>Temperatura em C°: {{.TempC}}</p>{{end}}
			{{if .TempF}}<p>Temperatura em F°: {{.TempF}}</p>{{end}}
			{{if .TempK}}<p>Temperatura em K: {{.TempK}}</p>{{end}}
			{{if .Error}}<p style="color:red;">Erro: {{.Error}}</p>{{end}}
		</body>
		</html>
	`))

	tmpl.Execute(w, data)
}

func brasilApiRequest(cep string) ([]byte, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CEP não encontrado ou inválido")
	}

	return body, nil
}

func weatherRequest(cidade string) ([]byte, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, cidade)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Cidade não encontrada ou inválida")
	}

	return body, nil
}

func removeAcentos(str string) string {
	t := norm.NFD.String(str)
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func sanitizeCidade(input string) string {
	semAcento := removeAcentos(input)

	var b strings.Builder
	for _, r := range semAcento {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}

	limpo := strings.Join(strings.Fields(b.String()), "+")
	return limpo
}

func CtoK(tempC float64) float64 {
	return tempC + 273
}
