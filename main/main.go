package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

var respuesta_cluster []string
var respuesta_api string

func handleRequests(csv [][]string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/sintoma", sintomas(csv))

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
func convertir_string_to_array(text string) []string {
	text = text[1 : len(text)-1]
	respuesta_cluster = strings.Fields(text)
	return respuesta_cluster
}
func retornar_resultado() {
	var sos = 0
	var no_sos = 0

	for i, s := range respuesta_cluster {
		s = strings.TrimSpace(s)
		fmt.Println(i, s)
		if s == "Sospechoso" {
			sos += 1
		} else {
			no_sos += 1
		}
	}

	if sos > no_sos {
		respuesta_api = "Sospechoso"
	} else {
		respuesta_api = "No Sospechoso"
	}
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
			retornar_resultado()
			fmt.Println(respuesta_api, "esto envia")
			resp.Header().Set("Content-Type", "application/json")
			io.WriteString(resp, `{"response":"`+respuesta_api+`"}`)
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

	// recibir respuesta
	ln, err := net.Listen("tcp", "localhost:9011")
	if err != nil {
		fmt.Printf("Error en la resolución de la dirección de red!! ", err.Error())
	}
	defer ln.Close()
	con2, err := ln.Accept()
	if err != nil {
		fmt.Printf("Error en la conexion!!!: ", err.Error())
	}
	defer con2.Close()
	bufferIn := bufio.NewReader(con2)
	msg, _ := bufferIn.ReadString('\n')
	msg = strings.TrimSpace(msg)
	println("Respuesta del cluster: ", msg)
	if msg != "" {

		convertir_string_to_array(msg)
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
	handleRequests(pacientes)
}

func main() {
	getDataTraining()
}
