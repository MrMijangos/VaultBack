package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"vault/src/core/cloudinary"
	"vault/src/core/config"
	"vault/src/core/eventbus"
	"vault/src/core/middleware"
	"vault/src/core/moderation"
	"vault/src/core/push"
	addressesInfra "vault/src/features/addresses/infrastructure"
	addressesRouter "vault/src/features/addresses/infrastructure/router"
	assetcommentsInfra "vault/src/features/assetcomments/infrastructure"
	assetcommentsRouter "vault/src/features/assetcomments/infrastructure/router"
	assetsInfra "vault/src/features/assets/infrastructure"
	assetsRouter "vault/src/features/assets/infrastructure/router"
	authInfra "vault/src/features/auth/infrastructure"
	authRouter "vault/src/features/auth/infrastructure/router"
	blockchaincertificatesInfra "vault/src/features/blockchaincertificates/infrastructure"
	blockchaincertificatesRouter "vault/src/features/blockchaincertificates/infrastructure/router"
	businessesInfra "vault/src/features/businesses/infrastructure"
	businessesRouter "vault/src/features/businesses/infrastructure/router"
	businessservicesInfra "vault/src/features/businessservices/infrastructure"
	businessservicesRouter "vault/src/features/businessservices/infrastructure/router"
	chatInfra "vault/src/features/chat/infrastructure"
	chatRouter "vault/src/features/chat/infrastructure/router"
	commentsInfra "vault/src/features/comments/infrastructure"
	commentsRouter "vault/src/features/comments/infrastructure/router"
	fcmtokensInfra "vault/src/features/fcmtokens/infrastructure"
	fcmtokensRouter "vault/src/features/fcmtokens/infrastructure/router"
	maintenancelogsInfra "vault/src/features/maintenancelogs/infrastructure"
	maintenancelogsRouter "vault/src/features/maintenancelogs/infrastructure/router"
	notificationsInfra "vault/src/features/notifications/infrastructure"
	notificationsRouter "vault/src/features/notifications/infrastructure/router"
	postsInfra "vault/src/features/posts/infrastructure"
	postsRouter "vault/src/features/posts/infrastructure/router"
	recommendationsInfra "vault/src/features/recommendations/infrastructure"
	recommendationsRouter "vault/src/features/recommendations/infrastructure/router"
	restorerprofilesInfra "vault/src/features/restorerprofiles/infrastructure"
	restorerprofilesRouter "vault/src/features/restorerprofiles/infrastructure/router"
	reviewsInfra "vault/src/features/reviews/infrastructure"
	reviewsAdapters "vault/src/features/reviews/infrastructure/adapters"
	reviewsRouter "vault/src/features/reviews/infrastructure/router"
	servicerequestsInfra "vault/src/features/servicerequests/infrastructure"
	servicerequestsRouter "vault/src/features/servicerequests/infrastructure/router"
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

	// Sender de notificaciones push (FCM). Si falta la cuenta de servicio, o
	// si Firebase no pudo inicializarse, se usa un sender "noop" y la API
	// sigue funcionando con normalidad -- las notificaciones dentro de la
	// app se siguen creando, solo no se manda el push.
	fcmTokenProvider := fcmtokensInfra.BuildFCMTokenProvider(pool)
	var pushSender push.Sender
	if cfg.FirebaseServiceAccountKey == "" {
		log.Println("advertencia: FIREBASE_SERVICE_ACCOUNT_KEY no configurado, las notificaciones push no se enviarán")
		pushSender = push.NoopSender{}
	} else {
		fcmSender, err := push.NewFCMSender(context.Background(), cfg.FirebaseServiceAccountKey, fcmTokenProvider)
		if err != nil {
			log.Printf("advertencia: no se pudo inicializar Firebase (%v), las notificaciones push no se enviarán", err)
			pushSender = push.NoopSender{}
		} else {
			pushSender = fcmSender
		}
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
		eventbus.StartBlockchainConfirmedConsumer(cfg.RabbitMQURL, pool, pushSender)
		eventbus.StartSubscriptionEventsConsumer(cfg.RabbitMQURL, pool, pushSender)
		eventbus.StartOrderConfirmedConsumer(cfg.RabbitMQURL, pool, pushSender)
		eventbus.StartOrderCreatedConsumer(cfg.RabbitMQURL, pool, pushSender)
		eventbus.StartOrderShippedConsumer(cfg.RabbitMQURL, pool, pushSender)
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
		usersInfra.BuildSetPublicKeyController(pool),
		usersInfra.BuildGetPublicKeyController(pool),
		usersInfra.BuildAddRolesController(pool),
		cfg.JWTSecret,
	)

	authRouter.RegisterRoutes(
		mux,
		authInfra.BuildLoginController(pool, cfg.JWTSecret, cfg.CookieSecure),
		authInfra.BuildGoogleLoginController(pool, cfg.GoogleOAuthClientIDs, cfg.JWTSecret, cfg.CookieSecure),
	)

	addressesRouter.RegisterRoutes(
		mux,
		addressesInfra.BuildCreateAddressController(pool),
		addressesInfra.BuildListAddressesController(pool),
		addressesInfra.BuildDeleteAddressController(pool),
		addressesInfra.BuildSetDefaultAddressController(pool),
		cfg.JWTSecret,
	)

	assetsRouter.RegisterRoutes(
		mux,
		assetsInfra.BuildCreateAssetController(pool, publisher),
		assetsInfra.BuildGetAllAssetsController(pool),
		assetsInfra.BuildGetMyAssetsController(pool),
		assetsInfra.BuildGetAssetByIdController(pool),
		assetsInfra.BuildUpdateAssetController(pool, publisher),
		assetsInfra.BuildDeleteAssetController(pool),
		assetsInfra.BuildUploadAssetPhotoController(pool, imageUploader),
		assetsInfra.BuildDeleteAssetPhotoController(pool),
		cfg.JWTSecret,
	)

	assetcommentsRouter.RegisterRoutes(
		mux,
		assetcommentsInfra.BuildCreateAssetCommentController(pool, moderationClient),
		assetcommentsInfra.BuildGetAssetCommentsController(pool),
		assetcommentsInfra.BuildDeleteAssetCommentController(pool),
		cfg.JWTSecret,
	)

	businessesRouter.RegisterRoutes(
		mux,
		businessesInfra.BuildCreateBusinessController(pool),
		businessesInfra.BuildGetAllBusinessesController(pool),
		businessesInfra.BuildGetBusinessByIdController(pool),
		businessesInfra.BuildUpdateBusinessController(pool),
		businessesInfra.BuildDeleteBusinessController(pool),
		businessesInfra.BuildUploadBusinessPhotoController(pool, imageUploader),
		businessesInfra.BuildDeleteBusinessPhotoController(pool),
		cfg.JWTSecret,
	)

	businessservicesRouter.RegisterRoutes(
		mux,
		businessservicesInfra.BuildCreateBusinessServiceController(pool),
		businessservicesInfra.BuildListBusinessServicesController(pool),
		businessservicesInfra.BuildUpdateBusinessServiceController(pool),
		businessservicesInfra.BuildDeleteBusinessServiceController(pool),
		cfg.JWTSecret,
	)

	maintenancelogsRouter.RegisterRoutes(
		mux,
		maintenancelogsInfra.BuildCreateMaintenanceLogController(pool, publisher),
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
		postsInfra.BuildLikePostController(pool, pushSender),
		postsInfra.BuildUnlikePostController(pool),
		postsInfra.BuildSavePostController(pool),
		postsInfra.BuildUnsavePostController(pool),
		postsInfra.BuildGetSavedPostsController(pool),
		cfg.JWTSecret,
	)

	commentsRouter.RegisterRoutes(
		mux,
		commentsInfra.BuildCreateCommentController(pool, moderationClient, pushSender),
		commentsInfra.BuildGetCommentsByPostController(pool),
		commentsInfra.BuildDeleteCommentController(pool),
		cfg.JWTSecret,
	)

	reviewsRouter.RegisterRoutes(
		mux,
		reviewsInfra.BuildCreateReviewController(pool, moderationClient, reviewsAdapters.NewHTTPPurchaseVerifier(cfg.PaymentServiceURL)),
		reviewsInfra.BuildGetReviewsByProviderController(pool),
		reviewsInfra.BuildGetProviderRatingController(pool),
		reviewsInfra.BuildGetReviewByIdController(pool),
		reviewsInfra.BuildDeleteReviewController(pool),
		reviewsInfra.BuildLikeReviewController(pool),
		reviewsInfra.BuildUnlikeReviewController(pool),
		cfg.JWTSecret,
	)

	restorerprofilesRouter.RegisterRoutes(
		mux,
		restorerprofilesInfra.BuildUpsertRestorerProfileController(pool),
		restorerprofilesInfra.BuildGetRestorerProfileController(pool),
		restorerprofilesInfra.BuildListRestorerProfilesController(pool),
		cfg.JWTSecret,
	)

	chatRouter.RegisterRoutes(
		mux,
		chatInfra.BuildSendChatMessageController(pool, pushSender),
		chatInfra.BuildGetConversationMessagesController(pool),
		chatInfra.BuildUpdateChatMessageStatusController(pool),
		chatInfra.BuildGetConversationsController(pool),
		chatInfra.BuildDeleteChatMessageController(pool),
		chatInfra.BuildDeleteConversationController(pool),
		cfg.JWTSecret,
	)

	// Recomendaciones (ML): reutiliza cfg.NLPServiceURL porque es el mismo
	// vault-ai-service que ya se usa para moderación de NLP, solo que aquí
	// se llama a /api/v1/ml/recommend-items en vez de /api/v1/nlp/analyze.
	recommendationsRouter.RegisterRoutes(
		mux,
		recommendationsInfra.BuildGetRecommendationsController(cfg.NLPServiceURL),
		cfg.JWTSecret,
	)

	notificationsRouter.RegisterRoutes(
		mux,
		notificationsInfra.BuildCreateNotificationController(pool, pushSender),
		notificationsInfra.BuildGetMyNotificationsController(pool, pushSender),
		notificationsInfra.BuildMarkNotificationAsReadController(pool, pushSender),
		notificationsInfra.BuildMarkAllNotificationsAsReadController(pool, pushSender),
		notificationsInfra.BuildDeleteNotificationController(pool, pushSender),
		cfg.JWTSecret,
	)

	fcmtokensRouter.RegisterRoutes(
		mux,
		fcmtokensInfra.BuildRegisterFCMTokenController(pool),
		fcmtokensInfra.BuildDeleteFCMTokenController(pool),
		cfg.JWTSecret,
	)

	servicerequestsRouter.RegisterRoutes(
		mux,
		servicerequestsInfra.BuildCreateServiceRequestController(pool, pushSender),
		servicerequestsInfra.BuildAcceptServiceRequestController(pool, pushSender),
		servicerequestsInfra.BuildStartServiceRequestController(pool, pushSender),
		servicerequestsInfra.BuildFinishServiceRequestController(pool, pushSender),
		servicerequestsInfra.BuildConfirmServiceRequestController(pool, pushSender),
		servicerequestsInfra.BuildListMyServiceRequestsController(pool),
		servicerequestsInfra.BuildListIncomingServiceRequestsController(pool),
		cfg.JWTSecret,
	)

	handler := middleware.CORS(cfg.CORSOrigin)(middleware.LimitBodySize(mux))

	fmt.Println("API Vault iniciada correctamente.")
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, handler))
}
