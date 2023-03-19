package main

import (
	"flag"
	"fmt"
	"github.com/snowmerak/gotor/actor"
	"github.com/snowmerak/gotor/config"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	initFlag := flag.String("init", "", "-init <file path>")
	genFlag := flag.String("gen", "", "-gen <file path>")
	flag.Parse()

	if *initFlag != "" {
		fmt.Println("initialize config file: " + *initFlag)
		f, err := os.Create(*initFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
				return
			}
		}()
		encoder := yaml.NewEncoder(f)
		encoder.SetIndent(2)
		if err := encoder.Encode(&config.Config{
			Actors: []struct {
				Path        string `yaml:"path"`
				PackageName string `yaml:"package_name"`
				ActorName   string `yaml:"actor_name"`
				Channels    []struct {
					Name string `yaml:"name"`
					Type string `yaml:"type"`
				} `yaml:"channels"`
			}{
				{
					Path:        filepath.Join(".", "test"),
					PackageName: "ActorMap",
					ActorName:   "Map",
					Channels: []struct {
						Name string `yaml:"name"`
						Type string `yaml:"type"`
					}{
						{
							Name: "get",
							Type: "tuple.Tuple[string, chan string]",
						},
						{
							Name: "set",
							Type: "string",
						},
						{
							Name: "delete",
							Type: "string",
						},
					},
				},
			},
		}); err != nil {
			fmt.Println(err)
			return
		}
	}

	if *genFlag != "" {
		fmt.Println("generate actor file: " + *genFlag)

		cfgData, err := os.ReadFile(*genFlag)
		if err != nil {
			fmt.Println(err)
			return
		}

		cfg := new(config.Config)
		if err := yaml.Unmarshal(cfgData, cfg); err != nil {
			fmt.Println(err)
			return
		}

		wg := new(sync.WaitGroup)
		for _, a := range cfg.Actors {
			a := a
			wg.Add(1)
			go func() {
				defer wg.Done()

				fmt.Println("generate actor: " + a.ActorName)
				dir := filepath.Join(a.Path, a.PackageName)
				if err := os.MkdirAll(dir, 0755); err != nil {
					fmt.Println(err)
					return
				}

				f, err := os.Create(filepath.Join(dir, a.ActorName+".go"))
				if err != nil {
					fmt.Println(err)
					return
				}
				defer func() {
					if err := f.Close(); err != nil {
						fmt.Println(err)
						return
					}
				}()

				channels := make(map[string]string, len(a.Channels))
				for _, c := range a.Channels {
					channels[c.Name] = c.Type
				}
				data := actor.Generate(a.PackageName, a.ActorName, channels)
				if _, err := f.Write(data); err != nil {
					fmt.Println(err)
					return
				}
			}()
		}
		wg.Wait()
	}
}
