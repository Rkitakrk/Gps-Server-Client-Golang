package main

import (
	"net/http"

	// "github.com/bmizerany/pat"
	"github.com/gorilla/mux"
	// "github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
	rau := alice.New(app.requireAuthenticatedUser)

	router := mux.NewRouter()
	router.Handle("/", dynamicMiddleware.Then(http.HandlerFunc(app.homePage))).Methods("GET")
	// router.HandleFunc("/", app.homePage).Methods("GET")

	router.Handle("/post/create", dynamicMiddleware.Then(rau.Then(http.HandlerFunc(app.createPostPageForm)))).Methods("GET")
	router.Handle("/post/create", dynamicMiddleware.Then(rau.Then(http.HandlerFunc(app.createPostPage)))).Methods("POST")
	router.Handle("/post/{id}", dynamicMiddleware.Then(http.HandlerFunc(app.showPostPage))).Methods("GET")

	router.Handle("/user/signup", dynamicMiddleware.Then(http.HandlerFunc(app.signupUserForm))).Methods("GET")
	router.Handle("/user/signup", dynamicMiddleware.Then(http.HandlerFunc(app.signupUser))).Methods("POST")
	router.Handle("/user/login", dynamicMiddleware.Then(http.HandlerFunc(app.loginUserForm))).Methods("GET")
	router.Handle("/user/login", dynamicMiddleware.Then(http.HandlerFunc(app.loginUser))).Methods("POST")
	router.Handle("/user/logout", dynamicMiddleware.Then(rau.Then(http.HandlerFunc(app.logoutUser)))).Methods("POST")

	//websocket for search pozition
	router.Handle("/searchGps", dynamicMiddleware.Then(http.HandlerFunc(app.searchGps))).Methods("GET")
	router.Handle("/wsSearchGps", dynamicMiddleware.Then(http.HandlerFunc(app.wsSearchGps)))

	router.Handle("/maps", dynamicMiddleware.Then(http.HandlerFunc(app.maps))).Methods("GET")

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer)).Methods("GET")

	return standardMiddleware.Then(router)

}
