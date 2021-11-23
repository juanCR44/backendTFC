package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"

	"github.com/navossoc/bayesian"
)

var localhostReg string
var localhostNot string
var localhostHp string
var remotehost string
var bitacoraAddr []string
var bitacoraAddr2 []string
var resp string = ""
var bitacoraResp []string
var csv [][]string
var sint []string
var respuesta_bitacora []string

const (
	sospechoso    bayesian.Class = "sospechoso"
	no_sospechoso bayesian.Class = "no_sospechoso"
)

func main() {
	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto de registro: ")
	port, _ := bufferIn.ReadString('\n')
	port = strings.TrimSpace(port)
	localhostReg = fmt.Sprintf("localhost:%s", port)

	fmt.Print("Ingrese el puerto de notificacion: ")
	port, _ = bufferIn.ReadString('\n')
	port = strings.TrimSpace(port)
	localhostNot = fmt.Sprintf("localhost:%s", port)

	fmt.Print("Ingrese el puerto de proceso HP: ")
	port, _ = bufferIn.ReadString('\n')
	port = strings.TrimSpace(port)
	localhostHp = fmt.Sprintf("localhost:%s", port)

	go registrarServer()
	go servicioHP()

	fmt.Print("Ingrese puerto del nodo a solicitar registro: ")
	puerto, _ := bufferIn.ReadString('\n')
	puerto = strings.TrimSpace(puerto)
	remotehost = fmt.Sprintf("localhost:%s", puerto)

	if puerto != "" {
		registrarSolicitud(remotehost)
	}
	recibeNotificarServer()
}

func registrarServer() {
	ln, _ := net.Listen("tcp", localhostReg)
	defer ln.Close()

	for {
		con, _ := ln.Accept()
		go manejadorRegistro(con)
	}

}

func manejadorRegistro(con net.Conn) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	ip, _ := bufferIn.ReadString('\n')
	ip = strings.TrimSpace(ip)

	bytes, _ := json.Marshal(append(bitacoraAddr, localhostNot))
	fmt.Fprintln(con, string(bytes)) //envia la bitácora

	ip2, _ := bufferIn.ReadString('\n')
	ip2 = strings.TrimSpace(ip2) ///localhost:puerto

	fmt.Println("IP2:", ip2)

	bytes, _ = json.Marshal(append(bitacoraAddr2, localhostHp))
	fmt.Fprintln(con, string(bytes))
	comunicarTodos(ip, ip2)

	bitacoraAddr = append(bitacoraAddr, ip)
	bitacoraAddr2 = append(bitacoraAddr2, ip2)

	fmt.Println(bitacoraAddr)
	fmt.Println(bitacoraAddr2)
}

func comunicarTodos(ip, ip2 string) {
	for _, addr := range bitacoraAddr {
		notificar(addr, ip, ip2)
	}
}

func notificar(addr, ip, ip2 string) {
	con, _ := net.Dial("tcp", addr)
	defer con.Close()
	fmt.Fprintln(con, ip)
	fmt.Fprintln(con, ip2)
}

func registrarSolicitud(remotehost string) {
	con, _ := net.Dial("tcp", remotehost)
	defer con.Close()
	fmt.Fprintln(con, localhostNot) //enviamos el puerto de notificacion

	//recuperar lo que responde el server
	bufferIn := bufio.NewReader(con)
	bitacoraServer, _ := bufferIn.ReadString('\n')

	var bitacoraTemp []string
	json.Unmarshal([]byte(bitacoraServer), &bitacoraTemp)

	bitacoraAddr = bitacoraTemp

	fmt.Fprintln(con, localhostHp) //enviamos el puerto de notificacion

	//recuperar lo que responde el server
	bitacoraServer, _ = bufferIn.ReadString('\n')

	var bitacoraTemp2 []string
	json.Unmarshal([]byte(bitacoraServer), &bitacoraTemp2)

	bitacoraAddr2 = bitacoraTemp2

	fmt.Println(bitacoraAddr)
	fmt.Println(bitacoraAddr2)
}

func recibeNotificarServer() {
	ln, _ := net.Listen("tcp", localhostNot)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejadorRecibeNotificar(con)
	}
}

