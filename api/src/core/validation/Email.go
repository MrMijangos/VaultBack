package validation

import "regexp"

// emailPattern es deliberadamente estricto (alfanuméricos y unos pocos
// símbolos comunes en el local-part, dominio con al menos un punto) --
// antes solo se checaba que el correo contuviera un "@", así que cualquier
// caracter pasaba, incluyendo cosas como "<script>...".
var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail valida formato, no si el correo existe de verdad.
func IsValidEmail(email string) bool {
	return len(email) <= 254 && emailPattern.MatchString(email)
}
