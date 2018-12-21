// A simple webserver to request an account on the MOC.
// Displays the account request page
//   upon completion, sends json to the MOC Adjutant Service
//   Displays the Thankyou for requesting an Account
//
// This is modeled after these examples:
//    https://golang.org/doc/articles/wiki/
//
// ToDo:
//    1) make https
//    2) template the service list
//    3) add the JSON call
//    4) move the template into the code (no point to having separate files)
//    5) containerize

package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("rbody")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	//username := r.FormValue("username")
	//kaizen_cb := r.FormValue("kaizen")
	//kshift_cb := r.FormValue("kshift")
	//klambda_cb := r.FormValue("klambda")
	//reason := r.FormValue("reason")
	//sponser := r.FormValue("sponser")

	//construct json here to send to adjutant

	http.Redirect(w, r, "/thankyou/", http.StatusFound)
}

var templates = template.Must(template.ParseFiles("AcntReq.html", "ThankYou.html"))

//var template

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func acntReqHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Dummy", Body: []byte("  ")}
	err := templates.ExecuteTemplate(w, "AcntReq.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func thanksHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Dummy", Body: []byte("  ")}
	err := templates.ExecuteTemplate(w, "ThankYou.html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	//http.HandleFunc("/view/", makeHandler(viewHandler))
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))

	http.HandleFunc("/AcntReq/", acntReqHandler)
	http.HandleFunc("/Request/", requestHandler)
	http.HandleFunc("/thankyou/", thanksHandler)

	log.Fatal(http.ListenAndServe(":8443", nil))
}
