package frontend

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/monishth/go-therm/internal/core"
	"github.com/monishth/go-therm/internal/models"
)

type ZoneData struct {
	Zone       models.Zone
	TargetTemp float64
}

func StartFrontend(app *core.App, ctx context.Context) {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", wrapHandler(app, handleIndex))
	http.HandleFunc("/zone/", wrapHandler(app, handleZone))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Server started at http://localhost:8080")

	<-ctx.Done()

	log.Println("Shutting down frontend server")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

func wrapHandler(app *core.App, h func(*core.App, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(app, w, r)
	}
}

func handleIndex(app *core.App, w http.ResponseWriter, r *http.Request) {
	// Parse the templates
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "index.html"),
		filepath.Join("templates", "zone.html"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Combine Zone data with TargetTemp for each Zone
	var zoneData []ZoneData
	for _, zone := range app.Entities.Zones {
		zoneData = append(zoneData, ZoneData{
			Zone:       zone,
			TargetTemp: app.GetTargets()[zone.ID],
		})
	}

	data := struct {
		Zones []ZoneData
	}{
		Zones: zoneData,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleZone(app *core.App, w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/zone/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid zone ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		for _, zone := range app.Entities.Zones {
			if zone.ID == id {
				tmpl, err := template.ParseFiles(filepath.Join("templates", "zone.html"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				data := ZoneData{
					Zone:       zone,
					TargetTemp: app.GetTargets()[zone.ID],
				}
				if err := tmpl.ExecuteTemplate(w, "zone", data); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		http.Error(w, "Zone not found", http.StatusNotFound)
	case http.MethodPut:
		targetTemp, err := strconv.ParseFloat(r.FormValue("temp"), 64)
		if err != nil {
			http.Error(w, "Invalid temperature value", http.StatusBadRequest)
			return
		}

		for _, zone := range app.Entities.Zones {
			if zone.ID == id {
				app.SetTarget(zone.ID, targetTemp)
				tmpl, err := template.ParseFiles(filepath.Join("templates", "zone.html"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				data := ZoneData{
					Zone:       zone,
					TargetTemp: app.GetTargets()[zone.ID],
				}
				if err := tmpl.ExecuteTemplate(w, "zone", data); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		http.Error(w, "Zone not found", http.StatusNotFound)
	}
}
