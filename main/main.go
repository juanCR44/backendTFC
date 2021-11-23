package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

func handleRequests(csv [][]string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/sintoma", sintomas(csv))

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func sintomas(csv [][]string) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if req.Method == "POST" {
			body, err := ioutil.ReadAll(req.Body)

			sintomas := string(body)

			if err != nil {
				http.Error(resp, "Error", http.StatusBadRequest)
			}
			conexionCluster(csv, sintomas)
			// luego de esta conexion debe eliminarse los [] para que quede: sospechoso, sospechoso, sospechoso,
			// luego se cuentan cual tiene mas o sospechoso o no sospechoso y se manda sospechoso o no sospechoso
			// abc = resp (sos o no sos)
			var abc = "Hola"
			resp.Header().Set("Content-Type", "application/json")
			io.WriteString(resp, `{"response":"`+abc+`"}`)
		}

	}
}

func conexionCluster(csv [][]string, sintomas string) {
	remotehost := "localhost:9003"
	con, _ := net.Dial("tcp", remotehost)
	defer con.Close()
	fmt.Fprintln(con, csv)
	fmt.Fprintln(con, sintomas)
	fmt.Fprintln(con, "")
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
	handleRequests(pacientes)
}

func main() {
	getDataTraining()
}
