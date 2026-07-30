package entities

import "time"

const (
	ServiceRequestTypeServicio   = "servicio"
	ServiceRequestTypeReparacion = "reparacion"

	// ServiceRequestStatusPendienteAceptacion: el dueño mandó la solicitud,
	// el negocio todavía no confirma que tiene el artículo en mano.
	ServiceRequestStatusPendienteAceptacion = "pendiente_aceptacion"
	// ServiceRequestStatusEnEspera: el negocio ya confirmó que lo recibió
	// físicamente, está en cola para empezar.
	ServiceRequestStatusEnEspera = "en_espera"
	// ServiceRequestStatusEnServicio: el negocio ya empezó a trabajarlo.
	ServiceRequestStatusEnServicio = "en_servicio"
	// ServiceRequestStatusTerminado: el negocio terminó -- el dueño todavía
	// no confirma que lo recibió de vuelta.
	ServiceRequestStatusTerminado = "terminado"
	// ServiceRequestStatusConfirmado: el dueño confirmó que ya tiene su
	// artículo de vuelta -- cierra el flujo y dispara la creación
	// automática del registro en maintenance_logs (ver
	// ConfirmServiceRequestUseCase).
	ServiceRequestStatusConfirmado = "confirmado"
)

// ServiceRequest es una solicitud de servicio/reparación iniciada desde el
// chat: el dueño manda un activo a un negocio y ambos van avanzando el
// estado hasta que el dueño confirma que lo recibió de vuelta. A diferencia
// de maintenancelogs (un historial que el dueño carga a mano, después del
// hecho), esto es un flujo en vivo con dos partes.
type ServiceRequest struct {
	ID          string
	AssetID     string
	OwnerID     string
	BusinessID  string
	Type        string
	Status      string
	CreatedAt   time.Time
	AcceptedAt  *time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
	ConfirmedAt *time.Time

	// Denormalizado desde assets/users/businesses para listar sin N+1 (ver
	// PostgreSQLServiceRequestRepository) -- no se persiste en esta tabla.
	AssetName     string
	AssetImageURL string
	OwnerName     string
	BusinessName  string
}
