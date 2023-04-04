package handler

import (
	Repository "L0/pkg/repository"
	"fmt"
	"html/template"
	"net/http"
)

type Handler struct {
	Repo Repository.Repository
}

func NewHandler(repo Repository.Repository) Handler {
	var handler Handler
	handler.Repo = repo

	return handler
}

func (h *Handler) orderHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	uid := r.FormValue("uid")
	order, _ := h.Repo.GetOrder(uid)

	template, _ := template.ParseFiles("static/order.html")
	template.Execute(w, order)
}

func (h *Handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "static/")
}

func (h *Handler) InitRouters() {
	fmt.Println("[LOG]: Initializing routers")
	http.HandleFunc("/", h.indexHandler)
	http.HandleFunc("/order", h.orderHandler)
}
