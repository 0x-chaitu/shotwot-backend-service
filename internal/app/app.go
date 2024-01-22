package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"shotwot_backend/internal/config"
	delivery "shotwot_backend/internal/delivery/http"
	"shotwot_backend/internal/repository"
	"shotwot_backend/internal/server"
	"shotwot_backend/internal/service"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/database/mongodb"
	"shotwot_backend/pkg/firebase"
	"shotwot_backend/pkg/logger"
	"shotwot_backend/pkg/s3"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)

	if err != nil {
		logger.Error(err)

		return
	}

	// sdkConfig, err := awsconfig.LoadDefaultConfig(context.Background(),
	// 	awsconfig.WithRegion("eu-central-1"),
	// 	awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("CYH4MC8ZAIQCNKXS0I9C", "T5ZBEVRccXndzP0wgkGPm2L1VcLHeeJVZr9WvJLD", "TOKEN")),
	// )
	// if err != nil {
	// 	fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
	// 	fmt.Println(err)
	// }

	// client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
	// 	o.BaseEndpoint = aws.String("https://s3.eu-central-1.wasabisys.com/")
	// })

	// x := s3.NewPresignClient(client)
	// bucketName := "warehouse-test"
	// objectKey := "warehouse.mp4"
	// request, err := x.PresignPutObject(context.TODO(), &s3.PutObjectInput{
	// 	Bucket:          aws.String(bucketName),
	// 	Key:             aws.String(objectKey),
	// 	ContentType:     aws.String("video/mp4"),
	// 	ContentEncoding: aws.String("base64"),
	// }, func(opts *s3.PresignOptions) {
	// 	opts.Expires = time.Duration(100000 * int64(time.Second))
	// })
	// if err != nil {
	// 	log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
	// 		bucketName, objectKey, err)
	// }
	// logger.Info(request)

	wasabiS3Client, err := s3.NewWasabiBucket("")
	if err != nil {
		logger.Error(err)
		return
	}

	mongoClient, err := mongodb.NewClient("mongodb+srv://goshotwot:Tokyo.3619@shotwottest.rpgjdmf.mongodb.net/?retryWrites=true&w=majority")
	if err != nil {
		logger.Error(err)
		return
	}

	db := mongoClient.Database("shotwot")

	repos := repository.NewRepositories(db)
	tokenManager, err := jwtauth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}
	adminTokenManager, err := jwtauth.NewAdminManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}
	authClient, err := firebase.NewAuthClient()
	if err != nil {
		logger.Error(err)
		return
	}
	services := service.NewServices(
		service.Deps{
			Repos:             repos,
			TokenManager:      tokenManager,
			AdminTokenManager: adminTokenManager,
			AccessTokenTTL:    cfg.Auth.JWT.AccessTokenTTL,
			RefreshTokenTTL:   cfg.Auth.JWT.RefreshTokenTTL,
			AuthClient:        authClient,
			WasabiS3Client:    wasabiS3Client,
		})
	handlers := delivery.NewHandler(services)

	srv := server.NewServer(cfg, handlers.Init(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
