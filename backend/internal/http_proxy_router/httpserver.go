package http_proxy_router

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	log.Printf("DEBUG: Starting proxy server with addr=%s", viper.GetString("http.addr"))
	//gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	r := InitRouter()
	HttpSrvHandler = &http.Server{
		Addr:           viper.GetString("http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(viper.GetInt("http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(viper.GetInt("http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(viper.GetInt("http.max_header_bytes")),
	}
	log.Printf(" [INFO] http_proxy_run %s\n", viper.GetString("http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] http_proxy_run %s err:%v\n", viper.GetString("http.addr"), err)
	}
}

// 启动https服务
func HttpsServerRun() {
	//gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	//r := InitRouter(middleware.RecoveryMiddleware(),
	//	middleware.RequestLog())
	r := InitRouter()
	HttpsSrvHandler = &http.Server{
		Addr:           viper.GetString("https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(viper.GetInt("https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(viper.GetInt("https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(viper.GetInt("https.max_header_bytes")),
	}
	log.Printf(" [INFO] https_proxy_run %s\n", viper.GetString("https.addr"))
	//todo 以下命令只在编译机有效，如果是交叉编译情况下需要单独设置路径
	//if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil && err!=http.ErrServerClosed {
	if err := HttpsSrvHandler.ListenAndServeTLS("./cert_file/server.crt", "./cert_file/server.key"); err != nil && err != http.ErrServerClosed {
		log.Fatalf(" [ERROR] https_proxy_run %s err:%v\n", viper.GetString("https.addr"), err)
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Printf(" [ERROR] http_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] http_proxy_stop %v stopped\n", viper.GetString("http.addr"))
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] https_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] https_proxy_stop %v stopped\n", viper.GetString("proxy.https.addr"))
}
