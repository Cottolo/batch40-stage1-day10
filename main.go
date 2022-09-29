package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/public/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	// route := mux.NewRouter()
	// route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("hello world"))
	// }).Methods("GET")

	connection.DatabaseConnect()

	route := mux.NewRouter()

	// path folder public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	// routing
	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/detail/{id}", detailProject).Methods("GET")

	route.HandleFunc("/project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")

	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")

	route.HandleFunc("/edite-project/{id}", formEditeProject).Methods("GET")
	route.HandleFunc("/edite-project/{id}", editeProject).Methods("POST")

	fmt.Println("server running at localhost:5000")
	http.ListenAndServe("localhost:5000", route)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, error = template.ParseFiles("views/index.html")

	if error != nil {
		w.Write([]byte("not found 404"))
		return
	}

	// response := map[string]interface{}{

	// 	"Projects": dataProject,
	// }

	data, _ := connection.Con.Query(context.Background(),
		"SELECT id ,project_name, description  FROM tb_project")
	// fmt.Println(data)

	var result []Project
	for data.Next() {
		var each = Project{}

		var err = data.Scan(
			&each.Id,
			&each.ProjectName,
			&each.ProjectDescription,
			// &each.StartDate,
			// &each.EndDate,
			// &each.NodeJs,
			// &each.NextJs,
			// &each.ReactJs,
			// &each.TypeScript,
		)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)

	}

	fmt.Println(result)

	resData := map[string]interface{}{
		"Projects": result,
	}

	tmpl.Execute(w, resData)
}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, error = template.ParseFiles("views/project.html")

	if error != nil {
		w.Write([]byte("not found 404"))
		return
	}

	tmpl.Execute(w, nil)
}

// Type Data
type Project struct {
	Id                 int
	ProjectName        string
	StartDate          time.Time
	EndDate            time.Time
	Duration           float64
	ProjectDescription string
	NodeJs             string
	NextJs             string
	ReactJs            string
	TypeScript         string
}

// ARRAY
// var dataProject = []Project{}

func addProject(w http.ResponseWriter, r *http.Request) {
	error := r.ParseForm()
	if error != nil {
		log.Fatal(error)
	}

	var projectName = r.PostForm.Get("project-name")
	var projectDescription = r.PostForm.Get("project-description")
	// var startDate = r.PostForm.Get("start-date")
	// var endDate = r.PostForm.Get("end-date")
	// var nodeJs = r.PostForm.Get("node-js")
	// var nextJs = r.PostForm.Get("next-js")
	// var reactJs = r.PostForm.Get("react-js")
	// var typeScript = r.PostForm.Get("typescript")
	// var layout = "2006-01-02"
	// var start, _ = time.Parse(layout, startDate)
	// var end, _ = time.Parse(layout, endDate)
	// var duration = math.Round(end.Sub(start).Hours() / 24 / 30)
	_, err := connection.Con.Exec(context.Background(),
		"INSERT INTO tb_project (project_name, description) VALUES ($1,$2)",
		projectName, projectDescription)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("suddenly :" + err.Error()))
		return
	}

	//OBJECT
	// var newProject = Project{
	// 	ProjectName:        projectName,
	// 	ProjectDescription: projectDescription,
	// 	StartDate:          startDate,
	// 	EndDate:            endDate,
	// 	NextJs:             nextJs,
	// 	ReactJs:            reactJs,
	// 	NodeJs:             nodeJs,
	// 	TypeScript:         typeScript,
	// 	Duration:           duration,
	// 	Id:                 len(dataProject),
	// }

	//PUSH
	// dataProject = append(dataProject, newProject)

	// fmt.Println(dataProject)

	//HALAMAN SETELAH POST
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, error = template.ParseFiles("views/contact.html")

	if error != nil {
		w.Write([]byte("not found 404"))
		return
	}

	tmpl.Execute(w, nil)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, error = template.ParseFiles("views/detail.html")

	if error != nil {
		w.Write([]byte("not found 404"))
		return
	}

	var ProjectDetail = Project{}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// fmt.Println(id)

	// for i, data := range dataProject {
	// 	if id == i {
	// 		ProjectDetail = Project{
	// 			ProjectName:        data.ProjectName,
	// 			ProjectDescription: data.ProjectDescription,
	// 			StartDate:          data.StartDate,
	// 			EndDate:            data.EndDate,
	// 			NextJs:             data.NextJs,
	// 			ReactJs:            data.ReactJs,
	// 			NodeJs:             data.NodeJs,
	// 			TypeScript:         data.TypeScript,
	// 			Duration:           data.Duration,
	// 		}
	// 	}
	// }

	err := connection.Con.QueryRow(context.Background(), "SELECT id, project_name, description FROM tb_project WHERE id=$1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.ProjectDescription)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ups : " + err.Error()))
	}
	// OBJECT
	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	// fmt.Println(data)

	tmpl.Execute(w, data)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// fmt.Println(id)

	_, err := connection.Con.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("hmmmm : " + err.Error()))
	}
	// dataProject = append(dataProject[:id], dataProject[id+1:]...)

	http.Redirect(w, r, "/home", http.StatusFound)

}

func formEditeProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("views/edite-project.html")

	if err != nil {
		w.Write([]byte("WHAT A PITTY : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	data := map[string]interface{}{
		"EditeId": id,
	}

	tmpl.Execute(w, data)
}

func editeProject(w http.ResponseWriter, r *http.Request) {
	error := r.ParseForm()
	if error != nil {
		log.Fatal(error)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	var projectName = r.PostForm.Get("project-name")
	var projectDescription = r.PostForm.Get("project-description")
	// var startDate = r.PostForm.Get("start-date")
	// var endDate = r.PostForm.Get("end-date")
	// var nodeJs = r.PostForm.Get("node-js")
	// var nextJs = r.PostForm.Get("next-js")
	// var reactJs = r.PostForm.Get("react-js")
	// var typeScript = r.PostForm.Get("typescript")
	// var layout = "2006-01-02"
	// var start, _ = time.Parse(layout, startDate)
	// var end, _ = time.Parse(layout, endDate)
	// var duration = math.Round(end.Sub(start).Hours() / 24 / 30)

	_, err := connection.Con.Exec(context.Background(),
		"UPDATE tb_project SET project_name = $1, description=$2 WHERE id=$3", projectName, projectDescription, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// var editeProject = Project{
	// 	ProjectName:        projectName,
	// 	ProjectDescription: projectDescription,
	// 	StartDate:          startDate,
	// 	EndDate:            endDate,
	// 	NextJs:             nextJs,
	// 	ReactJs:            reactJs,
	// 	NodeJs:             nodeJs,
	// 	TypeScript:         typeScript,
	// 	Duration:           duration,
	// }

	// dataProject[id] = editeProject
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