func manejadorRecibeNotificar(con net.Conn) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	ip, _ := bufferIn.ReadString('\n')
	ip = strings.TrimSpace(ip)
	bitacoraAddr = append(bitacoraAddr, ip)

	ip2, _ := bufferIn.ReadString('\n')
	ip2 = strings.TrimSpace(ip2)
	bitacoraAddr2 = append(bitacoraAddr2, ip2)

	fmt.Println(bitacoraAddr)
	fmt.Println(bitacoraAddr2)
}

////////////////////////////

func servicioHP() {
	ln, _ := net.Listen("tcp", localhostHp)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejadorHP(con)
	}
}

// func convertir_arra_string(csv string){
// 	tam := len(csv)
// 	csv = csv[1 : tam-1]

// 	arre_uni := strings.Split(csv, "]")

// 	for i,s := range arre_uni {

// 	}
// }

func convertir_string_to_array2(test string) [][]string {
	var listaT [][]string
	test = test[1 : len(test)-1]
	var s string
	for i := 0; i < len(test); i++ {
		if test[i] == ' ' && test[i-1] == ']' {
			s = s[1 : len(s)-1]
			xd := strings.Fields(s)
			listaT = append(listaT, xd)
			s = ""
		} else {
			s = s + string(test[i])
		}
	}
	return listaT
}

func convertir_string_to_array1(test string) []string {
	test = test[1 : len(test)-1]
	x := strings.Split(test, ",")
	for i := 0; i < len(x); i++ {
		x[i] = x[i][1 : len(x[i])-1]
	}
	return x
}

func manejadorHP(con net.Conn) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)

	csv1, _ := bufferIn.ReadString('\n')
	csv1 = strings.TrimSpace(csv1)
	csv = convertir_string_to_array2(csv1)

	sintomas, _ := bufferIn.ReadString('\n')
	sintomas = strings.TrimSpace(sintomas)
	sint = convertir_string_to_array1(sintomas)

	resp_nodo, _ := bufferIn.ReadString('\n')
	resp_nodo = strings.TrimSpace(resp_nodo)

	if resp_nodo != "" {
		// conversion
		fmt.Println("aca vuela")
		resp_nodo = resp_nodo[1 : len(resp_nodo)-1]
		respuesta_bitacora = strings.Fields(resp_nodo)
		for _, s := range respuesta_bitacora {
			bitacoraResp = append(bitacoraResp, s)
		}
	}

	//fmt.Println("Todas las resp: ", bitacoraResp)

	if resp == "" {
		algoritmo(sint)
		if len(bitacoraResp) < 3 {
			enviarProximo()
		} else {
			//fmt.Println("Todos los resultados:", bitacoraResp)
			// conecto con el api
			enviarApi()
		}
	}
}

func algoritmo(sintomas []string) {
	classifier := bayesian.NewClassifier(sospechoso, no_sospechoso)

	for i := 0; i < len(csv); i++ {
		if csv[i][0] == "Flag_sospechoso" {
			sospechosoSintomas := csv[i][1:]
			classifier.Learn(sospechosoSintomas, sospechoso)
		} else {
			no_sospechosoSintomas := csv[i]
			classifier.Learn(no_sospechosoSintomas, no_sospechoso)
		}
	}

	//Aquí se pasan la lista de sintomas que ingresa el usuario
	scores, likely, _ := classifier.LogScores(
		sintomas,
	)
	probs, likely, _ := classifier.ProbScores(
		sintomas,
	)

	_ = likely
	_ = scores

	//Print del resultado
	if probs[0] > probs[1] {
		fmt.Print("Sospechoso")
		resp = "Sospechoso"
	} else {
		fmt.Print("No sospechoso")
		resp = "No Sospechoso"
	}
	fmt.Print(probs)
	bitacoraResp = append(bitacoraResp, resp)
}

func enviarProximo() {
	if len(bitacoraAddr2) > 1 {
		indice := rand.Intn(len(bitacoraAddr2))
		con, _ := net.Dial("tcp", bitacoraAddr2[indice])
		defer con.Close()
		fmt.Fprintln(con, csv)
		fmt.Fprintln(con, sint)
		fmt.Fprintln(con, bitacoraResp)
	}
}

func enviarApi() {
	//bitacoraResp
	// esto deberia devolver al main.go bitacora Resp que es basicamente esto: [Sospechoso Sospechoso Sospechoso]
	// una vez que lo pase se cambia por "abc" que se esta enviando ahora en el main al front
}
