package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/kahnaisehC/hailhypermedia/contacts"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Contacts(c echo.Context) error {
	contactList, err := contacts.GetContacts()
	if err != nil {
		return err
	}

	queryParam := c.QueryParam("q")
	if queryParam != "" {
		contactList = contacts.FilterContacts(queryParam, contactList)
	}

	return c.Render(http.StatusOK, "contacts", struct {
		Contacts []contacts.Contact
		Q        string
	}{
		Contacts: contactList,
		Q:        queryParam,
	})
}

func GetNewContact(c echo.Context) error {
	return c.Render(http.StatusOK, "new", struct {
		Error   error
		Contact contacts.Contact
	}{
		nil,
		contacts.Contact{},
	})
}

// TODO: SOME WAY TO FLASH THE MESSAGE
func PostNewContact(c echo.Context) error {
	var newContact contacts.Contact
	newContact.Id = -1
	newContact.Name = c.FormValue("name")
	newContact.Email = c.FormValue("email")
	newContact.Phone = c.FormValue("phone")

	err := contacts.CreateContact(newContact)
	if err != nil {
		return c.Render(http.StatusOK, "new", struct {
			Error   error
			Contact contacts.Contact
		}{
			err,
			newContact,
		})
	}
	return c.Redirect(http.StatusSeeOther, "/contacts")
}

func GetContactView(c echo.Context) error {
	return c.String(http.StatusOK, c.Param("id"))
}

func GetContactEdit(c echo.Context) error {
	contactId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	contact, err := contacts.GetContact(contactId)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	return c.Render(http.StatusOK, "edit",
		struct {
			Contact contacts.Contact
			Error   error
		}{
			Contact: contact,
			Error:   nil,
		})
}

func PostContactEdit(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	name := c.FormValue("name")
	email := c.FormValue("email")
	phone := c.FormValue("phone")

	err = contacts.UpdateContact(id, name, email, phone)
	if err != nil {
		c.Render(http.StatusBadRequest, "edit", struct {
			Error   error
			Contact contacts.Contact
		}{
			Error: err,
			Contact: contacts.Contact{
				Id:    id,
				Name:  name,
				Email: email,
				Phone: phone,
			},
		})
	}
	return c.Redirect(http.StatusSeeOther, "/contacts")
}

func PostContactDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid id")
	}
	err = contacts.DeleteContact(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.Redirect(http.StatusSeeOther, "/contacts")
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

	// DONE
	e.GET("/contacts", Contacts)

	e.GET("/contacts/:id", GetContactView)

	// DONE
	e.GET("/contacts/:id/edit", GetContactEdit)
	e.POST("/contacts/:id/edit", PostContactEdit)

	// DONE
	e.POST("/contacts/:id/delete", PostContactDelete)

	// DONE
	e.GET("/contacts/new", GetNewContact)
	e.POST("/contacts/new", PostNewContact)

	e.Logger.Fatal(e.Start(":1323"))
}
