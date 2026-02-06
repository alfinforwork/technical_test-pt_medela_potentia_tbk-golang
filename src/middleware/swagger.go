package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/swaggo/swag"
)

// SwaggerHandler serves Swagger UI and swagger.json spec
func SwaggerHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := c.Path()

		// Serve swagger.json
		if strings.HasSuffix(path, "/swagger.json") {
			c.Type("json")
			spec, _ := swag.ReadDoc()
			return c.SendString(spec)
		}

		// Serve Swagger UI HTML
		swaggerUIHTML := getSwaggerUIHTML()
		c.Type("html")
		return c.SendString(swaggerUIHTML)
	}
}

func getSwaggerUIHTML() string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Workflow Management API - Swagger UI</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui.css">
    <link rel="icon" type="image/png" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/favicon-16x16.png" sizes="16x16" />
    <style>
      html{
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }
      *,
      *:before,
      *:after{
        box-sizing: inherit;
      }
      body{
        margin:0;
        background: #fafafa;
      }
    </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
<script>
    window.onload = function() {
        SwaggerUIBundle({
            url: "/swagger.json",
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            plugins: [
                SwaggerUIBundle.plugins.DownloadUrl
            ],
            layout: "StandaloneLayout"
        })
    }
</script>
</body>
</html>
`)
}
