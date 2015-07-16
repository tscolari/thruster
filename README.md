# thruster [![Build Status](https://travis-ci.org/tscolari/thruster.svg?branch=master)](https://travis-ci.org/tscolari/thruster)
Simple wrapper around gin.Engine, that reads configuration from a config file.

## Basic Usage

```go
  config := thruster.Config{
    Hostname: "localhost",
    Port:     3000,
  }
  server := thruster.NewServer(config)

  handler := func(c *gin.Context) {
    ...
  }
  server.AddHandler(thruster.GET, "my_handler", handler)
  server.Start()

  # GET http://localhost/my_handler
```

## JSON


```go
  ...
  server := thruster.NewServer(config)

  handler := func(c *gin.Context) (interface{}, error) {
    return map[string]string{"my":"response"}, nil
  }
  server.AddJSONHandler(thruster.GET, "my_handler", handler)
  server.Start()

  # GET http://localhost/my_handler
  # => {"my":"response"}
```

## HTTP Auth

```go
  config := thruster.Config{
    Hostname: "localhost",
    Port:     3000,
    HTTPAuth: []thruster.HTTPAuth{
      thurster.NewHTTPAuth("user1", "passwd1"),
      thurster.NewHTTPAuth("user2", "passwd2"),
      ...
    },
  }
  server := thruster.NewServer(config)

  # GET http://user1:passwd1@localhost/
  # GET http://user2:passwd2@localhost/
```

## TLS

#### Option1 - certificate in a file
```go
  config := thruster.Config{
    Hostname:    "localhost",
    Port:        3000,
    TLS:         true,
    Certificate: "/path/to/certificate",
    PublicKey:   "/path/to/public/key",
  }
  server := thruster.NewServer(config)

  # GET https://localhost/
```

#### Option2 - inline certificate
```go
  config := thruster.Config{
    Hostname:    "localhost",
    Port:        3000,
    TLS:         true,
    Certificate: "-----BEGIN CERTIFICATE----- .... ",
    PublicKey:   "-----BEGIN RSA PRIVATE KEY----- ....",
  }
  server := thruster.NewServer(config)

  # GET https://localhost/
```

## RESTful Resource

```go
  ...
  server.AddResource("users", controller)

  # GET /users -> controller.Index
  # GET /users/1 -> controller.Show
  # POST /users -> controller.Create
  # PUT /users/1 -> controller.Update
  # DELETE /users/1 -> controller.Destroy
```

```go
  ...
  server.AddJSONResource("users", jsonController)

  # GET /users -> jsonController.Index
  # GET /users/1 -> jsonController.Show
  # POST /users -> jsonController.Create
  # PUT /users/1 -> jsonController.Update
  # DELETE /users/1 -> jsonController.Destroy
```

## Reading configuration from YAML

```go
  config, err := thruster.NewConfig("/path/to/config.yml")
```

### Sample Config.yaml

```yaml
  hostname: localhost
  port: 8888
  http_auth:
  - username: admin
    password: 12345
  - username: user1
    password: 6666
  tls: true
  certificate: /etc/certificate1
  public_key: /etc/public_key
```
