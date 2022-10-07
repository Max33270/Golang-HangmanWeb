package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"piscine"
	"text/template"
)

type Test struct {
	Att  int
	Word string
	Jose string
	Rep  []string
	Win  []piscine.Score
}

var winners []piscine.Score
var attempt int    // Tentatatives
var UdScore []rune // Mot (_ A _ _ B O)
var pick string    // Mot (L A V A B O)
var boolean = true
var rep []string
var Name string
var level string

func Accueil(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/accueil.html")
}

func Choix(w http.ResponseWriter, r *http.Request) {
	UdScore = []rune{}
	if len(os.Args) == 1 {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		rep = []string{}
		piscine.Repetition = []string{}
		level = r.Form.Get("w")
		pick = piscine.ToUpper(piscine.Random(level))
		Name = r.Form.Get("nom_utilisateur")
		attempt = 10

		//Mot caché (_ _ _ _ _)
		for range pick {
			UdScore = append(UdScore, '_')
		}

		//Premières lettres dévoilées (_ I _ _ E)
		for v := 0; v < len(pick)/2-1; v++ {
			random := rand.Intn(len(pick))
			if UdScore[random] == '_' {
				UdScore[random] = rune(pick[random])
			} else {
				v--
			}
		}
	}
	http.Redirect(w, r, "/", 301)
}

func Redirect(w http.ResponseWriter, r *http.Request) {

	if boolean {
		boolean = false
		http.Redirect(w, r, "/accueil", 301)
	} else {
		new := Test{Att: attempt, Word: string(UdScore), Jose: piscine.Check(attempt), Rep: rep}
		tmpl := template.Must(template.ParseFiles("tmpl/index.html"))
		tmpl.Execute(w, new)
	}
}

func Hangman(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	letter := piscine.ToUpper(r.Form.Get("field2"))
	deja := true
	for i := range rep {
		if rep[i] == letter {
			deja = false
		}
	}

	udd := attempt
	UdScore, attempt = piscine.Compare(UdScore, attempt, pick, letter)
	if deja && udd != attempt {
		rep = append(rep, letter)
	}

	if attempt <= 0 && r.Method == "POST"{
		boolean = true
		if boolean {
			http.Redirect(w, r, "/loose", 301)
		}
	}
	if string(UdScore) == pick && r.Method == "POST" {
		boolean = true
		if boolean {
			http.Redirect(w, r, "/win", 301)
		}
	}
	http.Redirect(w, r, "/", 301)
}

func Loose(w http.ResponseWriter, r *http.Request) {
	new := Test{Word: pick}
	tmpl := template.Must(template.ParseFiles("tmpl/loose.html"))
	tmpl.Execute(w, new)

	//http.ServeFile(w, r, "tmpl/loose.html")

}

func Win(w http.ResponseWriter, r *http.Request) {

	if level == "EASY" {
		winners = append(winners, piscine.Score{Name, attempt})
	} else if level == "NORMAL" {
		winners = append(winners, piscine.Score{Name, attempt * 2})
	} else {
		winners = append(winners, piscine.Score{Name, attempt * 3})
	}
	winners = piscine.ScoreJoueur(winners)
	new := Test{Win: winners}
	tmpl := template.Must(template.ParseFiles("tmpl/win.html"))
	tmpl.Execute(w, new)
}

func main() {

	http.HandleFunc("/", Redirect)
	http.HandleFunc("/accueil", Accueil)
	http.HandleFunc("/win", Win)
	http.HandleFunc("/loose", Loose)
	http.HandleFunc("/hangman", Hangman)
	http.HandleFunc("/choix", Choix)

	//Show #CSS
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Listening at http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)

}
