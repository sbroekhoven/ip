# ip

This is a simple Go web application that returns the visitor's IP address and browser information. The response format is based on the `Accept` header and supports:

* `text/html` (for browsers)
* `text/plain`
* `application/json`


### Build and Run

#### Option 1: Using Docker Compose

```bash
docker-compose up --build
```

This builds the Go app, starts both the app and the Caddy reverse proxy.

Access the service at: [http://localhost](http://localhost)


### Updating the Go App

After making changes to `main.go` or the templates:

```bash
docker-compose up --build
```

To rebuild the app container and apply the updates.


### Testing with curl

#### Get plain text output:

```bash
curl -H "Accept: text/plain" http://localhost/
```

#### Get JSON output:

```bash
curl -H "Accept: application/json" http://localhost/
```

#### Simulate a browser request (HTML):

```bash
curl -H "Accept: text/html" http://localhost/
```
