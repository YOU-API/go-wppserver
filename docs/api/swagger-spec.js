var spec =
{
  "info": {
    "title": "WppServer - Send Messages with WhatsApp API",
    "version": "1.0.0",
    "description": "WppServer is an unofficial WhatsApp API server, which is a software programming platform created by third parties that allows developers to integrate WhatsApp into their existing applications or systems. This API offers a wide range of features and functionalities, including the ability to send and receive messages, images and documents, obtain user profile information and more. <br><br> Make sure you have a user. When starting the server for the first time, you can set up an initial email and password for the administrator user through the <code>.env</code> file. Alternatively, you can also use the <a href='#/User/post_user'>user</a> endpoint to register a new user via API. After having a user, it's necessary to provide the email and password in the <a href='#/Auth/post_auth'>auth</a> route to obtain an access token and thus perform the other operations in the API",
    "contact": {
      "email": "williansantanami@gmail.com"
    }
  },
  "externalDocs": {
    "description": "Find out more about",
    "url": "http://localhost:3000/docs/index.html"
  },
  "security": [
    {
      "oauth2": []
    }
  ],
  "paths": {
    "/auth": {
      "post": {
        "tags": [
          "Auth"
        ],
        "description": "A auth request is used to exchange authorization credentials for an access token. Requests to the auth endpoint are authenticated using client credentials (API Keys)  or through username and password authentication.",
        "summary": "Get Token",
        "parameters": [
          {
            "name": "grant_type",
            "description": "The grant_type is a way to specify which authorization flow the client wants to use to obtain an access token. The two supported types are: password and client_credentials.",
            "in": "query",
            "type": "string",
            "x-example": "password"
          },
          {
            "name": "username",
            "description": "Used when the client wants to obtain an access token using the user's credentials (grant_type is password).",
            "in": "query",
            "type": "string",
            "x-example": "admin"
          },
          {
            "name": "password",
            "description": "Used when the client wants to obtain an access token using the user's credentials (grant_type is password),",
            "in": "query",
            "type": "string",
            "x-example": "root"
          },
          {
            "name": "client_id",
            "description": "Used when the client wants to obtain an access token using API keys (grant_type is client_credentials).",
            "in": "query",
            "type": "string",
            "x-example": "cd_ce952554bf5a11edafa10242ac120002"
          },
          {
            "name": "client_secret",
            "description": "Used when the client wants to obtain an access token using API keys (grant_type is client_credentials).",
            "in": "query",
            "type": "string",
            "x-example": "cs_2242ec18a6c94f02bb90b4abbb6a3df5"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/AccessTokenModel"
            }
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "produces": [
          "application/json"
        ]
      }
    },
    "/user": {
      "post": {
        "tags": [
          "User"
        ],
        "summary": "Register User",
        "description": "The endpoint allows users to register. This endpoint typically receives user registration information such as name, email address, and password.",
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/UserModel"
            }
          },
          "400": {
            "description": "Bad request"
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/RegisterUserModel"
            }
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ]
      },
      "get": {
        "tags": [
          "User"
        ],
        "summary": "Get User",
        "description": "The endpoint allows users to obtain information about their accounts. This endpoint normally does not receive parameters.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/UserModel"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          }
        },
        "parameters": [],
        "produces": [
          "application/json"
        ],
        "consumes": [
          "application/json"
        ]
      },
      "put": {
        "tags": [
          "User"
        ],
        "summary": "Update User",
        "description": "The endpoint allows users to update information in their accounts based on the information provided in the request. This endpoint usually receives Name and Email parameters for updating. If the request is made by an administrator, the Type and Status parameters can also be received in the request, and in these case, the user ID parameter must be sent in the body of the request.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/UserModel"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/UpdateUserModel"
            }
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ]
      },
      "delete": {
        "tags": [
          "User"
        ],
        "summary": "Delete User",
        "description": "The endpoint is responsible for deleting a specific user from the system. When a DELETE request is sent to this endpoint with the user ID to be deleted, the system verifies that the user exists and then removes all information associated with the user.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "204": {
            "description": "Successful response. No content"
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/DeleteUserModel"
            }
          }
        ],
        "produces": [
          "application/json"
        ]
      }
    },
    "/user/password": {
      "put": {
        "tags": [
          "User"
        ],
        "summary": "Update Password User",
        "description": "The endpoint allows users to update their access passwords. By making a request to that endpoint with the correct credentials, the user can submit a new password that will replace the old password.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "204": {
            "description": "Successful response. No content"
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/UpdatePasswordModel"
            }
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ]
      }
    },
    "/users": {
      "get": {
        "tags": [
          "User"
        ],
        "summary": "Get All Users",
        "description": "The endpoint allows administrators to obtain a complete list of all registered users in the system. When making a request to this endpoint, the application will receive a response containing the details of all users, including information such as name, email address, encrypted password and other relevant information.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/UserModel"
              }
            }
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [],
        "produces": [
          "application/json"
        ]
      }
    },
    "/users/findusers": {
      "get": {
        "tags": [
          "User"
        ],
        "summary": "Find Users",
        "description": "The endpoint allows administrators to perform searches for specific users based on a keyword. Upon making a request to that endpoint with the correct search parameters the app will receive a response containing a list of users.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/UserModel"
              }
            }
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [
          {
            "name": "keyword",
            "description": "The field allows you to enter a relevant keyword that will be used to filter users.",
            "in": "query",
            "type": "string",
            "x-example": "Méndez"
          }
        ],
        "produces": [
          "application/json"
        ]
      }
    },
    "/device/login": {
      "get": {
        "tags": [
          "Device"
        ],
        "summary": "Login Device",
        "description": "When making a request to this endpoint, the system generates a unique QR code for the device in question. The user can then scan the QR code with their mobile device's camera, which will allow the system to authenticate the user's account and grant access to API operations on behalf of the number associated with the device.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/LoginDeviceModel"
            }
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [],
        "produces": [
          "application/json"
        ]
      }
    },
    "/device/logout": {
      "get": {
        "tags": [
          "Device"
        ],
        "summary": "Logout Device",
        "description": "The endpoint is used to allow a user to log out of a device associated with their API user account.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "204": {
            "description": "Successful response. No content"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [],
        "produces": [
          "application/json"
        ]
      }
    },
    "/session/connect": {
      "post": {
        "tags": [
          "Session"
        ],
        "summary": "Connect Session",
        "description": "The endpoint is used to establish a session connection between a specific device and the application server. When making a request to this endpoint, the system verifies that the device is registered with the system and has the necessary permissions to connect to the session. Once the session connection has been successfully established, the system can send and receive data on behalf of the connected device.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "204": {
            "description": "Successful response. No content"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [],
        "produces": [
          "application/json"
        ]
      }
    },
    "/session/disconnect": {
      "post": {
        "tags": [
          "Session"
        ],
        "summary": "Disconnect Session",
        "description": "The endpoint is used to terminate a previously established session connection between a specific device and the application server. When making a request to this endpoint, the system verifies that the session connection exists and then securely and efficiently terminates the connection. This allows the system to free up resources and memory used to manage the session, improving overall system efficiency. As long as the device remains logged into the application, it is possible to reconnect whenever necessary.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "204": {
            "description": "Successful response. No content"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [],
        "produces": [
          "application/json"
        ]
      }
    },
    "/session/status": {
      "get": {
        "tags": [
          "Session"
        ],
        "summary": "Status Session",
        "description": "The endpoint is used to obtain information about session status and the existence of devices logged into the API account.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/StatusSessionModel"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [],
        "produces": [
          "application/json"
        ]
      }
    },
    "/chat/send/text": {
      "post": {
        "tags": [
          "Send"
        ],
        "summary": "Send Text",
        "description": "The endpoint is used to send texts to a specific device number via the application server.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/MessageResponseModel"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/SendTextModel"
            }
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ]
      }
    },
    "/chat/send/image": {
      "post": {
        "tags": [
          "Send"
        ],
        "summary": "Send Image",
        "description": "The endpoint is used to send images to a specific device number via the application server. The contents of the file must be base64 encoded.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/MessageResponseModel"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/SendImageModel"
            }
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ]
      }
    },
    "/chat/send/document": {
      "post": {
        "tags": [
          "Send"
        ],
        "summary": "Send Document",
        "description": "The endpoint is used to send a document to a specific device number via the application server. The contents of the file must be base64 encoded.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/MessageResponseModel"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/SendDocumentModel"
            }
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ]
      }
    },
    "/phone/contacts": {
      "get": {
        "tags": [
          "Phone"
        ],
        "summary": "Get Contacts",
        "description": "The endpoint is used to retrieve saved contact information from a specific device. When making a request to this endpoint, the system retrieves contact information stored on the device, including name, phone number and other relevant information.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/PhoneInfoModel"
            }
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [],
        "produces": [
          "application/json"
        ]
      }
    },
    "/phone/scraping": {
      "post": {
        "tags": [
          "Phone"
        ],
        "summary": "Scraping Phones",
        "description": "The endpoint is used to extract contact information from one or more phone numbers. When making a request to this endpoint, the system uses data scraping techniques to extract public information from users.",
        "security": [
          {
            "oauth2": []
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/PhoneInfoModel"
              }
            }
          },
          "400": {
            "description": "Bad request"
          },
          "401": {
            "description": "Unauthorized"
          },
          "500": {
            "description": "Internal server error"
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              },
              "example": [
                "5491155553934",
                "5491155553935"
              ]
            }
          }
        ],
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ]
      }
    }
  },
  "swagger": "2.0",
  "host": "localhost:3000",
  "tags": [
    {
      "name": "Auth",
      "description": "User authentication route"
    },
    {
      "name": "User",
      "description": "User management route"
    },
    {
      "name": "Device",
      "description": "Device management route"
    },
    {
      "name": "Session",
      "description": "Device session management route"
    },
    {
      "name": "Send",
      "description": "Data sending route"
    },
    {
      "name": "Phone",
      "description": "Phone number management route"
    }
  ],
  "schemes": [
    "http",
    "https"
  ],
  "basePath": "/v1",
  "securityDefinitions": {
    "oauth2": {
      "type": "oauth2",
      "flow": "password",
      "tokenUrl": "http://localhost:3000/v1/auth",
      "scopes": {}
    }
  },
  "definitions": {
    "AccessTokenModel": {
      "type": "object",
      "properties": {
        "access_token": {
          "type": "string",
          "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImFwaWtleWlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwidXNlcmlkIjoiOGI1M2QwZDctYTM4Yy0xMWVkLWI5ZWUtNTBiN2MzMDA4NDMzIiwiU2NvcGUiOiJhZG1pbjoqIiwiZXhwIjoxNjc4Mzg1NjAwfQ.-k7Pm8E9R2hMPjFfjeun_7ZdKsLaT7RcCvzXlMDWR28"
        }
      }
    },
    "UserModel": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "58fd4474-8bb0-41a0-ad5f-5b492fbdfa2c"
        },
        "name": {
          "type": "string",
          "example": "Raúl Méndez"
        },
        "email": {
          "type": "string",
          "example": "rmendez1985@gmail.com"
        },
        "type": {
          "type": "string",
          "enum": [
            "admin",
            "agent"
          ],
          "example": "admin"
        },
        "status": {
          "type": "string",
          "enum": [
            "enabled",
            "disabled"
          ],
          "example": "enabled"
        }
      }
    },
    "RegisterUserModel": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "example": "Raúl Méndez"
        },
        "email": {
          "type": "string",
          "example": "rmendez1985@gmail.com"
        },
        "password": {
          "type": "string",
          "example": 5081985
        }
      },
      "required": [
        "name",
        "email",
        "password"
      ]
    },
    "UpdateUserModel": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "format": "uuid",
          "example": "58fd4474-8bb0-41a0-ad5f-5b492fbdfa2c"
        },
        "name": {
          "type": "string",
          "example": "Raúl Méndez"
        },
        "email": {
          "type": "string",
          "format": "email",
          "example": "rmendez1985@gmail.com"
        },
        "type": {
          "type": "string",
          "enum": [
            "admin",
            "agent"
          ],
          "example": "admin"
        },
        "status": {
          "type": "string",
          "enum": [
            "enabled",
            "disabled"
          ],
          "example": "enabled"
        }
      },
      "required": [
        "name",
        "email"
      ]
    },
    "DeleteUserModel": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "format": "uuid",
          "example": "58fd4474-8bb0-41a0-ad5f-5b492fbdfa2c"
        }
      },
      "required": [
        "id"
      ]
    },
    "UpdatePasswordModel": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "format": "uuid",
          "example": "58fd4474-8bb0-41a0-ad5f-5b492fbdfa2c"
        },
        "oldpassword": {
          "type": "string",
          "example": 5081985
        },
        "newpassword": {
          "type": "string",
          "example": "rmendez1985"
        }
      },
      "required": [
        "oldpassword",
        "newpassword"
      ]
    },
    "LoginDeviceModel": {
      "type": "object",
      "properties": {
        "base64qrcode": {
          "type": "string",
          "example": "data:image/png;base64,iVBORw0KGg...[base64-encoded data]"
        },
        "expiration": {
          "type": "string",
          "format": "date-time",
          "example": "2006-01-02T15:04:05-0700"
        }
      }
    },
    "StatusSessionModel": {
      "type": "object",
      "properties": {
        "islogged": {
          "type": "boolean",
          "example": true
        },
        "isconnected": {
          "type": "boolean",
          "example": true
        }
      }
    },
    "MessageResponseModel": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "f5a92da61d94481ea0e991b8868d716e"
        },
        "details": {
          "type": "string",
          "example": "Sent"
        },
        "timestamp": {
          "type": "string",
          "example": "1678652549"
        }
      },
      "required": [
        "phone",
        "filename",
        "document"
      ]
    },
    "SendTextModel": {
      "type": "object",
      "properties": {
        "phone": {
          "type": "string",
          "example": "+5511932486442"
        },
        "body": {
          "type": "string",
          "example": "Hello my dear friend...."
        }
      },
      "required": [
        "phone",
        "body"
      ]
    },
    "SendImageModel": {
      "type": "object",
      "properties": {
        "phone": {
          "type": "string",
          "example": "+5511932486442"
        },
        "image": {
          "type": "string",
          "example": "data:image/png;base64,UEsDBBQAAAAIAEcLMladeZbNGgAAABgAAAAIAAAAZG93bmxvYWTLyM9JVCgsTeUqSczhSs7PzedKLS5JLOYCAFBLAQIfABQAAAAIAEcLMladeZbNGgAAABgAAAAIACQAAAAAAAAAIAAAAAAAAABkb3dubG9hZAoAIAAAAAAAAQAYAOkgZwD1KtkB6SBnAPUq2QGL/Iv/9CrZAVBLBQYAAAAAAQABAFoAAABAAAAAAAA="
        }
      },
      "required": [
        "phone",
        "image"
      ]
    },
    "SendDocumentModel": {
      "type": "object",
      "properties": {
        "phone": {
          "type": "string",
          "example": "+5511932486442"
        },
        "filename": {
          "type": "string",
          "example": "My file.zip"
        },
        "document": {
          "type": "string",
          "example": "data:application/zip;base64,UEsDBBQAAAAIAEcLMladeZbNGgAAABgAAAAIAAAAZG93bmxvYWTLyM9JVCgsTeUqSczhSs7PzedKLS5JLOYCAFBLAQIfABQAAAAIAEcLMladeZbNGgAAABgAAAAIACQAAAAAAAAAIAAAAAAAAABkb3dubG9hZAoAIAAAAAAAAQAYAOkgZwD1KtkB6SBnAPUq2QGL/Iv/9CrZAVBLBQYAAAAAAQABAFoAAABAAAAAAAA="
        }
      },
      "required": [
        "phone",
        "filename",
        "document"
      ]
    },
    "PhoneInfoModel": {
      "type": "object",
      "properties": {
        "jid": {
          "type": "string",
          "example": "5511932486442.0:22@s.whatsapp.net"
        },
        "phone": {
          "type": "string",
          "example": "+5511932486442"
        },
        "onwhatsapp": {
          "type": "boolean",
          "example": true
        },
        "iscontact": {
          "type": "boolean",
          "example": true
        },
        "pessoalname": {
          "type": "string",
          "example": "Raúl Méndez"
        },
        "businessname": {
          "type": "string",
          "example": "Raúl Méndez Technology S.A"
        },
        "pictureurl": {
          "type": "string",
          "example": "data:image/png;base64,UEsDBBQAAAAIAEcLMladeZbNGgAAABgAAAAIAAAAZG93bmxvYWTLyM9JVCgsTeUqSczhSs7PzedKLS5JLOYCAFBLAQIfABQAAAAIAEcLMladeZbNGgAAABgAAAAIACQAAAAAAAAAIAAAAAAAAABkb3dubG9hZAoAIAAAAAAAAQAYAOkgZwD1KtkB6SBnAPUq2QGL/Iv/9CrZAVBLBQYAAAAAAQABAFoAAABAAAAAAAA="
        }
      }
    }
  }
}