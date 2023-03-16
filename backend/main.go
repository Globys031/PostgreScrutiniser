package main

import (
	"flag"
	"fmt"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/Globys031/PostgreScrutiniser/backend/web"
)

var (
	enableTls       = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
	tlsCertFilePath = flag.String("tls_cert_file", "../../misc/localhost.crt", "Path to the CRT/PEM file.")
	tlsKeyFilePath  = flag.String("tls_key_file", "../../misc/localhost.key", "Path to the private key file.")
	appUsername     = "postgrescrutiniser"
	hostname        = "localhost"
)

func main() {
	flag.Parse() // parses the above flag variables

	////////////////////////
	// Initialise logging
	logger := utils.InitLogging()

	////////////////////////
	// Save main postgres user and our app's user info
	appUserUid, appUserGid, err := utils.GetUserIds(appUsername, logger)
	appUser := &utils.User{
		Username: appUsername,
		Uid:      appUserUid,
		Gid:      appUserGid,
	}
	if err != nil {
		logger.LogError("Could not find our main application user: " + err.Error())
		return
	}

	_, port, _, postgreUsername, password, _ := utils.ParsePgpassFile(appUser, logger)
	postgreUid, postgreGid, err := utils.GetUserIds(postgreUsername, logger)
	postgresUser := &utils.User{
		Username: postgreUsername,
		Uid:      postgreUid,
		Gid:      postgreGid,
	}
	if err != nil {
		logger.LogError("Could not find main PostgreSql user: " + err.Error())
		return
	}

	////////////////////////
	// Initialise database connection
	dbHandler, _ := utils.InitDbConnection(hostname, postgresUser.Username, password, port, logger)

	//////////////////////////
	// Loads configs
	// config, err := LoadConfig()
	// if err != nil {
	// 	log.Fatalln("Failed at config", err)
	// }
	// port := config.Backend_port
	// host := config.Backend_host
	//////////////////////////

	// jwt := auth.JwtWrapper{
	// 	SecretKey:       config.JWT_secret_key,
	// 	Issuer:          "go-grpc-auth-svc",
	// 	ExpirationHours: 24 * 365,
	// }
	// authSvc := &routes.AuthService{
	// 	// Handler: db_handler,
	// 	Jwt: jwt,
	// }

	// //////////////////////////
	// // Initialise webserver and routes
	router := web.RegisterRoutes(dbHandler, appUser, postgresUser, logger)

	// router := web.RegisterRoutes(authSvc)
	Addr := ":8080"
	if *enableTls {
		if err := router.RunTLS(Addr, *tlsCertFilePath, *tlsKeyFilePath); err != nil {
			fmt.Errorf("failed starting http2 backend server: %v", err)
		}
	} else {
		if err := router.Run(Addr); err != nil {
			fmt.Errorf("failed starting http backend server: %v", err)
		}
	}
}
