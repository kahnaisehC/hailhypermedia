package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Contact struct {
	Id    int
	Name  string
	Email string
	Phone string
}

func filterContacts(reg string, contacts []Contact) []Contact {
	return contacts
}

func getContacts() ([]Contact, error) {
	var contacts []Contact
	f, err := os.Open("contacts.csv")
	if err != nil {
		f, err = os.Create("contacts.csv")
		if err != nil {
			return []Contact{}, err
		}
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		return []Contact{}, err
	}
	for _, csvContact := range data {
		contactId, err := strconv.Atoi(csvContact[0])
		if err != nil {
			fmt.Printf("cannot convert %s to integer\n", csvContact[0])
			continue
		}
		contactName := csvContact[1]
		contactEmail := csvContact[2]
		contactPhone := csvContact[3]

		contact := Contact{
			contactId,
			contactName,
			contactEmail,
			contactPhone,
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func Contacts(c echo.Context) error {
	contacts, err := getContacts()
	if err != nil {
		return err
	}

	queryParam := c.QueryParam("q")
	if queryParam != "" {
		contacts = filterContacts(queryParam, contacts)
	}

	return c.Render(http.StatusOK, "contacts", struct {
		Contacts []Contact
		Q        string
	}{
		Contacts: contacts,
		Q:        queryParam,
	})
}

func NewContactView(c echo.Context) error {
	return c.Render(http.StatusOK, "new", "")
}

func CreateContact(c echo.Context) error {
	return c.Render(http.StatusOK, "new", "")
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	e.Static("", "public")

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/contacts")
	})
	e.GET("/contacts", Contacts)

	e.GET("/contacts/new", NewContactView)
	e.POST("/contacts/new", CreateContact)

	e.Logger.Fatal(e.Start(":1323"))
}
