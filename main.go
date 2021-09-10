package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const logging bool = true
const fileName = "test"

// Device is the holding information for all PACP's in the system
type Device struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	ProcessorNum string `json:"pronum"`
	IPAddress    string `json:"ip"`
	Subnet       string `json:"sub"`
	Gateway      string `json:"gateway"`
}

//Sys1 is all of system 1's Devices
var (
	sysDisplay = make([]*Device, 0, 10)
	Sys1       = make([]*Device, 0, 10) // System 1 Devices
	Sys2       = make([]*Device, 0, 10) // System 2 Devices
	Sys3       = make([]*Device, 0, 0)  // System 3 Devices
)

var (
	ip [4]int
)

//Systems is a full map of all systems
var Systems = make(map[string]interface{})

func dummyData() { // this is to create some dummy data when the server first starts
	PACP1 := &Device{
		Name:         "Tom",
		Type:         "PACP",
		ProcessorNum: "1",
		IPAddress:    "10.101.1.216",
		Subnet:       "255.255.255.0",
		Gateway:      "10.101.1.1",
	}
	PACP2 := &Device{
		Name:         "Tom 2",
		Type:         "PACP",
		ProcessorNum: "2",
		IPAddress:    "10.101.1.243",
		Subnet:       "255.255.255.0",
		Gateway:      "10.101.1.1",
	}
	PACP3 := &Device{
		Name:         "John",
		Type:         "MSC 1",
		ProcessorNum: "1",
		IPAddress:    "10.101.1.242",
		Subnet:       "255.255.255.0",
		Gateway:      "10.101.1.1",
	}
	PACP4 := &Device{
		Name:         "Bob",
		Type:         "MSC 4",
		ProcessorNum: "1",
		IPAddress:    "10.101.1.95",
		Subnet:       "255.255.255.0",
		Gateway:      "10.101.1.1",
	}
	//fmt.Printf("PACP1 Addr:%p\n", PACP1)
	//fmt.Printf("PACP2 Addr:%p\n", PACP2)
	//fmt.Printf("PACP3 Addr:%p\n", PACP3)
	//fmt.Printf("PACP4 Addr:%p\n", PACP4)
	//sysexample := []*Device{PACP1, PACP2}
	Sys1 = append(Sys1, PACP1, PACP2)
	Sys2 = append(Sys2, PACP3)
	Sys3 = append(Sys3, PACP4)

	/*for _, Device := range sysexample {
		fmt.Printf("NAme %s", Device.Name)
		fmt.Printf("Addr: %p\n", Device)
	}
	for _, Device := range Sys1 {
		fmt.Printf("NAme %s", Device.Name)
		fmt.Printf("Addr: %p\n", Device)
	}
	for _, Device := range Sys2 {
		fmt.Printf("NAme %s", Device.Name)
		fmt.Printf("Addr: %p\n", Device)
	}
	for _, Device := range Sys3 {
		fmt.Printf("NAme %s", Device.Name)
		fmt.Printf("Addr: %p\n", Device)
	}*/
}

func main() {
	dummyData()
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/systems", systemHandler)
	r.HandleFunc("/project", projectHandler)
	r.HandleFunc("/new", newDevice)
	r.HandleFunc("/edit", editDevice)
	r.HandleFunc("/delete", deleteDevice)
	r.HandleFunc("/update", updateDevice)
	r.HandleFunc("/calc", multiAddrCalc)
	//r.HandleFunc("/save", save)
	r.HandleFunc("/system/{sysnum}", system)
	http.ListenAndServe(":8080", nil)

}

func homeHandler(w http.ResponseWriter, r *http.Request) { // Index Page
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, "unable to load template")
	}
	t.Execute(w, "Home")
}

func projectHandler(w http.ResponseWriter, r *http.Request) { // Project page
	t, err := template.ParseFiles("templates/projectsetup.html")
	if err != nil {
		fmt.Fprintf(w, "unable to load template")
	}
	t.Execute(w, "Project")
}

