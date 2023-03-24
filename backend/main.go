package main

import (
	"flag"
	"fmt"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/Globys031/PostgreScrutiniser/backend/web"
	"github.com/Globys031/PostgreScrutiniser/backend/web/auth"
)

var (
	enableTls       = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
	tlsCertFilePath = flag.String("tls_cert_file", "/usr/local/postgrescrutiniser/confs/scrutiniser.crt", "Path to the CRT/PEM file.")
	tlsKeyFilePath  = flag.String("tls_key_file", "/usr/local/postgrescrutiniser/confs/scrutiniser.key", "Path to the private key file.")
	appUsername     = "postgrescrutiniser"
	hostname        = "localhost"
	backupDir       = "/usr/local/postgrescrutiniser/backups"
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
		logger.LogError(fmt.Errorf("Could not find our main application user:  %v", err))
		return
	}

	_, postgrePort, _, postgreUsername, password, _ := utils.ParsePgpassFile(appUser, logger)
	postgreUid, postgreGid, err := utils.GetUserIds(postgreUsername, logger)
	postgresUser := &utils.User{
		Username: postgreUsername,
		Uid:      postgreUid,
		Gid:      postgreGid,
	}
	if err != nil {
		logger.LogError(fmt.Errorf("Could not find main PostgreSql user:  %v", err))
		return
	}

	////////////////////////
	// Initialise database connection
	dbHandler, _ := utils.InitDbConnection(hostname, postgresUser.Username, password, postgrePort, logger)

	//////////////////////////
	// Loads configs
	config, _ := LoadConfig(logger)
	appPort := config.Backend_port

	jwt := &auth.JwtWrapper{
		SecretKey:       config.JWT_secret_key,
		Issuer:          "postgre-scrutiniser",
		ExpirationHours: 4, // token expires after 4 hours
	}
	//////////////////////////

	//////////////////////////
	// Initialise webserver and routes
	router := web.RegisterRoutes(jwt, dbHandler, appUser, postgresUser, backupDir, logger)

	// router := web.RegisterRoutes(authSvc)
	Addr := fmt.Sprintf(":%d", appPort)
	if *enableTls {
		if err := router.RunTLS(Addr, *tlsCertFilePath, *tlsKeyFilePath); err != nil {
			logger.LogFatal(fmt.Errorf("failed starting https backend server: %v", err))
		}
	} else {
		if err := router.Run(Addr); err != nil {
			logger.LogFatal(fmt.Errorf("failed starting http backend server: %v", err))
		}
	}
}
