package regsitry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"
const ServerUrl = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) register(registration Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.registrations = append(r.registrations, registration)
	return nil
}

func (r *registry) unregister(url string) error {
	for i := range r.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service info %s not found", url)
}

// reg所依赖的服务
func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequiredService {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, PatchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	//
	if err := r.sendPatch(p, reg.ServiceUpdateURL); err != nil {
		return err
	}
	return nil
}

func (r *registry) sendPatch(p patch, url string) error {
	txt, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(txt))
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Send post to service response with code %v", resp.StatusCode)
	}
	return nil
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct {
}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Printf("Decode request body error, %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Add service: %v with URL: %v \n", r.ServiceName, r.ServiceURL)
		err = reg.register(r)
		if err != nil {
			log.Println("Register error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Read body error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := string(payload)
		log.Printf("Unregister service with URL: %v", url)
		err = reg.unregister(url)
		if err != nil {
			log.Println("Unregister service error.", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
