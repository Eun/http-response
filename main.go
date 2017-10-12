package main

import (
    "bytes"
    "html/template"
    "io"
    "log"
    "net/http"
    "fmt"
    "time"
)

var templates = template.Must(template.ParseFiles("index.html", "work_header.html", "work_footer.html"))

func indexHandler(w http.ResponseWriter, req *http.Request) {
    buffer := &bytes.Buffer{}
    if err := templates.ExecuteTemplate(buffer, "index.html", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    } else {
        io.Copy(w, buffer)
    }
}

func workHandler(w http.ResponseWriter, req *http.Request) {

    var flusher http.Flusher
    if f, ok := w.(http.Flusher); ok {
        flusher = f
    }

    bufferHeader := &bytes.Buffer{}
    bufferFooter := &bytes.Buffer{}
    if err := templates.ExecuteTemplate(bufferHeader, "work_header.html", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := templates.ExecuteTemplate(bufferFooter, "work_footer.html", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)

    io.Copy(w, bufferHeader)
    if flusher != nil {
        flusher.Flush()
    }
    
    // actual work
    max := 10
    for i := 1; i <= max; i++ {
        time.Sleep(1000 * time.Millisecond)
        io.WriteString(w, fmt.Sprintf(`<i style="width:%d%%"></i>`, 100*i/max))
        if flusher != nil {
            flusher.Flush()
        }
    }
    io.Copy(w, bufferFooter) 
    if flusher != nil {
        flusher.Flush()
    }

}

func main() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/work", workHandler)
    log.Println("Listening on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}