package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type Contato struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Mensagem string `json:"mensagem"`
}

func isEmailValido(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func contatoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var c Contato
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if c.Nome == "" || !isEmailValido(c.Email) || strings.TrimSpace(c.Mensagem) == "" {
		http.Error(w, "Campos inválidos", http.StatusBadRequest)
		return
	}

	// Variáveis de ambiente
	from := os.Getenv("EMAIL_REMETENTE")
	password := os.Getenv("EMAIL_SENHA")
	to := os.Getenv("EMAIL_DESTINATARIO")

	assunto := fmt.Sprintf("portifólio - %s", c.Email)
	corpo := fmt.Sprintf("Nome: %s\n\nMensagem:\n%s", c.Nome, c.Mensagem)

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "Contato do Portfólio")
	m.SetHeader("To", to)
	m.SetHeader("Subject", assunto)
	m.SetBody("text/plain", corpo)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Erro ao enviar e-mail:", err)
		http.Error(w, "Erro ao enviar e-mail", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("E-mail enviado com sucesso!"))
}

func main() {
	// Carrega variáveis do .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env:", err)
	}

	http.HandleFunc("/contato", contatoHandler)
	fmt.Println("API rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
