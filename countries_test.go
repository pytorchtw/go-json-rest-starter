package main

// code adapted from https://github.com/ant0ine/go-json-rest

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/gavv/httpexpect"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var api *rest.Api
var handler http.Handler
var server *httptest.Server

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func setup() {
	api = rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/countries/:code", HandleGetCountry()),
		rest.Delete("/countries/:code", HandleDeleteCountry()),
		rest.Get("/countries", HandleGetAllCountries()),
		rest.Post("/countries", HandlePostCountry()),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	handler = api.MakeHandler()
	server = httptest.NewServer(handler)
}

func teardown() {
	defer server.Close()
	fmt.Println("test completed...")
}

func TestHandlePostCountry(t *testing.T) {
	country := Country{Code: "US", Name: "United States"}
	recorded := test.RunRequest(t, handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/countries", country))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()

	country = Country{Code: "FR", Name: "France"}
	recorded = test.RunRequest(t, handler,
		test.MakeSimpleRequest("POST", "http://1.2.3.4/countries", country))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
}

func TestHandleGetCountry(t *testing.T) {
	recorded1 := test.RunRequest(t, handler,
		test.MakeSimpleRequest("GET", "http://1.2.3.4/countries/US", nil))
	recorded1.CodeIs(200)
}

func TestHandleGetAllCountries(t *testing.T) {
	recorded := test.RunRequest(t, handler,
		test.MakeSimpleRequest("GET", "http://1.2.3.4/countries", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
}

func TestGetAllCountries(t *testing.T) {
	e := httpexpect.New(t, server.URL)
	e.GET("/countries").Expect().Status(http.StatusOK).JSON().
		Array().Length().Equal(2)

	e.GET("/countries").Expect().Status(http.StatusOK).JSON().
		Array().
		Element(0).
		Object().
		ContainsKey("code").
		ValueEqual("code", "US")
}

func TestHandleDeleteCountry(t *testing.T) {
	recorded := test.RunRequest(t, handler,
		test.MakeSimpleRequest("DELETE", "http://1.2.3.4/countries/US", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()

	e := httpexpect.New(t, server.URL)
	e.DELETE("/countries/FR").Expect().
		Status(http.StatusOK)
}
