// Package main implementa un servicio HTTP simple con endpoints para
// registrar usuarios y realizar login mediante JWT.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
)

// Usuario representa la estructura de un usuario dentro del sistema.
// Esta implementación simula una base de datos en memoria.
type Usuario struct {
	Correo   string
	Telefono string
	Password string
}

// usuarios es una base de datos simulada en memoria.
var usuarios = []Usuario{}

// jwtKey es la clave secreta utilizada para firmar y verificar los tokens JWT.
var jwtKey = []byte("mi_clave_secreta")

// RegistroRequest define la estructura esperada para la petición
// del endpoint /registro.
type RegistroRequest struct {
	Correo   string `json:"correo"`
	Telefono string `json:"telefono"`
	Password string `json:"password"`
}

// LoginRequest define la estructura esperada para la petición
// del endpoint /login.
type LoginRequest struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

// ErrorResponse define la estructura estándar de respuesta de error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// LoginResponse define la respuesta del login, incluyendo el token
// generado y la fecha de inicio de sesión.
type LoginResponse struct {
	Token       string    `json:"token"`
	FechaInicio time.Time `json:"fecha_inicio"`
}

// validarCorreo revisa que el correo tenga un formato válido.
// Verifica la estructura básica de un email: texto@texto.texto
func validarCorreo(correo string) bool {
	correo = strings.TrimSpace(correo)
	if correo == "" {
		return false
	}

	// Dividir en partes por el @
	partes := strings.Split(correo, "@")
	if len(partes) != 2 {
		return false // Debe tener exactamente un @
	}

	parteLocal := partes[0]
	dominio := partes[1]

	// Validar parte local (antes del @)
	if len(parteLocal) == 0 || len(parteLocal) > 64 {
		return false
	}

	// No puede empezar o terminar con punto
	if strings.HasPrefix(parteLocal, ".") || strings.HasSuffix(parteLocal, ".") {
		return false
	}

	// No puede tener puntos consecutivos
	if strings.Contains(parteLocal, "..") {
		return false
	}

	// Validar dominio
	if len(dominio) == 0 || len(dominio) > 255 {
		return false
	}

	// Debe contener al menos un punto
	if !strings.Contains(dominio, ".") {
		return false
	}

	// No puede empezar o terminar con punto o guión
	if strings.HasPrefix(dominio, ".") || strings.HasSuffix(dominio, ".") ||
		strings.HasPrefix(dominio, "-") || strings.HasSuffix(dominio, "-") {
		return false
	}

	// Validar que tenga extensión después del último punto
	ultimoPunto := strings.LastIndex(dominio, ".")
	extension := dominio[ultimoPunto+1:]
	if len(extension) < 2 {
		return false // La extensión debe tener al menos 2 caracteres
	}

	// Validar caracteres permitidos en el correo
	caracteresValidos := func(s string) bool {
		for _, c := range s {
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
				(c >= '0' && c <= '9') || c == '.' || c == '-' || c == '_' || c == '@') {
				return false
			}
		}
		return true
	}

	return caracteresValidos(correo)
}

// validarTelefono valida que el teléfono tenga exactamente 10 dígitos numéricos.
func validarTelefono(telefono string) bool {
	if len(telefono) != 10 {
		return false
	}
	for _, c := range telefono {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// validarPassword revisa que la contraseña cumpla con:
// - Longitud entre 6 y 12 caracteres
// - Al menos una mayúscula
// - Al menos una minúscula
// - Al menos un número
// - Al menos un carácter especial de la lista "@$&"
func validarPassword(password string) bool {
	if len(password) < 6 || len(password) > 12 {
		return false
	}

	var tieneMayus, tieneMinus, tieneNumero, tieneEspecial bool
	especiales := "@$&"

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			tieneMayus = true
		case unicode.IsLower(c):
			tieneMinus = true
		case unicode.IsDigit(c):
			tieneNumero = true
		case strings.ContainsRune(especiales, c):
			tieneEspecial = true
		}
	}

	return tieneMayus && tieneMinus && tieneNumero && tieneEspecial
}

// registroHandler maneja la creación de nuevos usuarios.
// - Valida los campos recibidos
// - Revisa que no existan usuarios con el mismo correo o teléfono
// - Guarda al usuario en memoria si es válido
func registroHandler(w http.ResponseWriter, r *http.Request) {
	var req RegistroRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Cuerpo inválido"})
		return
	}

	// Validación de campos obligatorios
	if req.Correo == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Falta campo correo en el request.")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Falta el campo correo"})
		return
	}
	if req.Telefono == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Falta campo telefono en el request.")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Falta el campo telefono"})
		return
	}
	if req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Falta campo contraseña en el request.")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Falta el campo contraseña"})
		return
	}

	// Validación de formatos
	if !validarCorreo(req.Correo) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Correo inválido"})
		return
	}
	if !validarTelefono(req.Telefono) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Teléfono inválido"})
		return
	}
	if !validarPassword(req.Password) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Contraseña inválida"})
		return
	}

	// Revisión de duplicados
	for _, u := range usuarios {
		if u.Correo == req.Correo {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "El correo ya se encuentra registrado"})
			return
		}
		if u.Telefono == req.Telefono {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "El teléfono ya se encuentra registrado"})
			return
		}
	}

	// Registro exitoso
	usuarios = append(usuarios, Usuario{
		Correo:   req.Correo,
		Telefono: req.Telefono,
		Password: req.Password,
	})
	fmt.Println("Usuario registrado correctamente")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensaje":"Usuario registrado exitosamente"}`)
}

// loginHandler maneja la autenticación de usuarios.
// - Verifica las credenciales
// - Genera un token JWT válido por 24 horas
// - Responde con el token y la fecha de inicio
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("El cuerpo de la peticion es invalido!!")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Cuerpo inválido"})
		return
	}

	if req.Correo == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Falta campo correo en el request.")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Falta el campo correo"})
		return
	}
	if req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Falta campo contraseña en el request.")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Falta el campo contraseña"})
		return
	}

	// Búsqueda de usuario
	var usuario *Usuario
	for _, u := range usuarios {
		if u.Correo == req.Correo && u.Password == req.Password {
			usuario = &u
			break
		}
	}

	if usuario == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Usuario no encontrado.")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Correo o contraseña incorrectos"})
		return
	}

	// Generación de token JWT
	claims := jwt.MapClaims{
		"correo": usuario.Correo,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error al generar el token")
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error generando token"})
		return
	}

	// Respuesta exitosa
	resp := LoginResponse{
		Token:       tokenString,
		FechaInicio: time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// main inicializa el servidor HTTP en el puerto 8080
// y registra los handlers de /registro y /login.
func main() {
	http.HandleFunc("/registro", registroHandler)
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
