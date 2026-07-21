package main

import (
	"fmt"
	"log"
	"net/http"

	"vault/src/core/cloudinary"
	"vault/src/core/config"
	"vault/src/core/eventbus"
	"vault/src/core/middleware"
	"vault/src/core/moderation"
	assetsInfra "vault/src/features/assets/infrastructure"
	assetsRouter "vault/src/features/assets/infrastructure/router"
	authInfra "vault/src/features/auth/infrastructure"
	authRouter "vault/src/features/auth/infrastructure/router"
	blockchaincertificatesInfra "vault/src/features/blockchaincertificates/infrastructure"
	blockchaincertificatesRouter "vault/src/features/blockchaincertificates/infrastructure/router"
	businessesInfra "vault/src/features/businesses/infrastructure"
	businessesRouter "vault/src/features/businesses/infrastructure/router"
	commentsInfra "vault/src/features/comments/infrastructure"
	commentsRouter "vault/src/features/comments/infrastructure/router"
	maintenancelogsInfra "vault/src/features/maintenancelogs/infrastructure"
	maintenancelogsRouter "vault/src/features/maintenancelogs/infrastructure/router"
	notificationsInfra "vault/src/features/notifications/infrastructure"
	notificationsRouter "vault/src/features/notifications/infrastructure/router"
	postsInfra "vault/src/features/posts/infrastructure"
	postsRouter "vault/src/features/posts/infrastructure/router"
	reviewsInfra "vault/src/features/reviews/infrastructure"
	reviewsRouter "vault/src/features/reviews/infrastructure/router"
	usersInfra "vault/src/features/users/infrastructure"
	usersRouter "vault/src/features/users/infrastructure/router"
)

