package nginx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/golang/glog"
)

// Controller allows starting/stoping NGINX and applying a new configuration.
type Controller struct {
	doneCh  chan error
	started bool
}

// NewController creates a new controller.
func NewController() *Controller {
	return &Controller{
		doneCh: make(chan error),
	}
}

func (c *Controller) start() error {
	cmd := exec.Command("nginx", "-g", "daemon off;", "-c", "/etc/istio/proxy/nginx.conf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start NGINX: %v", err)
	}
	go func() {
		c.doneCh <- cmd.Wait()
		close(c.doneCh)
	}()
	return nil
}

// ExitStatus returns a chanel through which the exit status is returned when NGINX exits.
func (c *Controller) ExitStatus() <-chan error {
	return c.doneCh
}

// ApplyConfig updates NGINX configuration -- writes the config and reloads NGINX.
// On the first invocation ApplyConfig starts NGINX as well.
func (c *Controller) ApplyConfig(config Config) error {
	if err := writeConfig(config); err != nil {
		return fmt.Errorf("couldn't write NGINX configuration: %v", err)
	}

	if !c.started {
		if err := c.start(); err != nil {
			return fmt.Errorf("failed to start NGINX: %v", err)
		}
		c.started = true
	} else {
		if err := c.Reload(); err != nil {
			return fmt.Errorf("couldn't reload NGINX: %v", err)
		}
	}

	return nil
}

func writeConfig(config Config) error {
	for _, cfg := range config.HTTPConfigs {
		err := writeHTTPConfig(cfg)
		if err != nil {
			return fmt.Errorf("couldn't write the config file: %v", err)
		}
	}

	for _, cfg := range config.TCPConfigs {
		err := writeTCPConfig(cfg)
		if err != nil {
			return fmt.Errorf("couldn't write the config file: %v", err)
		}
	}

	err := writeMainConfig(config.Main)
	if err != nil {
		return fmt.Errorf("couldn't write the main config file: %v", err)
	}

	return nil
}

func writeHTTPConfig(cfg HTTPConfig) error {
	tmpl, err := template.New("config.tmpl").Parse(httpTemplate)
	if err != nil {
		return fmt.Errorf("couldn't parse the template: %v", err)
	}

	if glog.V(3) {
		err = tmpl.Execute(os.Stdout, cfg)
		if err != nil {
			return fmt.Errorf("couldn't execute the template: %v", err)
		}
	}

	name := "/etc/istio/proxy/conf.d/" + cfg.Name + ".conf"
	w, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("couldn't create the config file: %v", err)
	}

	defer w.Close()

	glog.V(2).Infof("Writing configuration to %v", name)

	err = tmpl.Execute(w, cfg)
	if err != nil {
		return fmt.Errorf("couldn't execute the template: %v", err)
	}

	return nil
}

func writeTCPConfig(cfg TCPConfig) error {
	tmpl, err := template.New("tcpconfig.tmpl").Parse(tcpTemplate)
	if err != nil {
		return fmt.Errorf("couldn't parse the template: %v", err)
	}

	if glog.V(3) {
		err = tmpl.Execute(os.Stdout, cfg)
		if err != nil {
			return fmt.Errorf("couldn't execute the template: %v", err)
		}
	}

	name := "/etc/istio/proxy/conf.d/" + cfg.Name + ".stream-conf"
	w, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("couldn't create the config file: %v", err)
	}

	defer w.Close()

	glog.V(2).Infof("Writing configuration to %v", name)

	err = tmpl.Execute(w, cfg)
	if err != nil {
		return fmt.Errorf("couldn't execute the template: %v", err)
	}

	return nil
}

// Reload reloads NGINX.
func (c *Controller) Reload() error {
	if err := shellOut("nginx -t -c /etc/istio/proxy/nginx.conf"); err != nil {
		return fmt.Errorf("invalid NGINX configuration detected, not reloading: %s", err)
	}
	if err := shellOut("nginx -s reload -c /etc/istio/proxy/nginx.conf"); err != nil {
		return fmt.Errorf("reloading NGINX failed: %s", err)
	}

	return nil
}

func shellOut(cmd string) (err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	glog.V(2).Infof("executing %s", cmd)

	command := exec.Command("sh", "-c", cmd)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err = command.Start()
	if err != nil {
		return fmt.Errorf("failed to execute %v, err: %v", cmd, err)
	}

	err = command.Wait()
	if err != nil {
		return fmt.Errorf("command %v stdout: %q\nstderr: %q\nfinished with error: %v", cmd,
			stdout.String(), stderr.String(), err)
	}
	return nil
}

// Quit shutdowns NGINX gracefully.
func (c *Controller) Quit() error {
	if err := shellOut("nginx -s quit -c /etc/istio/proxy/nginx.conf"); err != nil {
		return fmt.Errorf("failed to quit NGINX: %v", err)
	}
	return nil
}

func writeMainConfig(cfg Main) error {
	tmpl, err := template.New("main.tmpl").Parse(mainTemplate)
	if err != nil {
		return fmt.Errorf("couldn't parse the template: %v", err)
	}

	if glog.V(3) {
		err = tmpl.Execute(os.Stdout, cfg)
		if err != nil {
			return fmt.Errorf("couldn't execute the template: %v", err)
		}
	}

	name := "/etc/istio/proxy/nginx.conf"
	w, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("couldn't create the config file: %v", err)
	}

	defer w.Close()

	glog.Infof("Writing configuration to %v", name)

	err = tmpl.Execute(w, cfg)
	if err != nil {
		return fmt.Errorf("couldn't execute the template: %v", err)
	}

	return nil
}
