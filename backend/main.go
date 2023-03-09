package main

import (
	"flag"
	"fmt"

	// "github.com/Globys031/PostgreScrutiniser/backend/cmd"
	// "github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/Globys031/PostgreScrutiniser/backend/web"
)

var (
	enableTls       = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
	tlsCertFilePath = flag.String("tls_cert_file", "../../misc/localhost.crt", "Path to the CRT/PEM file.")
	tlsKeyFilePath  = flag.String("tls_key_file", "../../misc/localhost.key", "Path to the private key file.")
)

func main() {
	flag.Parse() // parses the above flag variables

	// ////////////////////////
	// // Initialise logging
	logger := utils.InitLogging()

	// ////////////////////////
	// // Run the actual application
	// cmd.RunChecks(logger)

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
	router := web.RegisterRoutes(logger)

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
