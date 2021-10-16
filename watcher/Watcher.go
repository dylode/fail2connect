package watcher

import (
	. "fail2connect/config"
	"github.com/hpcloud/tail"
	"io"
	"log"
	"regexp"
	"sync"
	"time"
)

type ConnectionInfo struct {
	connectedAt     time.Time
	connectionTries int
	banned          bool
}

type Watcher struct {
	config      *WatcherConfig
	connections map[string]ConnectionInfo
	lock        sync.Mutex
	known       []string
}

func New(config WatcherConfig) *Watcher {
	instance := &Watcher{
		config:      &config,
		connections: make(map[string]ConnectionInfo),
		lock:        sync.Mutex{},
		known:       []string{},
	}

	go instance.Start()

	return instance
}

func (w *Watcher) Start() {
	for true {
		t, err := tail.TailFile(w.config.PathToLogFile, tail.Config{Follow: true, Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
			Offset: 0,
		}})

		if err != nil {
			log.Fatalln(err)
		}

		for line := range t.Lines {
			w.Analyze(line.Text)
		}
	}
}

func (w *Watcher) Analyze(line string) {
	connectionRegex := regexp.MustCompile(w.config.ConnectionRegex)
	successRegex := regexp.MustCompile(w.config.SuccessRegex)

	if connectionRegex.MatchString(line) {
		ip := connectionRegex.FindStringSubmatch(line)[1]
		log.Println("New connection from ", ip)

		// If this Watcher trust known IPs, check if it is known. If it is, skip, since we trust this IP
		if w.config.TrustKnown && Find(w.known, ip) {
			return
		}

		w.lock.Lock()
		defer w.lock.Unlock()

		// Check if this IP address tries to connect multiple times within the ultimatum
		if connection, found := w.connections[ip]; found {
			connection.connectionTries++
			w.connections[ip] = connection
		} else {
			w.connections[ip] = ConnectionInfo{
				connectedAt:     time.Now(),
				connectionTries: 1,
				banned:          false,
			}
		}
	} else if successRegex.MatchString(line) {
		ip := successRegex.FindStringSubmatch(line)[1]
		log.Println("Successful connection from ", ip)
		w.lock.Lock()
		defer w.lock.Unlock()

		if _, found := w.connections[ip]; found {
			delete(w.connections, ip)

			// If this Watcher trusts known IPs, add it to the known list
			if w.config.TrustKnown && !Find(w.known, ip) {
				w.known = append(w.known, ip)
			}
		}
	}
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
