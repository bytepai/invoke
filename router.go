package invoke

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// NodeType represents the type of trie node.
type NodeType int

const (
	Static NodeType = iota // Static node type for regular string nodes.
	Param                  // Param node type for parameter nodes, e.g., /user/:id.
	Regex                  // Regex node type for regex pattern nodes, e.g., /product/{regexp}.
)

// TrieNode represents a node in the trie.
type TrieNode struct {
	Children []*TrieNode            `json:"children"`  // Child nodes of the current node.
	Handler  func(ctx *HttpContext) `json:"-"`         // Handler function for the node.
	Level    int                    `json:"level"`     // Depth level of the node in the trie.
	Pattern  string                 `json:"pattern"`   // Pattern of the node.
	NodeType NodeType               `json:"node_type"` // Type of the node (Static, Param, Regex).
	Method   string                 `json:"method"`    // HTTP method associated with the route.
	FullPath string                 `json:"full_path"` // Full path to the node.
	Path     string                 `json:"path"`      // Name of the current node.
}

// Router represents a trie-based router.
type router struct {
	Root            *TrieNode                               `json:"root"`  // Root node of the trie.
	Param           map[string]string                       `json:"param"` // Map of route parameters.
	BeforeHooks     []func(ctx *HttpContext) bool           `json:"-"`     // Global before hooks.
	AfterHooks      []func(ctx *HttpContext)                `json:"-"`     // Global after hooks.
	NotFound        func(ctx *HttpContext)                  `json:"-"`     // Handler for 404 Not Found.
	Prefix          string                                  // Prefix for the routes in the group.
	GroupBefore     []func(ctx *HttpContext) bool           // Group-specific before hooks.
	GroupAfter      []func(ctx *HttpContext)                // Group-specific after hooks.
	RecoveryHandler func(ctx *HttpContext, err interface{}) // Custom recovery handler
	Assets          func(ctx *HttpContext) bool             `json:"-"` // Handler for serving static files.
}

var Router = NewRouter()

// NewRouter creates a new router with an empty root node.
func NewRouter() *router {
	return &router{
		Root: &TrieNode{
			Children: []*TrieNode{},
			Level:    0,
			Pattern:  "",
			NodeType: Static,
			FullPath: "",
		},
		Param:    make(map[string]string),
		NotFound: defaultNotFoundHandler,
		Assets:   defaultAssetsHandler,
	}

}

// AddRoute adds a route to the router.
func (r *router) AddRoute(method, path string, handler func(ctx *HttpContext)) {
	parts := splitPath(path) // Split the path into parts.
	curr := r.Root           // Start from the root node.

	for _, part := range parts {
		nodeType, pattern, paramName := getNodeTypeAndPattern(part) // Determine node type and pattern.
		found := false

		for _, child := range curr.Children {
			if child.Pattern == pattern && child.NodeType == nodeType && child.Method == method {
				curr = child // Move to the matching child node.
				found = true
				break
			}
		}

		if !found {
			newNode := &TrieNode{
				Children: []*TrieNode{},
				Level:    curr.Level + 1,
				Pattern:  pattern,
				NodeType: nodeType,
				FullPath: curr.FullPath + "/" + pattern,
				Method:   method, // Store the HTTP method.
			}
			if nodeType == Regex {
				newNode.Pattern = paramName + ":" + pattern
			}
			curr.Children = append(curr.Children, newNode) // Add new node if not found.
			curr = newNode
		}
	}
	if curr.Handler != nil {
		info := fmt.Sprintf("Warning: Route '%s' with method '%s' is already registered.\n", path, method)
		panic(info)
	}
	curr.Handler = handler // Assign the handler to the leaf node.
}

