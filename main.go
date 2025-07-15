package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

type VisitorInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Country   string `json:"country,omitempty"`
	City      string `json:"city,omitempty"`
}

var geoDB *geoip2.Reader

// getClientIP gets the client's IP address from the request. It first tries to use the value
// of the X-Real-IP header, then falls back to the RemoteAddr field of the request object.
func getClientIP(c *gin.Context) string {
	realIP := c.Request.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}

// enrichGeoInfo takes an IP address in string format and returns the country and city names
// associated with that IP address.
func enrichGeoInfo(ipStr string) (country string, city string) {
	ip := net.ParseIP(ipStr)
	if ip == nil || geoDB == nil {
		return "", ""
	}

	record, err := geoDB.City(ip)
	if err != nil {
		return "", ""
	}

	country = record.Country.Names["en"]
	if len(record.City.Names) > 0 {
		city = record.City.Names["en"]
	}

	return country, city
}

// main sets up a web server and handles HTTP requests.
func main() {
	// Configure logger
	gin.DefaultWriter = os.Stdout

	// Load port from env, fallback to 5001
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	// Load GeoIP database
	var err error
	geoDB, err = geoip2.Open("geoip/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatalf("error=geoip_db_open_failed detail=%v", err)
	}
	defer geoDB.Close()

	// Create router
	router := gin.New()

	// Set trusted proxies from env
	proxyCIDRs := os.Getenv("TRUSTED_PROXIES")
	if proxyCIDRs == "" {
		proxyCIDRs = "172.18.0.0/16" // default: Docker bridge network
	}
	cidrs := strings.Split(proxyCIDRs, ",")

	if err := router.SetTrustedProxies(cidrs); err != nil {
		log.Fatalf("error=trusted_proxy_config_failed detail=%v", err)
	}
	// Basic structured logging
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("status=%d method=%s path=%s ip=%s ua=%q duration=%s\n",
			param.StatusCode,
			param.Method,
			param.Path,
			param.ClientIP,
			param.Request.UserAgent(),
			param.Latency,
		)
	}))
	router.Use(gin.Recovery())

	// Load HTML templates
	router.LoadHTMLGlob("templates/*.html")

	// Route
	router.GET("/", func(c *gin.Context) {
		start := time.Now()

		ip := getClientIP(c)
		country, city := enrichGeoInfo(ip)

		info := VisitorInfo{
			IP:        ip,
			UserAgent: c.Request.UserAgent(),
			Country:   country,
			City:      city,
		}

		accept := c.GetHeader("Accept")
		switch {
		case strings.Contains(accept, "application/json"):
			c.JSON(http.StatusOK, info)
		case strings.Contains(accept, "text/plain"):
			c.Data(http.StatusOK, "text/plain", []byte(
				fmt.Sprintf("Visitor IP: %s\nUser-Agent: %s\nCountry: %s\nCity: %s\n",
					info.IP, info.UserAgent, info.Country, info.City)))
		default:
			c.HTML(http.StatusOK, "template.html", info)
		}

		log.Printf("event=handled ip=%s method=%s path=%s duration=%s",
			ip, c.Request.Method, c.Request.URL.Path, time.Since(start))
	})

	log.Printf("event=server_start port=%s", port)
	router.Run(":" + port)
}
