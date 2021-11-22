package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type Importacion struct {
	Anio     string `json:"anio"`
	Mes      string `json:"mes"`
	Producto string `json:"producto"`
	Peso     string `json:"peso"`
}

var importaciones []Importacion

func handleRequests() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", cargarImportaciones)
	mux.HandleFunc("/importaciones", mostrarImportaciones)
	mux.HandleFunc("/sintoma", cargarSintomas)

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func mostrarImportaciones(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(importaciones, "", " ")
	io.WriteString(resp, string(jsonBytes))
}

func cargarSintomas(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var listaSintoma []string
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(resp, "Error", http.StatusBadRequest)
	}

	json.Unmarshal([]byte(body), &listaSintoma)
	fmt.Print(listaSintoma)
}

func conexionCluster(importaciones []Importacion) {
	remotehost := "localhost:9003"
	con, err := net.Dial("tcp", remotehost)
	fmt.Print(err, " ERRORORO")
	defer con.Close()
	fmt.Fprintln(con, importaciones)
	fmt.Fprintln(con, "")
}

func cargarSintomas(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var listaSintoma []string
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		http.Error(resp, "Error", http.StatusBadRequest)
	}

	json.Unmarshal([]byte(body), &listaSintoma)
	fmt.Print(listaSintoma)
}

func cargarImportaciones(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	resp.Header().Set("Access-Control-Allow-Credentials", "true")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if req.Method == "POST" {
		fs, _, err := req.FormFile("file")
		if err != nil {
			http.Error(resp, "Error", http.StatusBadRequest)
		}

		fmt.Printf("Correcto")
		if err != nil {
			fmt.Printf("get form err: %s", err)
		}
		defer fs.Close()
		r := csv.NewReader(fs)

		for {
			row, err := r.Read()
			if err != nil && err != io.EOF {
				fmt.Printf("can not read, err is %+v", err)
			}
			if err == io.EOF {
				break
			}
			importacion := Importacion{
				Anio:     row[0],
				Mes:      row[1],
				Producto: row[2],
				Peso:     row[3],
			}

			importaciones = append(importaciones, importacion)
		}

		//fmt.Print(importaciones)
		conexionCluster(importaciones)
		resp.Header().Set("Content-Type", "application/json")

		jsonBytes, _ := json.MarshalIndent(importaciones, "", " ")
		io.WriteString(resp, string(jsonBytes))

	} else {
		http.Error(resp, "Método inválido", http.StatusMethodNotAllowed)
	}

}

func main() {
	handleRequests()
}
