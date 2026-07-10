package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/reviews/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createReview *controllers.CreateReviewController,
	getReviewsByProvider *controllers.GetReviewsByProviderController,
	getReviewById *controllers.GetReviewByIdController,
	deleteReview *controllers.DeleteReviewController,
	likeReview *controllers.LikeReviewController,
	unlikeReview *controllers.UnlikeReviewController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/reviews", auth(http.HandlerFunc(createReview.Handle)))
	mux.HandleFunc("GET /api/v1/reviews", getReviewsByProvider.Handle)
	mux.HandleFunc("GET /api/v1/reviews/{id}", getReviewById.Handle)
	mux.Handle("DELETE /api/v1/reviews/{id}", auth(http.HandlerFunc(deleteReview.Handle)))
	mux.Handle("POST /api/v1/reviews/{id}/likes", auth(http.HandlerFunc(likeReview.Handle)))
	mux.Handle("DELETE /api/v1/reviews/{id}/likes", auth(http.HandlerFunc(unlikeReview.Handle)))
}
