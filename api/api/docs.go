// Package api GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package api

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/greeting/{name}": {
            "get": {
                "description": "Greeting by given name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Greeting"
                ],
                "summary": "Say greeting",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name for greeting",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/http.getGreetingRes"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/users": {
            "post": {
                "description": "Create a user with basic information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Create a user",
                "parameters": [
                    {
                        "description": "request to create a new user",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.createUserReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/http.createUserRes"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/users/{id}": {
            "get": {
                "description": "Get a user with given id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/http.getUserRes"
                        }
                    },
                    "400": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.createUserReq": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "familyName": {
                    "type": "string"
                },
                "givenName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "http.createUserRes": {
            "type": "object",
            "properties": {
                "credit": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "familyName": {
                    "type": "string"
                },
                "givenName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "http.getGreetingRes": {
            "type": "object",
            "properties": {
                "msg": {
                    "type": "string"
                }
            }
        },
        "http.getUserRes": {
            "type": "object",
            "properties": {
                "credit": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "familyName": {
                    "type": "string"
                },
                "givenName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Imageer Endpoint API",
	Description: "Endpoint API for image processing service.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
