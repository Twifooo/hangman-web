package main

import (
    "html/template"
    "net/http"
    "strconv"
    "regexp"
)

var count int 

type Etudiant struct {
    Nom    string
    Prenom string
    Age    int
    Sexe   string
}

type Classe struct {
    Nom           string
    Filiere       string
    Niveau        string
    NbEtudiants   int
    ListeEtudiants []Etudiant
}

func promoHandler(w http.ResponseWriter, r *http.Request) {
    etudiants := []Etudiant{
        {"Dupont", "Jean", 20, "masculin"},
        {"Martin", "Claire", 19, "féminin"},
        {"Durand", "Alex", 21, "masculin"},
    }

    classe := Classe{
        Nom:          "B1 Informatique",
        Filiere:      "Informatique",
        Niveau:       "Bachelor 1",
        NbEtudiants:  len(etudiants),
        ListeEtudiants: etudiants,
    }

    tmpl := template.Must(template.ParseFiles("templates/promo.html"))
    tmpl.Execute(w, classe)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
    count++
    tmpl := template.Must(template.ParseFiles("templates/change.html"))

    data := struct {
        Count int
        Pair  bool
    }{
        Count: count,
        Pair:  count%2 == 0,
    }

    tmpl.Execute(w, data)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/form.html"))
    tmpl.Execute(w, nil)
}

func treatmentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        nom := r.FormValue("nom")
        prenom := r.FormValue("prenom")
        naissance := r.FormValue("naissance")
        sexe := r.FormValue("sexe")

        validName := regexp.MustCompile(`^[a-zA-Z]{1,32}$`).MatchString
        validSexe := sexe == "masculin" || sexe == "féminin" || sexe == "autre"

        if validName(nom) && validName(prenom) && validSexe {
            http.Redirect(w, r, "/user/display?nom="+nom+"&prenom="+prenom+"&naissance="+naissance+"&sexe="+sexe, http.StatusSeeOther)
        } else {
            http.Error(w, "Données non valides", http.StatusBadRequest)
        }
    }
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
    nom := r.URL.Query().Get("nom")
    prenom := r.URL.Query().Get("prenom")
    naissance := r.URL.Query().Get("naissance")
    sexe := r.URL.Query().Get("sexe")

    data := struct {
        Nom       string
        Prenom    string
        Naissance string
        Sexe      string
    }{
        Nom:       nom,
        Prenom:    prenom,
        Naissance: naissance,
        Sexe:      sexe,
    }

    tmpl := template.Must(template.ParseFiles("templates/display.html"))
    tmpl.Execute(w, data)
}

func main() {
    http.HandleFunc("/promo", promoHandler)

    http.HandleFunc("/change", changeHandler)

    http.HandleFunc("/user/form", formHandler)
    http.HandleFunc("/user/treatment", treatmentHandler)
    http.HandleFunc("/user/display", displayHandler)

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    http.ListenAndServe(":8080", nil)
}