func systemHandler(w http.ResponseWriter, r *http.Request) { // Systems Page
	Systems["S1"] = Sys1
	Systems["S2"] = Sys2
	Systems["S3"] = Sys3
	if logging == true {
		fmt.Println(Systems)
	}

	t, err := template.ParseFiles("templates/allSystems.html")
	if err != nil {
		fmt.Fprintf(w, "unable to load template")
	}
	t.Execute(w, Systems)
}

func system(w http.ResponseWriter, r *http.Request) { // Show system dependign on system number selected.
	t, err := template.ParseFiles("templates/indvidualSystem.html")
	if err != nil {
		fmt.Fprintf(w, "unable to load template")
	}
	vars := mux.Vars(r)
	sys, err := strconv.Atoi(vars["sysnum"])
	if err != nil {
		fmt.Println("Could not Convert String")
	}

	switch {
	case sys == 1:
		t.Execute(w, Sys1)
	case sys == 2:
		t.Execute(w, Sys2)
	case sys == 3:
		t.Execute(w, Sys3)
	}
	if r.Method == "EXPORT" {
		downLoad(sys)
	}
}

func newDevice(w http.ResponseWriter, r *http.Request) { // New Device
	t, err := template.ParseFiles("templates/newDevice.html")
	if err != nil {
		fmt.Fprintf(w, "unable to load template")
	}
	if r.Method == "POST" {
		if logging == true {
			fmt.Println("Right Function")
		}
		sysNum := r.FormValue("SNum")
		newDevice := &Device{
			Name:         r.FormValue("Nam"),
			Type:         r.FormValue("Typ"),
			ProcessorNum: r.FormValue("PNum"),
			IPAddress:    r.FormValue("IPAddr"),
			Subnet:       r.FormValue("Snet"),
			Gateway:      r.FormValue("Gty"),
		}
		switch {
		case sysNum == "1":
			Sys1 = append(Sys1, newDevice)
			if logging == true {
				log.Println("New Device Added to System 1")
			}
			http.Redirect(w, r, "/system/1", http.StatusSeeOther)
		case sysNum == "2":
			Sys2 = append(Sys2, newDevice)
			if logging == true {
				log.Println("New Device Added to System 2")
			}
			http.Redirect(w, r, "/system/2", http.StatusSeeOther)
		case sysNum == "3":
			Sys3 = append(Sys3, newDevice)
			if logging == true {
				log.Println("New Device Added to System 3")
			}
			http.Redirect(w, r, "/system/3", http.StatusSeeOther)
		}
	} else {
		fmt.Println("Page Load")
	}
	t.Execute(w, nil)
}

func editDevice(w http.ResponseWriter, r *http.Request) { // Edit Device
	t, err := template.ParseFiles("templates/editDevice.html")
	if err != nil {
		fmt.Fprintf(w, "unable to load template")
	}
	t.Execute(w, "Edit Device")
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {

}

func updateDevice(w http.ResponseWriter, r *http.Request) {

}

func downLoad(x int) {

}

func multiAddrCalc(w http.ResponseWriter, r *http.Request) { // Index Page
	t, err := template.ParseFiles("templates/MultiCastAddrCalc.html")
	if err != nil {
		fmt.Fprintf(w, "unable to load template")
	}
	ip[0] = 239
	ip[1] = 255
	fmt.Println(ip)
	if r.Method == "POST" {
		if logging == true {
			fmt.Println("Right Function")
		}
		uniNum, err := strconv.Atoi(r.FormValue("uni"))
		if err != nil {
			fmt.Println("Could not Convert String")
		}
		x, y := uniCal(uniNum)
		ip[2] = x
		ip[3] = y
	}
	fmt.Println(ip)
	t.Execute(w, ip)
}

func uniCal(x int) (int, int) {
	a := x / 255
	b := x%255 - a
	return a, b
}
