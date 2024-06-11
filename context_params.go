package invoke

import (
	"errors"
	"mime/multipart"
	"net/http"
	"time"
)

const MaxMultipartBytes = 32 << 20 // 32 MB

// FormDataHandler is a callback function type for handling form data key-value pairs.
type FormDataHandler func(key string, value string)

// FileHandler is a callback function type for handling file data.
type FileHandler func(key string, file multipart.File, fileHeader *multipart.FileHeader)

// parmQuery is a helper function to retrieve a parameter value from various sources, including URL parameters, form data, and multipart form data.
func (ctx *HttpContext) parmQuery(key string) string {
	if v, ok := ctx.Params[key]; ok {
		return v
	}
	// Parse form data if not already parsed
	if ctx.Req.Form == nil {
		ctx.Req.ParseForm()
	}

	// Retrieve value from form data
	if formValue := ctx.Req.Form.Get(key); formValue != "" {
		return formValue
	}

	// Parse multipart form data if not already parsed
	if ctx.Req.MultipartForm == nil {
		if err := ctx.Req.ParseMultipartForm(MaxMultipartBytes); err != nil && err != http.ErrNotMultipart {
			return ""
		}
	}

	// Retrieve value from multipart form data
	if ctx.Req.MultipartForm != nil {
		if multipartValue, ok := ctx.Req.MultipartForm.Value[key]; ok && len(multipartValue) > 0 {
			return multipartValue[0]
		}
	}
	return ""
}

// handleFormString iterates over all form data (including multipart form data) and applies the handler function.
func (ctx *HttpContext) handleFormString(handler FormDataHandler) {
	// Parse standard form data if not already parsed
	if ctx.Req.Form == nil {
		ctx.Req.ParseForm()
	}

	// Iterate over standard form data
	for key, values := range ctx.Req.Form {
		for _, value := range values {
			handler(key, value)
		}
	}

	// Parse multipart form data if not already parsed
	if ctx.Req.MultipartForm == nil {
		ctx.Req.ParseMultipartForm(MaxMultipartBytes)
	}

	// Iterate over multipart form data
	if ctx.Req.MultipartForm != nil {
		for key, values := range ctx.Req.MultipartForm.Value {
			for _, value := range values {
				handler(key, value)
			}
		}
	}
}

// handleFormFile iterates over multipart form files and applies the handler function.
func (ctx *HttpContext) handleFormFile(handler FileHandler) {
	// Parse multipart form data if not already parsed
	if ctx.Req.MultipartForm == nil {
		ctx.Req.ParseMultipartForm(MaxMultipartBytes)
	}

	// Iterate over multipart form files
	if ctx.Req.MultipartForm != nil {
		for key, files := range ctx.Req.MultipartForm.File {
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					continue
				}
				handler(key, file, fileHeader)
				file.Close() // Ensure the file is closed after handling
			}
		}
	}
}

// handleFormData iterates over all form data (including multipart form data) and applies the handlers for both form data and files.
func (ctx *HttpContext) handleFormData(formDataHandler FormDataHandler, fileHandler FileHandler) {
	// Parse standard form data if not already parsed
	if ctx.Req.Form == nil {
		ctx.Req.ParseForm()
	}

	// Iterate over standard form data
	for key, values := range ctx.Req.Form {
		for _, value := range values {
			formDataHandler(key, value)
		}
	}

	// Parse multipart form data if not already parsed
	if ctx.Req.MultipartForm == nil {
		ctx.Req.ParseMultipartForm(MaxMultipartBytes)
	}

	// Iterate over multipart form data
	if ctx.Req.MultipartForm != nil {
		// Handle form values
		for key, values := range ctx.Req.MultipartForm.Value {
			for _, value := range values {
				formDataHandler(key, value)
			}
		}

		// Handle files
		for key, files := range ctx.Req.MultipartForm.File {
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					continue
				}
				defer file.Close()
				fileHandler(key, file, fileHeader)
			}
		}
	}
}

// ParmStr parses a string parameter from the request.
func (ctx *HttpContext) ParmStr(key string) string {
	return ctx.parmQuery(key)
}

// ParmInt parses an integer parameter from the request.
func (ctx *HttpContext) ParmInt(key string) (int64, error) {
	return Str2Int64(ctx.parmQuery(key))
}

// ParmInt_ is an alternate method to parse an integer parameter from the request.
func (ctx *HttpContext) ParmInt_(key string) int64 {
	i, _ := Str2Int64(ctx.parmQuery(key))
	return i
}

// ParmFloat parses a float parameter from the request.
func (ctx *HttpContext) ParmFloat(key string) (float64, error) {
	return Str2Float(ctx.parmQuery(key))
}

// ParmBool parses a boolean parameter from the request.
func (ctx *HttpContext) ParmBool(key string) (bool, error) {
	return Str2Bool(ctx.parmQuery(key))
}

// ParmDate parses a date parameter (yyyy-mm-dd) from the request.
func (ctx *HttpContext) ParmDate(key string) (time.Time, error) {
	return time.Parse("2006-01-02", ctx.parmQuery(key))
}

// ParmTime parses a time parameter (HH:MM) from the request.
func (ctx *HttpContext) ParmTime(key string) (time.Time, error) {
	var t time.Time
	p := ctx.parmQuery(key)
	if len(p) >= 5 {
		return time.Parse("15:04", p[len(p)-5:])
	}
	return t, errors.New("ParmTime error")
}

// ParmDataTime parses a date-time parameter (yyyy-mm-dd hh:mm:ss) from the request.
func (ctx *HttpContext) ParmDataTime(key string) (time.Time, error) {
	p := ctx.parmQuery(key)
	return time.Parse("2006-01-02 15:04:05", p)
}

// ParmMonth parses a month parameter (yyyy-mm) from the request.
func (ctx *HttpContext) ParmMonth(key string) (time.Time, error) {
	return time.Parse("2006-01", ctx.parmQuery(key))
}

// ParmYear parses a year parameter (yyyy) from the request.
func (ctx *HttpContext) ParmYear(key string) (time.Time, error) {
	return time.Parse("2006", ctx.parmQuery(key))
}

// ParmTimeParse parses a custom time format from the request.
func (ctx *HttpContext) ParmTimeParse(key string, format string) (time.Time, error) {
	return time.Parse(format, ctx.parmQuery(key))
}

// ParmStrings retrieves multiple values for a parameter from the request.
func (ctx *HttpContext) ParmStrings(key string) []string {
	if ctx.Req.Form == nil {
		ctx.Req.ParseForm()
	}
	if ctx.Req.Form != nil && len(ctx.Req.Form[key]) > 0 {
		return ctx.Req.Form[key]
	}
	return nil
}
