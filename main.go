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
		log.Println("Método não permitido:", r.Method)
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Responder OPTIONS (preflight)
	if r.Method == http.MethodOptions {
		log.Println("Preflight OPTIONS recebido")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Decodificar JSON
	var c Contato
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	log.Printf("Payload recebido: %+v\n", c)

	// Validar campos
	if c.Nome == "" || !isEmailValido(c.Email) || strings.TrimSpace(c.Mensagem) == "" {
		log.Println("Validação falhou:", c)
		http.Error(w, "Campos inválidos", http.StatusBadRequest)
		return
	}

	// Variáveis de ambiente
	from := os.Getenv("EMAIL_REMETENTE")
	password := os.Getenv("EMAIL_SENHA")
	to := os.Getenv("EMAIL_DESTINATARIO")

	if from == "" || password == "" || to == "" {
		log.Println("Variáveis de ambiente ausentes:")
		log.Println("EMAIL_REMETENTE:", from)
		log.Println("EMAIL_SENHA:", password != "") // não loga senha em texto puro
		log.Println("EMAIL_DESTINATARIO:", to)
		http.Error(w, "Configuração de e-mail ausente", http.StatusInternalServerError)
		return
	}

	// Montar mensagem
	assunto := fmt.Sprintf("portfólio - %s", c.Email)
	corpo := fmt.Sprintf("Nome: %s\n\nMensagem:\n%s", c.Nome, c.Mensagem)

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "Contato do Portfólio")
	m.SetHeader("To", to)
	m.SetHeader("Subject", assunto)
	m.SetBody("text/plain", corpo)

	// Enviar e-mail
	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)
	if err := d.DialAndSend(m); err != nil {
		log.Println("Erro ao enviar e-mail:", err)
		http.Error(w, "Erro ao enviar e-mail", http.StatusInternalServerError)
		return
	}

	log.Println("E-mail enviado com sucesso!")
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
