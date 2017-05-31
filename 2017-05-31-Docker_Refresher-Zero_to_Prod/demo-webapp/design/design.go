package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("WorstSocialNetwork", func() {
	Title("The worst social network")
	Description("A terrible social network")
	Scheme("http")
	Host("localhost:8000")
})

var _ = Resource("health", func() {
	BasePath("/health")
	DefaultMedia(HealthMedia)

	Action("show", func() {
		Description("show health")
		Routing(GET("/"))
		Response(OK)
	})
})

// HealthMedia -
var HealthMedia = MediaType("application/vnd.health+json", func() {
	Description("service health")
	Attributes(func() {
		Attribute("hostname", String, "the hostname")
		Attribute("version", String, "the version")
		Required("hostname", "version")
	})
	View("default", func() {
		Attribute("hostname")
		Attribute("version")
	})
})

var _ = Resource("post", func() {
	BasePath("/posts")
	DefaultMedia(PostMedia)

	Action("show", func() {
		Description("Get post by id")
		Routing(GET("/:postID"))
		Params(func() {
			Param("postID", Integer, "Post ID")
		})
		Response(OK)
		Response(NotFound)
	})
})

// PostMedia -
var PostMedia = MediaType("application/vnd.wsn.post+json", func() {
	Description("A post")
	Attributes(func() {
		Attribute("id", Integer, "Unique post ID")
		Attribute("title", String, "post title")
		Attribute("body", String, "post body")
		Attribute("userId", Integer, "Owner's ID")
		Required("id", "title", "body", "userId")
	})
	View("default", func() {
		Attribute("id")
		Attribute("userId")
		Attribute("title")
		Attribute("body")
	})
})

var _ = Resource("public", func() {
	Origin("*", func() {
		Methods("GET")
	})

	Files("/swagger.json", "public/swagger/swagger.json")
	Files("/schema.json", "public/schema/schema.json")
	Files("/", "public/html/index.html")
	Files("/js/*filepath", "public/js/")
})
