package regsitry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	serviceUpdateUrl, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateUrl.Path, &serviceUpdateHandler{})

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err = enc.Encode(&r)
	if err != nil {
		return err
	}

	// 向注册中心注册自己
	resp, err := http.Post(ServerUrl, "application/json", &buf)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service.Regitry service responsed with code %v", resp.StatusCode)
	}

	return nil
}

type serviceUpdateHandler struct {
}

// ServeHTTP 用来更新依赖的服务
func (h serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	provide.Update(p)
}

func UnregisterService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServerUrl, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unregister service. Regitry service responsed with code %v", resp.StatusCode)
	}
	return nil
}

type providers struct {
	services map[serviceName][]string
	mutex    *sync.Mutex
}

// Update 更新服务依赖 patch中包含了需要新增的依赖和需要剔除的依赖
func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	for _, patchEntry := range pat.Removed {
		if providerUrls, ok := p.services[patchEntry.Name]; ok {
			for i := range providerUrls {
				if providerUrls[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(providerUrls[:i], providerUrls[i+1:]...)
				}
			}
		}
	}
}

func (p providers) get(name serviceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no provider available for service %v", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

// GetProvider 获取服务URL
func GetProvider(name serviceName) (string, error) {
	return provide.get(name)
}

var provide = providers{
	services: make(map[serviceName][]string),
	mutex:    new(sync.Mutex),
}
