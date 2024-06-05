package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Estrutura para representar um usuário
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

// Função para ler os usuários do arquivo JSON
func loadUsers() error {
	file, err := ioutil.ReadFile("users.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &users)
}

// Função para salvar os usuários no arquivo JSON
func saveUsers() error {
	file, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("users.json", file, 0644)
}

// Função para verificar se um usuário existe
func userExists(username, password string) bool {
	for _, user := range users {
		if user.Username == username && user.Password == password {
			return true
		}
	}
	return false
}

// Manipulador para a rota de login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Se o método da requisição for GET, renderizamos o formulário de login
	if r.Method == "GET" {
		fmt.Fprintf(w, `
            <h1>Login</h1>
            <form method="post" action="/login">
                <label for="username">Nome de Usuário:</label><br>
                <input type="text" id="username" name="username"><br>
                <label for="password">Senha:</label><br>
                <input type="password" id="password" name="password"><br><br>
                <input type="submit" value="Login">
            </form>
            <p>Não possui uma conta? <a href="/register">Registre-se aqui</a>.</p>
        `)
		return
	}

	// Se o método da requisição for POST, processamos os dados do formulário
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Verificando se o usuário existe
	if !userExists(username, password) {
		// Se o usuário não existir, exibimos uma mensagem de erro
		http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}

	// Redirecionando para o serviço de catálogo
	catalogURL := "http://localhost:8081/catalogo"
	http.Redirect(w, r, catalogURL, http.StatusSeeOther)
}

// Manipulador para a rota de registro
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Se o método da requisição for GET, renderizamos o formulário de registro
	if r.Method == "GET" {
		fmt.Fprintf(w, `
            <h1>Registro</h1>
            <form method="post" action="/register">
                <label for="username">Nome de Usuário:</label><br>
                <input type="text" id="username" name="username"><br>
                <label for="password">Senha:</label><br>
                <input type="password" id="password" name="password"><br><br>
                <input type="submit" value="Registrar">
            </form>
        `)
		return
	}

	// Se o método da requisição for POST, processamos os dados do formulário
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Verificando se o nome de usuário já está em uso
	for _, user := range users {
		if user.Username == username {
			http.Error(w, "Usuário já registrado", http.StatusBadRequest)
			return
		}
	}

	// Adicionando o novo usuário à lista
	newUser := User{Username: username, Password: password}
	users = append(users, newUser)

	// Salvando os usuários atualizados no arquivo JSON
	err := saveUsers()
	if err != nil {
		http.Error(w, "Erro ao salvar os usuários", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Usuário registrado com sucesso: %s", username)
}

// Manipulador para servir arquivos estáticos
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	// Carregando os usuários do arquivo JSON
	err := loadUsers()
	if err != nil {
		fmt.Println("Erro ao carregar os usuários", err)
		return
	}

	// Configurando rota para o login
	http.HandleFunc("/login", loginHandler)

	// Configurando rota para o registro de usuários
	http.HandleFunc("/register", registerHandler)

	// Configurando rota para servir arquivos estáticos (página HTML)
	http.HandleFunc("/", staticHandler)

	// Iniciando o servidor na porta 8080
	fmt.Println("Servidor rodando na porta 8080, Link: http://localhost:8080/login")
	http.ListenAndServe(":8080", nil)
}