// ServeHTTP handles HTTP requests.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			if r.RecoveryHandler == nil {
				// Log the error and return a 500 Internal Server Error response
				http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
			} else {
				// Use the custom recovery handler if provided
				r.RecoveryHandler(&HttpContext{W: w, Req: req}, err)
			}
		}
	}()

	path := strings.ToLower(req.URL.Path)
	parts := splitPath(path) // Split the request path into parts.
	curr := r.Root           // Start from the root node.
	params := make(map[string]string)
	method := req.Method

	// Create HttpContext
	ctx := &HttpContext{
		W:      w,
		Req:    req,
		Params: params,
	}

	// Execute global before hooks
	for _, hook := range r.BeforeHooks {
		if !hook(ctx) {
			return
		}
	}

	// Execute group before hooks
	for _, hook := range r.GroupBefore {
		if !hook(ctx) {
			return
		}
	}

	for _, part := range parts {
		found := false
		for _, child := range curr.Children {
			switch child.NodeType {
			case Static:
				if child.Method == method && child.Pattern == part {
					curr = child
					found = true
				}
			case Param:
				if child.Method == method {
					curr = child
					params[child.Pattern] = part // Add param to the map.
					found = true
				}
			case Regex:
				//if child.Method == method {
				patternParts := strings.SplitN(child.Pattern, ":", 2)
				if child.Method == method && len(patternParts) == 2 {
					paramName, regexPattern := patternParts[0], patternParts[1]
					if match := regexp.MustCompile(regexPattern).FindString(part); match == part {
						curr = child
						params[paramName] = match // Add param to the map.
						found = true
					}
				}
				//}
			}
			if found {
				break
			}
		}
		if !found {
			if !r.Assets(ctx) {
				return
			}
			r.NotFound(ctx) // Handle 404 Not Found.
			return
		}
	}

	r.Param = params
	req = req.WithContext(contextWithParams(req.Context(), params)) // Add params to context.

	if curr.Handler != nil {
		curr.Handler(ctx)
	} else {
		r.NotFound(ctx)
	}

	// Execute global after hooks
	for _, hook := range r.AfterHooks {
		hook(ctx)
	}

	// Execute group after hooks
	for _, hook := range r.GroupAfter {
		hook(ctx)
	}
}

// SetRecoveryHandler sets the custom recovery handler.
func (r *router) SetRecoveryHandler(handler func(ctx *HttpContext, err interface{})) {
	r.RecoveryHandler = handler
}

// RegisterBeforeHook registers a global before hook.
func (r *router) RegisterBeforeHook(hook func(ctx *HttpContext) bool) {
	r.BeforeHooks = append(r.BeforeHooks, hook)
}

// RegisterAfterHook registers a global after hook.
func (r *router) RegisterAfterHook(hook func(ctx *HttpContext)) {
	r.AfterHooks = append(r.AfterHooks, hook)
}

// SetNotFoundHandler sets the 404 Not Found handler.
func (r *router) SetNotFoundHandler(handler func(ctx *HttpContext)) {
	r.NotFound = handler
}

// SetAssetsHandler sets the handler for serving static files.
func (r *router) SetAssetsHandler(handler func(ctx *HttpContext) bool) {
	r.Assets = handler
}

// defaultAssetsHandler is the default handler for serving static files.
func defaultAssetsHandler(ctx *HttpContext) bool {
	filePath := "./" + strings.Trim(ctx.Req.URL.Path, "/")

	// If the requested path is a directory, try to serve index.html in that directory.
	if fileInfo, err := os.Stat(filePath); err == nil && fileInfo.IsDir() {
		indexFilePath := filepath.Join(filePath, "index.html")
		if _, err := os.Stat(indexFilePath); err == nil {
			filePath = indexFilePath
		}
	}

	// Check if the file exists.
	if fileExists(filePath) {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return true
		}
		if fileInfo.IsDir() {
			return true
		}

		http.ServeFile(ctx.W, ctx.Req, filePath)
		return false // Return false to allow further routing.
	}
	return true // Return true to indicate file not found.
}

