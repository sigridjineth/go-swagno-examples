package handlers

import (
	"github.com/go-swagno/examples/fiber/models"
	. "github.com/go-swagno/swagno"
	swaggerUi "github.com/go-swagno/swagno-files"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) SetRoutes(a *fiber.App) {
	a.Get("/hello", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello World!")
	}).Name("index")
}

func (h *Handler) SetSwagger(a *fiber.App) {
	endpoints := []Endpoint{
		EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil),
		EndPoint(GET, "/product", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),

		// no return
		EndPoint(POST, "/product-no-return", "product", Params(), nil, nil, models.ErrorResponse{}, "", nil),
		// no error
		EndPoint(POST, "/product-no-error", "product", Params(), nil, nil, nil, "", nil),

		// ids query enum
		EndPoint(GET, "/products", "product", Params(IntEnumQuery("ids", []int64{1, 2, 3}, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		// ids path enum
		EndPoint(GET, "/products2", "product", Params(IntEnumParam("ids", []int64{1, 2, 3}, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		// with fields
		EndPoint(GET, "/productsMinMax", "product", Params(IntArrQuery("ids", nil, true, "test", Fields{Min: 0, Max: 10, Default: 5})), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		// string array query
		EndPoint(GET, "/productsArr", "product", Params(StrArrQuery("strs", nil, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(GET, "/productsArrWithEnums", "product", Params(StrArrQuery("strs", []string{"test1", "test2"}, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(GET, "/productsArrWithEnumsInPath", "product", Params(StrArrParam("strs", []string{"test1", "test2"}, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),

		// /merchant/{merchantId}?id={id} -> get product of a merchant
		EndPoint(GET, "/merchant", "merchant", Params(StrParam("merchant", true, ""), IntQuery("id", true, "product id")), nil, models.Product{}, models.ErrorResponse{}, "", nil),

		// with headers
		EndPoint(POST, "/product-header", "header params", Params(IntHeader("header1", false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(POST, "/product2-header", "header params", Params(IntEnumHeader("header1", []int64{1, 2, 3}, false, ""), StrEnumHeader("header2", []string{"a", "b", "c"}, false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(POST, "/product3-header", "header params", Params(IntArrHeader("header1", []int64{1, 2, 3}, false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),

		// with file
		EndPoint(POST, "/productUpload", "upload", Params(FileParam("file", true, "File to upload")), nil, models.Product{}, models.ErrorResponse{}, "", nil),

		// without EndPoint function
		{Method: "GET", Path: "/product4", Description: "product", Params: Params(IntParam("id", true, "")), Return: models.Product{}, Error: models.ErrorResponse{}, Tags: []string{"WithStruct"}},
		// without EndPoint function and without Params
		{Method: "GET", Path: "/product5", Description: "product", Params: []Parameter{{Name: "id", Type: "integer", In: "path", Required: true}}, Return: models.Product{}, Error: models.ErrorResponse{}, Tags: []string{"WithStruct"}},

		// with security
		EndPoint(POST, "/secure-product", "Secure", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "Only Basic Auth", BasicAuth()),
		EndPoint(POST, "/multi-secure-product", "Secure", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "Basic Auth + Api Key Auth", Security(ApiKeyAuth("api_key"), BasicAuth())),
		EndPoint(POST, "/secure-product-oauth", "Secure", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "OAuth", OAuth("oauth2_name", "read:pets")),
	}

	sw := CreateNewSwagger("Swagger API", "1.0")
	AddEndpoints(endpoints)

	// set auth
	sw.SetBasicAuth()
	sw.SetApiKeyAuth("api_key", "query")
	sw.SetOAuth2Auth("oauth2_name", "password", "http://localhost:8080/oauth2/token", "http://localhost:8080/oauth2/authorize", Scopes(Scope("read:pets", "read your pets"), Scope("write:pets", "modify pets in your account")))

	// 3 alternative way for describing tags with descriptions
	sw.AddTags(Tag("product", "Product operations"), Tag("merchant", "Merchant operations"))
	sw.AddTags(SwaggerTag{Name: "WithStruct", Description: "WithStruct operations"})
	sw.Tags = append(sw.Tags, SwaggerTag{Name: "headerparams", Description: "headerparams operations"})

	// if you want to export your swagger definition to a file
	swaggerDocs := string(sw.GenerateDocs())
	a.Use("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendString(swaggerDocs)
	})

	SwaggerHandler(a, sw.GenerateDocs())
}

var swaggerDoc string

func SwaggerHandler(a *fiber.App, doc []byte) {
	if swaggerDoc == "" {
		swaggerDoc = string(doc)
	}

	// Redirect /swagger to the correct Swagger UI index page
	a.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html?url=/swagger/doc.json", fiber.StatusSeeOther)
	})

	// Serve the Swagger JSON at /swagger/doc.json
	a.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendString(swaggerDoc)
	})

	// Serve the Swagger UI assets under /swagger
	a.Use("/swagger", filesystem.New(filesystem.Config{
		Root:       swaggerUi.HTTP,
		PathPrefix: "", // Ensure this is empty
	}))
}