func main() {
	fmt.Println("¡Módulo Vault inicializado correctamente!")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("configuracion invalida: %v", err)
	}

	pool, err := config.NewPoolConnection(cfg)
	if err != nil {
		log.Fatalf("error de conexion a base de datos: %v", err)
	}
	defer pool.Close()

	imageUploader, err := cloudinary.NewImageUploader(cfg)
	if err != nil {
		log.Fatalf("error al inicializar cloudinary: %v", err)
	}

	if err := config.RunMigrations(pool); err != nil {
		log.Fatalf("error al migrar el esquema: %v", err)
	}

	// Publisher de eventos hacia vault-ai-service (NLP/ML). Si RabbitMQ no
	// está disponible, se usa un publisher "noop" y la API sigue funcionando
	// con normalidad (solo no se dispara NLP/ML para contenido nuevo).
	var publisher eventbus.Publisher
	rabbitPublisher, err := eventbus.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("advertencia: RabbitMQ no disponible (%v), NLP/ML no se disparará hasta reiniciar la API", err)
		publisher = eventbus.NoopPublisher{}
	} else {
		defer rabbitPublisher.Close()
		publisher = rabbitPublisher
		eventbus.StartNLPAnalyzedConsumer(cfg.RabbitMQURL, pool)
	}

	// Cliente síncrono de moderación: se llama antes de guardar un post,
	// comentario o reseña. Si el contenido es tóxico, o si el servicio de
	// NLP no responde, la publicación se rechaza (no queda guardada).
	moderationClient := moderation.NewClient(cfg.NLPServiceURL)

	mux := http.NewServeMux()

	usersRouter.RegisterRoutes(
		mux,
		usersInfra.BuildCreateUserController(pool, cfg.JWTSecret, cfg.CookieSecure),
		usersInfra.BuildGetAllUsersController(pool),
		usersInfra.BuildGetUserByIdController(pool),
		usersInfra.BuildUpdateUserController(pool),
		usersInfra.BuildDeleteUserController(pool),
		usersInfra.BuildUploadUserImageController(pool, imageUploader),
		cfg.JWTSecret,
	)

	authRouter.RegisterRoutes(mux, authInfra.BuildLoginController(pool, cfg.JWTSecret, cfg.CookieSecure))

	assetsRouter.RegisterRoutes(
		mux,
		assetsInfra.BuildCreateAssetController(pool, publisher),
		assetsInfra.BuildGetAllAssetsController(pool),
		assetsInfra.BuildGetAssetByIdController(pool),
		assetsInfra.BuildUpdateAssetController(pool, publisher),
		assetsInfra.BuildDeleteAssetController(pool),
		assetsInfra.BuildUploadAssetPhotoController(pool, imageUploader),
		cfg.JWTSecret,
	)

	businessesRouter.RegisterRoutes(
		mux,
		businessesInfra.BuildCreateBusinessController(pool),
		businessesInfra.BuildGetAllBusinessesController(pool),
		businessesInfra.BuildGetBusinessByIdController(pool),
		businessesInfra.BuildUpdateBusinessController(pool),
		businessesInfra.BuildDeleteBusinessController(pool),
		cfg.JWTSecret,
	)

	maintenancelogsRouter.RegisterRoutes(
		mux,
		maintenancelogsInfra.BuildCreateMaintenanceLogController(pool),
		maintenancelogsInfra.BuildGetLogsByAssetController(pool),
		maintenancelogsInfra.BuildGetMaintenanceLogByIdController(pool),
		maintenancelogsInfra.BuildUpdateMaintenanceLogController(pool),
		maintenancelogsInfra.BuildDeleteMaintenanceLogController(pool),
		cfg.JWTSecret,
	)

	blockchaincertificatesRouter.RegisterRoutes(
		mux,
		blockchaincertificatesInfra.BuildCreateBlockchainCertificateController(pool),
		blockchaincertificatesInfra.BuildGetCertificatesByAssetController(pool),
		blockchaincertificatesInfra.BuildGetCertificateByIdController(pool),
		cfg.JWTSecret,
	)

	postsRouter.RegisterRoutes(
		mux,
		postsInfra.BuildCreatePostController(pool, moderationClient),
		postsInfra.BuildGetAllPostsController(pool),
		postsInfra.BuildGetPostByIdController(pool),
		postsInfra.BuildUpdatePostController(pool),
		postsInfra.BuildDeletePostController(pool),
		postsInfra.BuildUploadPostPhotoController(pool, imageUploader),
		postsInfra.BuildLikePostController(pool),
		postsInfra.BuildUnlikePostController(pool),
		postsInfra.BuildSavePostController(pool),
		postsInfra.BuildUnsavePostController(pool),
		postsInfra.BuildGetSavedPostsController(pool),
		cfg.JWTSecret,
	)

	commentsRouter.RegisterRoutes(
		mux,
		commentsInfra.BuildCreateCommentController(pool, moderationClient),
		commentsInfra.BuildGetCommentsByPostController(pool),
		commentsInfra.BuildDeleteCommentController(pool),
		cfg.JWTSecret,
	)

	reviewsRouter.RegisterRoutes(
		mux,
		reviewsInfra.BuildCreateReviewController(pool, moderationClient),
		reviewsInfra.BuildGetReviewsByProviderController(pool),
		reviewsInfra.BuildGetProviderRatingController(pool),
		reviewsInfra.BuildGetReviewByIdController(pool),
		reviewsInfra.BuildDeleteReviewController(pool),
		reviewsInfra.BuildLikeReviewController(pool),
		reviewsInfra.BuildUnlikeReviewController(pool),
		cfg.JWTSecret,
	)

	notificationsRouter.RegisterRoutes(
		mux,
		notificationsInfra.BuildCreateNotificationController(pool),
		notificationsInfra.BuildGetMyNotificationsController(pool),
		notificationsInfra.BuildMarkNotificationAsReadController(pool),
		notificationsInfra.BuildDeleteNotificationController(pool),
		cfg.JWTSecret,
	)

	handler := middleware.CORS(cfg.CORSOrigin)(mux)

	fmt.Println("API Vault iniciada correctamente.")
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, handler))
}
