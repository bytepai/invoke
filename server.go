package invoke

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	middlewares = make(map[string]func(http.Handler) http.Handler)
	servers     []ServerConfig
	serverConfigPath = "server_conf.json"
	serverConfigLock sync.Mutex

)

type TLSConfig struct {
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type RateLimitConfig struct {
	RequestsPerSecond int `json:"requests_per_second"`
}

type LoggingConfig struct {
	AccessLog string `json:"access_log"`
	ErrorLog  string `json:"error_log"`
	LogLevel  string `json:"log_level"`
}

type CORSConfig struct {
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
}

type SecurityConfig struct {
	AllowedHosts   []string   `json:"allowed_hosts"`
	CORS           CORSConfig `json:"cors"`
	CSRFProtection bool       `json:"csrf_protection"`
}

type TimeoutsConfig struct {
	IdleTimeout           time.Duration `json:"idle_timeout"`
	HeaderTimeout         time.Duration `json:"header_timeout"`
	ResponseHeaderTimeout time.Duration `json:"response_header_timeout"`
}

type KeepAliveConfig struct {
	Enabled bool          `json:"enabled"`
	Timeout time.Duration `json:"timeout"`
}

type CompressionConfig struct {
	EnableGzip       bool `json:"enable_gzip"`
	CompressionLevel int  `json:"compression_level"`
}

type StaticFilesConfig struct {
	StaticDir string `json:"static_dir"`
	IndexFile string `json:"index_file"`
}

type ServerConfig struct {
	Domain         string            `json:"domain"`
	Port           int               `json:"port"`
	ReadTimeout    time.Duration     `json:"read_timeout"`
	WriteTimeout   time.Duration     `json:"write_timeout"`
	MaxHeaderBytes int               `json:"max_header_bytes"`
	TLS            TLSConfig         `json:"tls"`
	Limits         RateLimitConfig   `json:"limits"`
	RateLimit      RateLimitConfig   `json:"rate_limit"`
	Logging        LoggingConfig     `json:"logging"`
	Security       SecurityConfig    `json:"security"`
	Timeouts       TimeoutsConfig    `json:"timeouts"`
	KeepAlive      KeepAliveConfig   `json:"keep_alive"`
	Compression    CompressionConfig `json:"compression"`
	StaticFiles    StaticFilesConfig `json:"static_files"`
	Middleware     []string          `json:"middleware"`
}

// Middleware registration function
func RegisterMiddleware(name string, middleware func(http.Handler) http.Handler) {
	middlewares[name] = middleware
}

// Server registration function with configuration synchronization
func RegisterServer(config ServerConfig) {
	serverConfigLock.Lock()
	defer serverConfigLock.Unlock()

	servers = append(servers, config)
	updateConfigFile()
}

func updateConfigFile() {
	config := Config{Servers: servers}
	file, err := os.Create(serverConfigPath)
	if err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(config); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}
}

type Config struct {
	Servers []ServerConfig `json:"servers"`
}

