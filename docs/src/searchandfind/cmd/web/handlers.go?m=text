package main

import (
	"fmt"
	"log"
	"net/http"
	"searchandfind/pkg/forms"
	"searchandfind/pkg/models"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	gubrak "github.com/novalagung/gubrak/v2"
)

type M map[string]interface{}

const MESSAGE_NEW_USER = "New User"
const MESSAGE_CHAT = "Chat"
const MESSAGE_LEAVE = "Leave"

var connections = make([]*WebSocketConnection, 0)

type SocketPayload struct {
	Message string
}

type SocketResponse struct {
	From    string
	Type    string
	Message string
}

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

func (app *application) homePage(w http.ResponseWriter, r *http.Request) {
	// panic("oops! Something went wrong!")

	posts, err := app.posts.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{
		Posts: posts,
	})
}

func (app *application) showPostPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// fmt.Println(vars["id"])
	// var userid = r.URL.Query().Get(":id")
	// fmt.Println(userid)

	var userid = vars["id"]
	post, err := app.posts.Get(userid)
	// fmt.Println(post.Title)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	// fmt.Println(app.authenicatedUser(r))
	app.render(w, r, "show.page.html", &templateData{
		Post: post,
	})
}

func (app *application) createPostPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("title", "content")
	form.MaxLength("title", 100)

	if !form.Valid() {
		app.render(w, r, "create.page.html", &templateData{Form: form})
		return
	}

	var post models.Post

	post.Title = form.Get("title")
	post.Content = form.Get("content")
	post.Created = time.Now()

	id, err := app.posts.Insert(&post)
	if err != nil {
		app.serverError(w, err)
	}

	app.session.Put(r, "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/post/"+id), http.StatusSeeOther)
}

func (app *application) createPostPageForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", &templateData{
		Form: forms.New(nil),
	})

}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)
	form.Required("name", "surname", "email", "password")
	form.MatchesPattern("email", forms.EmailRx)
	form.MinLength("password", 10)

	existsEmail := app.users.Exists(form.Get("email"))
	if existsEmail != nil {
		form.ExistsEmail("email")
	}

	if !form.Valid() {
		app.render(w, r, "signup.page.html", &templateData{
			Form: form,
		})
		return
	}

	var user models.User

	user.Name = form.Get("name")
	user.Surname = form.Get("surname")
	user.Email = form.Get("email")
	password := form.Get("password")
	user.Created = time.Now()

	err = app.users.Insert(&user, password)
	if err != nil {
		app.serverError(w, err)
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("11")
	err := r.ParseForm()
	fmt.Println("1")
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		// fmt.Println("I'm here!")
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.html", &templateData{Form: form})
		return
	}
	// fmt.Println(id)
	app.session.Put(r, "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	// fmt.Fprintln(w, "This is loginUser")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "OK")
	// fmt.Println("Nothing!")

	app.session.Remove(r, "userID")

	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", 303)
	// fmt.Fprintln(w, "This is logoutUser")
}

// Websocket GPS

func (app *application) searchGps(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "searchgps.page.html", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) wsSearchGps(w http.ResponseWriter, r *http.Request) {
	currentGorillaConn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	username := r.URL.Query().Get("username")
	currentConn := WebSocketConnection{Conn: currentGorillaConn, Username: username}
	connections = append(connections, &currentConn)

	go handleIO(&currentConn, connections)

}

func handleIO(currentConn *WebSocketConnection, connections []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	broadcastMessage(currentConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				broadcastMessage(currentConn, MESSAGE_LEAVE, "")
				ejectConnection(currentConn)
				return
			}

			log.Println("ERROR", err.Error())
			continue
		}

		broadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
	}
}
func ejectConnection(currentConn *WebSocketConnection) {
	filtered := gubrak.From(connections).Reject(func(each *WebSocketConnection) bool {
		return each == currentConn
	}).Result()
	connections = filtered.([]*WebSocketConnection)
}

func broadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	for _, eachConn := range connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})
	}
}

func (app *application) maps(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "maps.page.html", &templateData{
		Form: forms.New(nil),
	})
}
