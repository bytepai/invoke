# Invoke

Invoke is a lightweight trie-based HTTP router for Go. It supports static routes, parameterized routes, regex routes, and middleware hooks for before and after route handling. This router is designed to be simple and efficient, making it ideal for small to medium-sized web applications.

> **Note**: This project was 100% generated by ChatGPT-4.

## Features

- **Static Routes**: Define static routes with fixed paths.
- **Parameterized Routes**: Define dynamic routes with parameters, e.g., `/user/:id`.
- **Regex Routes**: Define routes with regex patterns, e.g., `/product/{id:[0-9]+}`.
- **Middleware Hooks**: Support for global and group-specific before and after hooks.
- **404 Not Found Handler**: Customizable handler for 404 errors.
- **Static File Serving**: Built-in support for serving static files.
- **Custom Recovery handler**:An optional custom handler for panics, allowing for specific error handling logic.

## Installation

```bash
go get github.com/bytepai/invoke

```
## Usage
### Basic Example
```
package main

import (
	"fmt"
	. "github.com/bytepai/invoke"
)

func main() {
	r := NewRouter()

	// Set custom 404 handler
	r.SetNotFoundHandler(func(ctx *HttpContext) {
		ctx.WriteString("Custom 404 - Not Found")
	})

	r.GET("/", func(ctx *HttpContext) {
		ctx.WriteString("index")
	})

	r.GET("/hello", func(ctx *HttpContext) {
		ctx.WriteString("Hello, world!")
	})

	r.GET("/user/:name", func(ctx *HttpContext) {
		name := ctx.Params["name"]
		ctx.WriteString(fmt.Sprintf("Hello, %s!", name))
	})

	r.GET("/order/:id", func(ctx *HttpContext) {
		id := ctx.Params["id"]
		ctx.WriteString(fmt.Sprintf("Order ID: %s", id))
	})
	r.POST("/product/{id:[0-9]+}", func(ctx *HttpContext) {
		id := ctx.Params["id"]
		ctx.WriteSuccessJSON(id)
	})
	r.GET("/phone/{phone:1[3456789]\\d{9}}", func(ctx *HttpContext) {
		phone := ctx.Params["phone"]
		ctx.WriteString(fmt.Sprintf("Hello, %s!", phone))
	})
	fmt.Println("Server listening on port 8080...")
	r.ListenAndServe(":8080")
}
```
### Grouping Routes
```
package main

import (
	"fmt"
	. "github.com/bytepai/invoke"
)

func main() {
	r := NewRouter()

	// Create a new group
	apiGroup := r.Group("/api")

	// Register group before hook
	apiGroup.RegisterGroupBeforeHook(func(ctx *HttpContext) bool {
		fmt.Println("API group before hook")
		return true // Return false to terminate
	})

	// Register group after hook
	apiGroup.RegisterGroupAfterHook(func(ctx *HttpContext) {
		fmt.Println("API group after hook")
	})

	// Add routes to the group
	apiGroup.GET("/info", func(ctx *HttpContext) {
		ctx.WriteString("API Info")
	})

	fmt.Println("Server listening on port 8080...")
	r.ListenAndServe(":8080")
}
```

### Middleware Hooks
```
package main

import (
	"fmt"
	. "github.com/bytepai/invoke"
)

func main() {
	r := NewRouter()

	// Register global before hook
	r.RegisterBeforeHook(func(ctx *HttpContext) bool {
		fmt.Println("Global before hook")
		return true
	})

	// Register global after hook
	r.RegisterAfterHook(func(ctx *HttpContext) {
		fmt.Println("Global after hook")
	})

	r.SetRecoveryHandler(func(ctx *HttpContext, err interface{}) {
		errString := fmt.Sprintf("Custom recovery handler: %v\n", err)
		ctx.WriteString(errString)
	})
	
	r.POST("/panic", func(ctx *HttpContext) {
		var s []string
		fmt.Println(s[1])
	})

	fmt.Println("Server listening on port 8080...")
	r.ListenAndServe(":8080")
}
```
### Serve
#### Router->Serve(config)
```
	Router.StartServer(":8080") 
	Router.StartServer("localhost:8080")
	Router.StartServerConfig(ServerConfig)
	
```
#### Serve(config,router)
```
    Server.Start("localhost:8080",Router)
	Server.StartConfig(ServerConfig,Router)

```




