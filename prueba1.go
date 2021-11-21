package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
)

var localhostReg string
var localhostNot string
var localhostHp string
var remotehost string
var bitacoraAddr []string
var bitacoraAddr2 []string
var resp string = ""
var bitacoraResp []string

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

func manejadorHP(con net.Conn) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	importaciones, _ := bufferIn.ReadString('\n')
	importaciones = strings.TrimSpace(importaciones)
	resp_nodo, _ := bufferIn.ReadString('\n')
	resp_nodo = strings.TrimSpace(resp_nodo)

	if resp_nodo != "" {
		bitacoraResp = append(bitacoraResp, resp_nodo)
	}
	fmt.Println("Respuesta recibida: ", resp_nodo)
	fmt.Println("Todas las resp: ", bitacoraResp)
	if resp == "" {
		resp = algoritmo(importaciones)
		enviarProximo(importaciones)
	}
}
func algoritmo(importaciones string) string {
	//fmt.Println(importaciones, "Fin importaciones")
	return "respuesta de esta cosa" + localhostReg
}

func enviarProximo(importaciones string) {
	indice := rand.Intn(len(bitacoraAddr2))
	con, _ := net.Dial("tcp", bitacoraAddr2[indice])
	fmt.Printf("Enviando hacia %s", bitacoraAddr2[indice], importaciones)
	defer con.Close()
	fmt.Fprintln(con, importaciones)
	fmt.Fprintln(con, resp)

}
