package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//go:embed index.html
var indexHTML embed.FS

func main() {
	a := app.New()
	w := a.NewWindow("6s2l IoT Device Dashboard")

	t, err := template.New("index").ParseFS(indexHTML, "index.html")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(indexHTML)))

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	w.Resize(fyne.NewSize(400, 300))
	w.SetContent(fyne.NewContainer(
		fyne.NewWidgetRenderer(fyne.NewRasterizer(func(w fyne.CanvasObject, c fyne.Canvas) fyne.Widget {
			return fyne.NewWidget(func() fyne.Widget {
				return templatePage(t)
			})
		})),
	))

	w.ShowAndRun()
}

func templatePage(t *template.Template) fyne.CanvasObject {
	data := struct{}{
		Title: "6s2l IoT Device Dashboard",
		Devices: []struct{}{
			{
				Name:  "Living Room Sensor",
				Status: "Online",
			},
			{
				Name:  "Kitchen Sensor",
				Status: "Offline",
			},
		},
	}

	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		log.Fatal(err)
	}

	return fyne.NewWidgetRenderer(fyne.NewRasterizer(func(w fyne.CanvasObject, c fyne.Canvas) fyne.Widget {
		return widget.NewHTML(string(b.Bytes()))
	}))
}