# 🌤️ CEP Weather

Aplicação simples em Go que permite ao usuário digitar um CEP brasileiro, consultar o endereço via [ViaCEP](https://viacep.com.br), e retornar a temperatura atual da cidade correspondente via [WeatherAPI](https://www.weatherapi.com/).

🚀 Acesse agora:  
👉 **[cep-weather-997632606102.us-central1.run.app](https://cep-weather-997632606102.us-central1.run.app)**

---

## 🧪 Exemplo de uso

1. Acesse o link acima.
2. Digite um CEP válido (ex: `01001-000`).
3. Clique em **Buscar**.
4. O sistema mostrará:
   - CEP pesquisado
   - Cidade correspondente
   - Temperatura atual (°C)

---

## 🛠️ Como rodar localmente

### Pré-requisitos

- Go instalado (versão 1.20 ou superior)
- Variável de ambiente `WEATHER_API_KEY` com sua chave da [WeatherAPI](https://www.weatherapi.com/)
- (Opcional) Arquivo `.env`

### Rodando

```bash
go run main.go
