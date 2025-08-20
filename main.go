package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/gomail.v2"
)

type Contato struct {
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Mensagem string `json:"mensagem"`
}

func isEmailValido(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodOptions {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Responder OPTIONS (preflight)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
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

	from := os.Getenv("EMAIL_REMETENTE")
	password := os.Getenv("EMAIL_SENHA")
	to := os.Getenv("EMAIL_DESTINATARIO")

	if from == "" || password == "" || to == "" {
		http.Error(w, "Variáveis de ambiente não configuradas", http.StatusInternalServerError)
		return
	}

	assunto := fmt.Sprintf("portfólio - %s", c.Email)
	corpo := fmt.Sprintf("Nome: %s\n\nMensagem:\n%s", c.Nome, c.Mensagem)

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "Contato do Portfólio")
	m.SetHeader("To", to)
	m.SetHeader("Subject", assunto)
	m.SetBody("text/plain", corpo)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Erro ao enviar email: %v", err)
		http.Error(w, fmt.Sprintf("Erro ao enviar e-mail: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("E-mail enviado com sucesso!"))
}

func main() {
	http.HandleFunc("/contato", Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado na porta %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
