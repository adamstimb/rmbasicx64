package game

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/adamstimb/rmbasicx64/pkg/nimgobus"
	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Boot bool
}

type Game struct {
	Count int
	nimgobus.Nimbus
	Config            AppConfig
	PrettyPrintIndent string
}

// LoadConfig attempts to load settings from the config file.  If the file does not
// exist or is unreadable it will be ignored and default settings will be used.
func (g *Game) LoadConfig() {
	// Default settings
	g.Config = AppConfig{Boot: true}
	// Attempt to load settings from config file
	c, err := g.ReadConf()
	if err == nil {
		// Overwrite default settings
		g.Config = c
	}
}

func (g *Game) WriteConf(c AppConfig) bool {
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Printf("Error resolving directory of executable: %v", err)
		return false
	}
	configPath := filepath.Join(exeDir, "rmbasicx64config.yaml")
	data, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatal(err)
		return false
	}
	err = ioutil.WriteFile(configPath, data, 0666)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func (g *Game) ReadConf() (AppConfig, error) {
	// Try to get directory of executable
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Printf("Error resolving directory of executable: %v", err)
	}
	// If the config file doesn't exist, create one with default settings
	c := AppConfig{Boot: true}
	configPath := filepath.Join(exeDir, "rmbasicx64config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		g.WriteConf(c)
		return c, nil
	}
	// If the config file exists, try to load and parse it
	buf, err := ioutil.ReadFile(configPath)
	data := make(map[interface{}]interface{})
	if err != nil {
		log.Printf("Error loading config file: %v", err)
		return c, err
	}
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		log.Printf("Error parsing config file: %v", err)
		return c, err
	}
	// Successfully opened and parsed the yaml.  Now try to match it
	// up with config keys but fail silently if the keys don't match.
	for k, v := range data {
		switch k {
		case "boot":
			bootVal, ok := v.(bool)
			if ok {
				c.Boot = bootVal
			}
		}
	}
	return c, nil
}

func (g *Game) Update() error {
	g.Nimbus.Update()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw the Nimbus monitor on the screen and scale to current window size.
	monitorWidth, monitorHeight := g.Monitor.Size()

	// Get ebiten window size so we can scale the Nimbus screen up or down
	// but if (0, 0) is returned we're not running on a desktop so don't do any scaling
	windowWidth, windowHeight := ebiten.WindowSize()

	// Calculate aspect ratios of Nimbus monitor and ebiten screen
	monitorRatio := float64(monitorWidth) / float64(monitorHeight)
	windowRatio := float64(windowWidth) / float64(windowHeight)

	// If windowRatio > monitorRatio then clamp monitorHeight to windowHeight otherwise
	// clamp monitorWidth to screenWidth
	var scale, offsetX, offsetY float64
	switch {
	case windowRatio > monitorRatio:
		scale = float64(windowHeight) / float64(monitorHeight)
		offsetX = (float64(windowWidth) - float64(monitorWidth)*scale) / 2
		offsetY = 0
	case windowRatio <= monitorRatio:
		scale = float64(windowWidth) / float64(monitorWidth)
		offsetX = 0
		offsetY = (float64(windowHeight) - float64(monitorHeight)*scale) / 2
	}

	// Apply scale and centre monitor on screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(g.Monitor, op)
}
