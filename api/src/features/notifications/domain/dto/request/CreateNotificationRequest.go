package request

import (
	"encoding/json"
	"errors"
)

var allowedNotificationTypes = map[string]bool{
	"servicio":   true,
	"reparacion": true,
	"venta":      true,
	"blockchain": true,
	"comunidad":  true,
}

var allowedNotificationSubtypes = map[string]bool{
	"entro_servicio":   true,
	"salio_servicio":   true,
	"entro_reparacion": true,
	"salio_reparacion": true,
	"pedido_recibido":  true,
	"pedido_enviado":   true,
	"nueva_compra":     true,
	"asset_verificado": true,
	"likes_post":       true,
}

type CreateNotificationRequest struct {
	Type    string          `json:"type"`
	Subtype string          `json:"subtype"`
	Title   string          `json:"title"`
	Body    string          `json:"body"`
	Data    json.RawMessage `json:"data"`
}

func (r CreateNotificationRequest) Validate() error {
	if !allowedNotificationTypes[r.Type] {
		return errors.New("el tipo de notificacion no es valido")
	}
	if !allowedNotificationSubtypes[r.Subtype] {
		return errors.New("el subtipo de notificacion no es valido")
	}
	if r.Title == "" {
		return errors.New("el titulo es obligatorio")
	}
	if r.Body == "" {
		return errors.New("el cuerpo es obligatorio")
	}
	return nil
}
