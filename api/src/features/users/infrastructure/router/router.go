package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/users/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createUser *controllers.CreateUserController,
	getAllUsers *controllers.GetAllUsersController,
	getUserById *controllers.GetUserByIdController,
	updateUser *controllers.UpdateUserController,
	deleteUser *controllers.DeleteUserController,
	uploadUserImage *controllers.UploadUserImageController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.HandleFunc("POST /api/v1/users", createUser.Handle)
	mux.Handle("GET /api/v1/users", auth(http.HandlerFunc(getAllUsers.Handle)))
	mux.Handle("GET /api/v1/users/{id}", auth(http.HandlerFunc(getUserById.Handle)))
	mux.Handle("PUT /api/v1/users/{id}", auth(http.HandlerFunc(updateUser.Handle)))
	mux.Handle("DELETE /api/v1/users/{id}", auth(http.HandlerFunc(deleteUser.Handle)))
	mux.Handle("PUT /api/v1/users/{id}/image", auth(http.HandlerFunc(uploadUserImage.Handle)))
}
