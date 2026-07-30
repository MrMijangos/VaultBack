package application

import "vault/src/features/servicerequests/domain/entities"

// typeLabel da el texto en español (con tilde) para mostrar en notificaciones
// -- el valor crudo ("servicio"/"reparacion") es el que se guarda y valida
// contra el CHECK de la base, este es solo para el texto legible.
func typeLabel(t string) string {
	if t == entities.ServiceRequestTypeReparacion {
		return "reparación"
	}
	return "servicio"
}

// subtypeFor arma el subtype de notificación siguiendo el patrón ya
// reservado en notifications_subtype_check (entro_servicio/entro_reparacion,
// salio_servicio/salio_reparacion) y lo extiende con el mismo prefijo para
// el estado intermedio (en_proceso_servicio/en_proceso_reparacion).
func subtypeFor(prefix string, t string) string {
	if t == entities.ServiceRequestTypeReparacion {
		return prefix + "_reparacion"
	}
	return prefix + "_servicio"
}

// articleFor/definiteArticleFor existen porque "servicio" es masculino y
// "reparación" es femenino -- concatenar typeLabel directo después de "un"/
// "el" en los mensajes produce "un reparación"/"el reparación", mal
// concordados en español.
func articleFor(t string) string {
	if t == entities.ServiceRequestTypeReparacion {
		return "una"
	}
	return "un"
}

func definiteArticleFor(t string) string {
	if t == entities.ServiceRequestTypeReparacion {
		return "la"
	}
	return "el"
}
