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

func handleRequests(csv [][]string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", cargarImportaciones)
	mux.HandleFunc("/importaciones", mostrarImportaciones)
	mux.HandleFunc("/sintoma", sintomas(csv))

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func mostrarImportaciones(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(importaciones, "", " ")
	io.WriteString(resp, string(jsonBytes))
}

func sintomas(csv [][]string) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		body, err := ioutil.ReadAll(req.Body)

		var sintomas []string
		json.Unmarshal([]byte(body), &sintomas)

		if err != nil {
			http.Error(resp, "Error", http.StatusBadRequest)
		}

		conexionCluster(csv, sintomas)
		resp.Header().Set("Content-Type", "application/json")

		jsonBytes, _ := json.MarshalIndent(csv, "", " ")
		io.WriteString(resp, string(jsonBytes))
	}
}

func conexionCluster(csv [][]string, sintomas []string) {
	remotehost := "localhost:9003"
	con, err := net.Dial("tcp", remotehost)
	fmt.Print(err, " ERRORORO")
	defer con.Close()
	fmt.Fprintln(con, csv)
	fmt.Fprintln(con, sintomas)
	fmt.Fprintln(con, "")
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
		/*conexionCluster(importaciones)
		resp.Header().Set("Content-Type", "application/json")

		jsonBytes, _ := json.MarshalIndent(importaciones, "", " ")
		io.WriteString(resp, string(jsonBytes))

		_ = likely
		_ = scores

		//Print del resultado
		if probs[0] > probs[1] {
			fmt.Print("Sospechoso")
		} else {
			fmt.Print("No sospechoso")
		}
		fmt.Print(probs)
		*/

	} else {
		http.Error(resp, "Método inválido", http.StatusMethodNotAllowed)
	}

}

func getDataTraining() {
	url := "https://raw.githubusercontent.com/juanCR44/backendTFC/main/COVID_DF.csv"
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var pacientes [][]string
	csvReader := csv.NewReader(resp.Body)
	for {
		paciente, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//Eliminando 0's de las listas de sintomas
		aux := 0
		for _, n := range paciente {
			if n != "0" {
				paciente[aux] = n
				aux++
			}
		}

		//Generación de lista de filas
		paciente = paciente[:aux]
		pacientes = append(pacientes, paciente)
	}
	/*

		//Creación del clasificador bayesiano
		classifier := bayesian.NewClassifier(sospechoso, no_sospechoso)

		//Entrenamiento con la data del csv
		for i := 0; i < len(pacientes); i++ {
			if pacientes[i][0] == "Flag_sospechoso" {
				sospechosoSintomas := pacientes[i][1:]
				classifier.Learn(sospechosoSintomas, sospechoso)
			} else {
				no_sospechosoSintomas := pacientes[i]
				classifier.Learn(no_sospechosoSintomas, no_sospechoso)
			}
		}

		fmt.Print("Finalizado entrenamiento")*/
	handleRequests(pacientes)
}

func main() {
	getDataTraining()
}
