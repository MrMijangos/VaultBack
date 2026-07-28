package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/posts/application"
)

type GetAllPostsController struct {
	useCase *application.GetAllPostsUseCase
}

func NewGetAllPostsController(useCase *application.GetAllPostsUseCase) *GetAllPostsController {
	return &GetAllPostsController{useCase: useCase}
}

func (c *GetAllPostsController) Handle(w http.ResponseWriter, r *http.Request) {
	// Ruta pública (OptionalAuth): claims solo está presente si mandaron un
	// token válido -- sin sesión, userID queda vacío y is_liked/is_saved dan
	// false para todos en el use case, sin romper el feed para visitantes.
	userID := ""
	if claims, ok := security.ClaimsFromContext(r.Context()); ok {
		userID = claims.UserID
	}

	list, err := c.useCase.Execute(r.Context(), userID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}
