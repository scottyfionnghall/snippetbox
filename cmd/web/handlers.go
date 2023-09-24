package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com.scottyfionnghall.snippetbox/internal/models"
	"github.com.scottyfionnghall.snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// Define a snippetCreateForm struct to represent the form data
// and validation errors from fields. All fields must be exported to be used
// by html/templates
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// This handler returns home page.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/".

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.html", data)
}

// This handler shows user particular snippet based on the passed ID.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// httprouter extracts all parameters passed in the request in a form
	// of a slice
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)
}

// This handler handels POST requests to create a new snipppet in the database
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use Validator to check all fields.
	form.CheckField(validator.NotBlank(form.Title), "title",
		"This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title",
		"This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content",
		"This field cannot be blank")
	form.CheckField(validator.PermitedValue(form.Expires, 1, 7, 365), "expires",
		"This field must equal 1,7 or 365")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// Pass the data to the SnippetModel.Insert() method, reciving the
	// ID of the new record back
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Use the Put() method to add a string value and the correspondin key
	// to the session data
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully create!")
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// This handler handels GET requests to show user a form to create a snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

// This handler handels DELETE requests to remove created snippets from database
// by their ID
func (app *application) snippetDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.badRequest(w)
		return
	}
	err = app.snippets.Delete(id)
	if err != nil {
		app.badRequest(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// This handler handels GET requests to show user signup form
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

// This handler handels POST requests to save user info in the database
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	// Parse the form data into the userSignupForm struct
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Validate the form contents using helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name",
		"This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email",
		"This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email",
		"This field must be a valid email adress")
	form.CheckField(validator.NotBlank(form.Password), "password",
		"This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password",
		"This field must be at least 8 characters long")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}
	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// This handler handels GET requests to show user login form
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

// This handler handels POST requests to authinticate and login the user
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// Decode the form data into the userLoginForm struct
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Do validation checks on the form.
	form.CheckField(validator.NotBlank(form.Email), "email",
		"This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email",
		"This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password",
		"This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-dispaly the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Use the RenewToken() method on the current session to change the session ID.
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Add the ID of the current user to the session, so that they are now
	// "logged in"
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

// This handler handels POST requests to logout the user
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change the session
	// ID again
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Remove the authenticatedUserID from the session data so that the user
	// is "logged out"
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	// Add a flash message to the session to confirm to the user they've been
	// logged out
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	// Redirect the user to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
