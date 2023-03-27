package main

import (
	"flag"
	"log"
	"os"
	"time"

	"wppserver/internal/app"
	"wppserver/pkg/config"
	"wppserver/pkg/utils"

	"github.com/kardianos/service"
)

var envfile *string

type program struct {
	startTime time.Time
}

func (p *program) Start(s service.Service) error {
	p.startTime = time.Now()

	log.Printf("Running in service manager: %v", !service.Interactive())
	log.Printf("System servicer: %v", service.Platform())

	go p.run()
	return nil
}

func (p *program) run(args ...string) error {
	app := &app.App{}
	if *envfile == "" {
		// Get .env file in root dir
		path, err := utils.GetWorkDir()
		if err != nil {
			log.Fatal(err)
		}
		*envfile = path + "/.env"
	}

	if err := app.Initialize(config.GetConfig(*envfile)); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (p *program) Stop(s service.Service) error {
	log.Print("Stopping execution")
	log.Printf("Server uptime: %s", time.Since(p.startTime))
	return nil
}

func main() {
	path, err := utils.GetWorkDir()
	if err != nil {
		log.Fatal(err)
	}
	if err = os.Chdir(path); err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	logFile, err := os.OpenFile(path+"/debug.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	mode := flag.String("mode", "run", "Run server or configure background service manager. Valid actions: ['run', 'start', 'stop', 'restart', 'install', 'uninstall']")
	debug := flag.String("debug", "disabled", "Enable file debug output mode. Valid actions: ['enabled', 'disabled']")
	envfile = flag.String("envfile", "", "Load an .env file for this process. Accepts an absolute path with the file name.")
	flag.Parse()

	svConfig := &service.Config{
		Name:        "wppserver",
		DisplayName: "Wppserver Service v1.0.0",
		Description: "Wppserver API Server v1.0.0",
		Arguments: []string{
			"-debug=enabled", "-envfile=" + *envfile},
	}

	prg := &program{}
	s, err := service.New(prg, svConfig)
	if err != nil {
		log.Fatal(err)
	}
	if *debug == "enabled" {
		log.SetOutput(logFile)
	}

	switch *mode {
	case "run":
		if err = s.Run(); err != nil {
			log.Fatal(err)
		}

	case "start":
		if err = s.Start(); err != nil {
			log.Fatal(err)
		}

	case "stop":
		if err = s.Stop(); err != nil {
			log.Fatal(err)
		}

	case "restart":
		if err = s.Restart(); err != nil {
			log.Fatal(err)
		}

	case "install":
		if err = s.Install(); err != nil {
			log.Fatal(err)
		}

	case "uninstall":
		s.Stop()

		if err = s.Uninstall(); err != nil {
			log.Fatal(err)
		}

	default:
		log.Printf("The %v value is not valid: ", *mode)
	}
}