// fileExists checks if a file exists.
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// Group creates a new router group with the specified prefix.
func (r *router) Group(prefix string) *router {
	return &router{
		Root:        r.Root,
		Param:       r.Param,
		BeforeHooks: r.BeforeHooks,
		AfterHooks:  r.AfterHooks,
		NotFound:    r.NotFound,
		Prefix:      r.Prefix + prefix,
		GroupBefore: append([]func(ctx *HttpContext) bool{}, r.GroupBefore...), // Copy hooks from parent group.
		GroupAfter:  append([]func(ctx *HttpContext){}, r.GroupAfter...),       // Copy hooks from parent group.
	}
}

// RegisterGroupBeforeHook registers a before hook for the group.
func (r *router) RegisterGroupBeforeHook(hook func(ctx *HttpContext) bool) {
	r.GroupBefore = append(r.GroupBefore, hook)
}

// RegisterGroupAfterHook registers an after hook for the group.
func (r *router) RegisterGroupAfterHook(hook func(ctx *HttpContext)) {
	r.GroupAfter = append(r.GroupAfter, hook)
}

// registerRoute registers a route.
func (r *router) registerRoute(method, path string, handler func(ctx *HttpContext)) {
	fullPath := r.Prefix + path
	fullPath = strings.ToLower(fullPath)
	r.AddRoute(method, fullPath, func(ctx *HttpContext) {
		// Execute group before hooks
		for _, hook := range r.GroupBefore {
			if !hook(ctx) {
				return
			}
		}

		handler(ctx)

		// Execute group after hooks
		for _, hook := range r.GroupAfter {
			hook(ctx)
		}
	})
}

// GET registers a GET route.
func (r *router) GET(path string, handler func(ctx *HttpContext)) {
	r.registerRoute("GET", path, handler)
}

// POST registers a POST route.
func (r *router) POST(path string, handler func(ctx *HttpContext)) {
	r.registerRoute("POST", path, handler)
}

// DELETE registers a DELETE route.
func (r *router) DELETE(path string, handler func(ctx *HttpContext)) {
	r.registerRoute("DELETE", path, handler)
}

// PUT registers a PUT route.
func (r *router) PUT(path string, handler func(ctx *HttpContext)) {
	r.registerRoute("PUT", path, handler)
}

// PATCH registers a PATCH route.
func (r *router) PATCH(path string, handler func(ctx *HttpContext)) {
	r.registerRoute("PATCH", path, handler)
}

// HEAD registers a HEAD route.
func (r *router) HEAD(path string, handler func(ctx *HttpContext)) {
	r.registerRoute("HEAD", path, handler)
}

// OPTIONS registers a OPTIONS route.
func (r *router) OPTIONS(path string, handler func(ctx *HttpContext)) {
	r.registerRoute("OPTIONS", path, handler)
}

// defaultNotFoundHandler is the default 404 Not Found handler.
func defaultNotFoundHandler(ctx *HttpContext) {
	http.Error(ctx.W, "404 - Not Found", http.StatusNotFound)
}

// splitPath splits the path into parts.
func splitPath(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

// getNodeTypeAndPattern determines the node type and pattern.
func getNodeTypeAndPattern(part string) (NodeType, string, string) {
	if strings.HasPrefix(part, ":") {
		return Param, part[1:], ""
	}
	if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
		// Extract the parameter name and regex pattern
		content := part[1 : len(part)-1]
		patternParts := strings.SplitN(content, ":", 2)
		if len(patternParts) < 2 {
			return Static, part, ""
		}
		return Regex, patternParts[1], patternParts[0]
	}
	return Static, part, ""
}

// contextKey is a type for context keys to avoid conflicts.
type contextKey string

// ParamsKey is the context key for route parameters.
const ParamsKey = contextKey("params")

// contextWithParams returns a new context with the given parameters.
func contextWithParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, ParamsKey, params)
}

// GetParams retrieves the route parameters from the request context.
func GetParams(req *http.Request) map[string]string {
	if params, ok := req.Context().Value(ParamsKey).(map[string]string); ok {
		return params
	}
	return nil
}

// ListenAndServe starts an HTTP server with the provided address and handler.
func (r *router) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, r)
}

// StartServe starts an HTTP server with the provided address and handler.
func (r *router) StartServe(addr string) error {
	return http.ListenAndServe(addr, r)
}
