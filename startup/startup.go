package startup

import (
	"net/http"

	"github.com/tanyoAH/tanyo-server/config"
	"github.com/tanyoAH/tanyo-server/controllers"
	// Import xwsenv for it's init routine that's relied upon for startup
	"github.com/gorilla/handlers"
	"github.com/tanyoAH/tanyo-server/models" /* MySQL */
	_ "github.com/tanyoAH/tanyo-server/twsproto"
	"os"

	"github.com/tanyoAH/tanyo-server/tws"
)

var Log = config.Conf.GetLogger()

var CorsHandler = handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}), handlers.AllowCredentials(), handlers.AllowedHeaders([]string{"content-type"}), handlers.AllowedOrigins([]string{config.Conf.Home, config.Conf.Origin}))

func StartServer() {
	Log.Info("Server Loaded")

	err := models.Setup()
	if err != nil {
		Log.WithField("error", err).Fatal("Couldn't setup models")
	}

	Log.Info("Initializing WS environment(s)")
	ws, err := tws.CreateAndSetupNewState()
	if err != nil {
		Log.WithField("error", err).Fatal("Couldn't setup WS Eevironment(s)")
	}
	go startTWSServer(ws)

	Log.WithField("hostname", config.Conf.ApiUrl).Info("API Server starting")
	if config.Conf.IsDebugMode() {
		http.ListenAndServe(config.Conf.ApiUrl, CorsHandler(handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateRouter(ws))))
	} else {
		http.ListenAndServe(config.Conf.ApiUrl, handlers.CombinedLoggingHandler(os.Stdout, controllers.CreateRouter(ws)))
	}
}

func startTWSServer(w *tws.State) {
	Log.WithField("hostname", config.Conf.WsUrl).Info("TWS Server starting")
	if config.Conf.IsDebugMode() {
		http.ListenAndServe(config.Conf.WsUrl, CorsHandler(handlers.CombinedLoggingHandler(os.Stdout, w.CreateRouter())))
	} else {
		http.ListenAndServe(config.Conf.WsUrl, handlers.CombinedLoggingHandler(os.Stdout, w.CreateRouter()))
	}
}
