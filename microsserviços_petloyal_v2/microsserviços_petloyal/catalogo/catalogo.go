package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/mux"
)

type Catalogo struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Idade    int    `json:"idade"`
	Raca     string `json:"raca"`
	Genero   string `json:"genero"`
	Porte    string `json:"porte"`
	Cidade   string `json:"cidade"`
	Telefone string `json:"telefone"`
}

type Catalogos struct {
	Catalogo []Catalogo `json:"catalogo"`
}

func loadData() ([]Catalogo, error) {
	file, err := os.ReadFile("catalogos.json")
	if err != nil {
		return nil, err
	}

	var catalogos Catalogos
	err = json.Unmarshal(file, &catalogos)
	if err != nil {
		return nil, err
	}

	return catalogos.Catalogo, nil
}

func ListCatalogos(w http.ResponseWriter, r *http.Request) {
	catalogos, err := loadData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Obter parâmetros de consulta da URL
	city := r.URL.Query().Get("city")
	gender := r.URL.Query().Get("gender")
	age := r.URL.Query().Get("age")

	// Aplicar filtros se os parâmetros estiverem presentes
	if city != "" {
		catalogos = filterByCity(catalogos, city)
	}
	if gender != "" {
		catalogos = filterByGender(catalogos, gender)
	}
	if age != "" {
		catalogos = filterByAge(catalogos, age)
	}

	// Usar a variável catalogos
	response, err := json.MarshalIndent(catalogos, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func filterByCity(catalogos []Catalogo, city string) []Catalogo {
	var filteredCatalogos []Catalogo

	for _, catalogo := range catalogos {
		if catalogo.Cidade == city {
			filteredCatalogos = append(filteredCatalogos, catalogo)
		}
	}

	return filteredCatalogos
}

func filterByGender(catalogos []Catalogo, gender string) []Catalogo {
	var filteredCatalogos []Catalogo

	for _, catalogo := range catalogos {
		if catalogo.Genero == gender {
			filteredCatalogos = append(filteredCatalogos, catalogo)
		}
	}

	return filteredCatalogos
}

func filterByAge(catalogos []Catalogo, age string) []Catalogo {
	// Se age for "youngest", ordenar por idade crescente
	if age == "youngest" {
		sort.Slice(catalogos, func(i, j int) bool {
			return catalogos[i].Idade < catalogos[j].Idade
		})
	} else if age == "oldest" {
		// Se age for "oldest", ordenar por idade decrescente
		sort.Slice(catalogos, func(i, j int) bool {
			return catalogos[i].Idade > catalogos[j].Idade
		})
	}

	return catalogos
}

func main() {
	r := mux.NewRouter()

	// Servir o arquivo HTML estático
	r.HandleFunc("/catalogo", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "catalogo.html")
	}).Methods("GET")

	// Rota para a API de catálogo de animais
	r.HandleFunc("/catalogos", ListCatalogos).Methods("GET")

	// Servir arquivos estáticos (por exemplo, CSS, JavaScript)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server is running on port 8081...")

	http.ListenAndServe(":8081", r)
}
