package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func PageRenderer(page string, data map[string]string, c echo.Context) error {
	file, _ := os.ReadFile("templates/" + page)
	output := string(file)
	for key, value := range data {
		key = "<!-- " + key + " -->"
		output = strings.Replace(output, key, value, -1)
	}
	return c.HTML(200, output)
}

func main() {
	// Initialization
	e := echo.New()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL is not set. Please set it to your database URL.")
	}
	dbURL, _ = pq.ParseURL(dbURL)
	db, _ := sql.Open("postgres", dbURL)
	defer db.Close()

	// Ping the database to make sure we can connect.
	err := db.Ping()
	fmt.Println("Connecting to database...")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Successfully connected to the database.")
	}

	// Setup Database
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS data (key TEXT, url TEXT, date TEXT)")
	if err != nil {
		panic(err.Error())
	}

	// Index Handler
	e.GET("/", func(c echo.Context) error {
		file, _ := os.ReadFile("templates/index.html")
		return c.HTML(200, string(file))
	})

	// Redirect Handler
	e.GET("/:key", func(c echo.Context) error {
		key := c.Param("key")
		var url string
		err := db.QueryRow("SELECT url FROM data WHERE key = $1", key).Scan(&url)
		if err != nil {
			fmt.Println(err.Error())
			return c.String(404, "404: Not Found")
		}
		return c.Redirect(301, url)
	})

	// API Handler
	e.POST("/api", func(c echo.Context) error {
		url := c.FormValue("url")
		key := c.FormValue("key")
		date := time.Now().Format("2006-01-02 15:04:05")
		_, err := db.Exec("INSERT INTO data (key, url, date) VALUES ($1, $2, $3)", key, url, date)
		c.Response().Header().Set("Content-Type", "application/json")
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(500, map[string]string{"status": "error"})
		} else {
			return c.JSON(200, map[string]string{"status": "ok"})
		}
	})

	// Admin Handler
	e.GET("/login", func(c echo.Context) error {
		file, _ := os.ReadFile("templates/login.html")
		return c.HTML(200, string(file))
	})

	e.POST("/login", func(c echo.Context) error {
		passwd, _ := os.ReadFile("passwd.txt")
		if c.FormValue("passwd") == string(passwd) {
			return c.String(200, "200: Login Success")
		} else {
			return c.String(401, "401: Unauthorized")
		}
	})

	e.GET("/admin", func(c echo.Context) error {
		cookie, err := c.Cookie("passwd")
		if err != nil {
			return c.Redirect(302, "/login")
		}
		passwd, _ := os.ReadFile("passwd.txt")
		if cookie.Value != string(passwd) {
			return c.Redirect(302, "/login")
		}

		rows, err := db.Query("SELECT * FROM data")
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()

		table := ""
		for rows.Next() {
			var key, url, date string
			rows.Scan(&key, &url, &date)
			deleteUrl := fmt.Sprintf("<a href=\"/admin/delete/%s\">Delete</a>", key)
			keyUrl := fmt.Sprintf("<a href=\"/%s\">%s</a>", key, key)
			table += fmt.Sprintf("\n<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", keyUrl, url, date, deleteUrl)
		}
		PageRenderer("admin.html", map[string]string{"table": table}, c)
		return nil
	})

	e.GET("/admin/delete/:key", func(c echo.Context) error {
		cookie, err := c.Cookie("passwd")
		if err != nil {
			return c.Redirect(302, "/login")
		}
		passwd, _ := os.ReadFile("passwd.txt")
		if cookie.Value != string(passwd) {
			return c.Redirect(302, "/login")
		}

		key := c.Param("key")
		_, err = db.Exec("DELETE FROM data WHERE key = $1", key)
		if err != nil {
			fmt.Println(err.Error())
			return c.String(500, "500: Internal Server Error")
		} else {
			return c.Redirect(301, "/admin")
		}
	})

	// Start Web Server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	e.Logger.Fatal(e.Start(":" + PORT))
}
