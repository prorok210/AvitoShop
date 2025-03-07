// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/api/auth": {
            "post": {
                "description": "Аутентификация с помощью имени пользователя и пароля и возвращение токена.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Аутентификация и получение JWT-токена.",
                "parameters": [
                    {
                        "description": "Auth credentials",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/models.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная аутентификация",
                        "schema": {
                            "$ref": "#/definitions/models.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "$ref": "#/definitions/models.Error400Response"
                        }
                    },
                    "401": {
                        "description": "Неавторизован",
                        "schema": {
                            "$ref": "#/definitions/models.Error401Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/models.Error500Response"
                        }
                    }
                }
            }
        },
        "/api/buy/{item}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Покупка предмета за монеты: списывается стоимость предмета с баланса пользователя и создается заказ.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Merch"
                ],
                "summary": "Купить предмет за монеты.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Название предмета",
                        "name": "item",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный запрос",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос.",
                        "schema": {
                            "$ref": "#/definitions/models.Error400Response"
                        }
                    },
                    "401": {
                        "description": "Неавторизован.",
                        "schema": {
                            "$ref": "#/definitions/models.Error401Response"
                        }
                    },
                    "404": {
                        "description": "Предмет не найден.",
                        "schema": {
                            "$ref": "#/definitions/models.Error404Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера.",
                        "schema": {
                            "$ref": "#/definitions/models.Error500Response"
                        }
                    }
                }
            }
        },
        "/api/info": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Получение баланса, инвентаря и истории транзакций (отправленных и полученных монет) для авторизованного пользователя.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Получить информацию о монетах, инвентаре и истории транзакций.",
                "responses": {
                    "200": {
                        "description": "Успешный ответ",
                        "schema": {
                            "$ref": "#/definitions/models.InfoResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "$ref": "#/definitions/models.Error400Response"
                        }
                    },
                    "401": {
                        "description": "Неавторизован",
                        "schema": {
                            "$ref": "#/definitions/models.Error401Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/models.Error500Response"
                        }
                    }
                }
            }
        },
        "/api/sendCoin": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Перевод монет от одного пользователя к другому.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "Отправить монеты другому пользователю.",
                "parameters": [
                    {
                        "description": "SendCoinRequest",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.TransactionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный запрос",
                        "schema": {
                            "$ref": "#/definitions/models.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "$ref": "#/definitions/models.Error400Response"
                        }
                    },
                    "401": {
                        "description": "Неавторизован",
                        "schema": {
                            "$ref": "#/definitions/models.Error401Response"
                        }
                    },
                    "404": {
                        "description": "Не найдено",
                        "schema": {
                            "$ref": "#/definitions/models.Error404Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/models.Error500Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AuthRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "secret123"
                },
                "username": {
                    "type": "string",
                    "example": "username"
                }
            }
        },
        "models.AuthResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "models.CoinHistory": {
            "type": "object",
            "properties": {
                "received": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ReceivedTx"
                    }
                },
                "sent": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.SentTx"
                    }
                }
            }
        },
        "models.Error400Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Неверный запрос"
                }
            }
        },
        "models.Error401Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Неавторизован"
                }
            }
        },
        "models.Error404Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Не найдено"
                }
            }
        },
        "models.Error500Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Внутренняя ошибка сервера"
                }
            }
        },
        "models.InfoResponse": {
            "type": "object",
            "properties": {
                "coinHistory": {
                    "$ref": "#/definitions/models.CoinHistory"
                },
                "coins": {
                    "type": "integer",
                    "example": 100
                },
                "inventory": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.InventoryItem"
                    }
                }
            }
        },
        "models.InventoryItem": {
            "type": "object",
            "properties": {
                "quantity": {
                    "type": "integer",
                    "example": 2
                },
                "type": {
                    "type": "string",
                    "example": "weapon"
                }
            }
        },
        "models.ReceivedTx": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 50
                },
                "fromUser": {
                    "type": "string",
                    "example": "Alice"
                }
            }
        },
        "models.SentTx": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 30
                },
                "toUser": {
                    "type": "string",
                    "example": "Bob"
                }
            }
        },
        "models.SuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Успешный запрос"
                }
            }
        },
        "models.TransactionRequest": {
            "type": "object",
            "required": [
                "amount",
                "toUser"
            ],
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 100
                },
                "toUser": {
                    "type": "string",
                    "example": "username"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type \"Bearer\" followed by a space and JWT token",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "API Avito shop",
	Description:      "API для отбора на стажировку в Авито",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