// LoadConfig loads configuration from a JSON file
func loadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// createDefaultConfig Create default configuration file
func createDefaultConfig(filePath string) error {
	defaultConfig := Config{
		Servers: []ServerConfig{
			{
				Domain:         "localhost",
				Port:           8080,
				ReadTimeout:    5 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1048576,
				TLS:            TLSConfig{},
				Limits:         RateLimitConfig{RequestsPerSecond: 100},
				RateLimit:      RateLimitConfig{RequestsPerSecond: 100},
				Logging:        LoggingConfig{LogLevel: "info"},
				Security:       SecurityConfig{CSRFProtection: true},
				Timeouts:       TimeoutsConfig{IdleTimeout: 120 * time.Second},
				KeepAlive:      KeepAliveConfig{Enabled: true, Timeout: 30 * time.Second},
				Compression:    CompressionConfig{EnableGzip: true, CompressionLevel: 5},
				StaticFiles:    StaticFilesConfig{StaticDir: "./static", IndexFile: "index.html"},
				Middleware:     []string{"logging", "rateLimiting"},
			},
		},
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(defaultConfig); err != nil {
		return err
	}
	return nil
}

// InitializeConfig checks if the config file exists, creates it with default values if not
func init() {
	if _, err := os.Stat(serverConfigPath); os.IsNotExist(err) {
		log.Println("Config file not found, creating default config")
		if err := createDefaultConfig(serverConfigPath); err != nil {
			log.Fatalf("Failed to create default config: %v", err)
		}
	} else {
		config, err := loadConfig(serverConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		servers = config.Servers
	}
	RegisterMiddleware("logging", loggingDefaultMiddleware)
	RegisterMiddleware("rateLimiting", rateLimitingDefaultMiddleware)
}

// Server represents a single server instance
// type Server struct {
// 	Config ServerConfig
// 	Router *Router
// }

// // MultiServer handles multiple servers and graceful shutdown
// type MultiServer struct {
// 	servers []*Server
// }

// // AddServer adds a new server to the MultiServer
// func (ms *MultiServer) AddServer(config ServerConfig, router *Router) {
// 	ms.servers = append(ms.servers, &Server{
// 		Config: config,
// 		Router: router,
// 	})
// }

// startServer Add Middleware and TLS Support (Optional)
func startServer(config ServerConfig) {
	mux := http.NewServeMux()

	if config.StaticFiles.StaticDir != "" {
		fileServer := http.FileServer(http.Dir(config.StaticFiles.StaticDir))
		mux.Handle("/", http.StripPrefix("/", fileServer))
	}

	handler := applyMiddlewares(mux, config.Middleware)

	srv := &http.Server{
		Addr:           config.Domain + ":" + string(config.Port),
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	log.Printf("Starting server on %s:%d", config.Domain, config.Port)
	if config.TLS.CertFile != "" && config.TLS.KeyFile != "" {
		if err := srv.ListenAndServeTLS(config.TLS.CertFile, config.TLS.KeyFile); err != nil {
			log.Fatalf("HTTPS server failed: %v", err)
		}
	} else {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}
}
func applyMiddlewares(handler http.Handler, middleware []string) http.Handler {
	for _, m := range middleware {
		if mw, ok := middlewares[m]; ok {
			handler = mw(handler)
		}
	}
	return handler
}

func loggingDefaultMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func rateLimitingDefaultMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement rate limiting logic here
		next.ServeHTTP(w, r)
	})
}

func StartServer(r ...*router) {
	var wg sync.WaitGroup
	for _, srv := range servers {
		wg.Add(1)
		go func(s ServerConfig) {
			defer wg.Done()
			addr := s.Domain + ":" + strconv.Itoa(s.Port)
			httpServer := &http.Server{
				Addr:           addr,
				Handler:        Router,
				ReadTimeout:    s.ReadTimeout * time.Second,
				WriteTimeout:   s.WriteTimeout * time.Second,
				MaxHeaderBytes: s.MaxHeaderBytes,
			}
			go func() {
				fmt.Printf("Server listening: %s\n", addr)
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			<-waitForShutdown()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(ctx); err != nil {
				panic(err)
			}
		}(srv)
	}
	wg.Wait()
}

// StartAll starts all servers and listens for shutdown signals
// func (ms *MultiServer) StartAll() {
// 	var wg sync.WaitGroup
// 	for _, srv := range ms.servers {
// 		wg.Add(1)
// 		go func(s *Server) {
// 			defer wg.Done()
// 			httpServer := &http.Server{
// 				Addr:           s.Config.Domain + ":" + string(s.Config.Port),
// 				Handler:        s.Router,
// 				ReadTimeout:    s.Config.ReadTimeout * time.Second,
// 				WriteTimeout:   s.Config.WriteTimeout * time.Second,
// 				MaxHeaderBytes: s.Config.MaxHeaderBytes,
// 			}
// 			go func() {
// 				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 					panic(err)
// 				}
// 			}()
// 			<-waitForShutdown()
// 			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 			defer cancel()
// 			if err := httpServer.Shutdown(ctx); err != nil {
// 				panic(err)
// 			}
// 		}(srv)
// 	}
// 	wg.Wait()
// }

// waitForShutdown listens for interrupt signals and returns a channel
func waitForShutdown() <-chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	return shutdown
}
