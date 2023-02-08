package backend

import (
	"flag"
	// "github.com/Globys031/plotzemis/go/auth"
	// "github.com/Globys031/plotzemis/go/db"
	// "github.com/Globys031/plotzemis/go/routes"
)

var (
	enableTls       = flag.Bool("enable_tls", false, "Use TLS - required for HTTP2.")
	tlsCertFilePath = flag.String("tls_cert_file", "../../misc/localhost.crt", "Path to the CRT/PEM file.")
	tlsKeyFilePath  = flag.String("tls_key_file", "../../misc/localhost.key", "Path to the private key file.")
)

func main() {
	flag.Parse() // parses the above flag variables
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
	// router := routes.RegisterRoutes(authSvc)
	// Addr := fmt.Sprintf("%s:%d", host, port)
	// Addr := fmt.Sprintf(":%d", port)
	// if *enableTls {
	// 	if err := router.RunTLS(Addr, *tlsCertFilePath, *tlsKeyFilePath); err != nil {
	// 		fmt.Errorf("failed starting http2 backend server: %v", err)
	// 	}
	// } else {
	// 	if err := router.Run(Addr); err != nil {
	// 		fmt.Errorf("failed starting http backend server: %v", err)
	// 	}
	// }
}
