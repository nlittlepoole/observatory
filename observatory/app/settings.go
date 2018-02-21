package main

import(
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type Settings struct {
     Driver string
     Location string
     Threshold int64
     SampleRate float64
     Granularity string
}

var DEFAULT_SETTINGS Settings = Settings{
      Driver: "mon0",
      Location: "Wayne Manor",
      Threshold: 100,
      SampleRate: 1.0,
      Granularity: "hour",
}

const SETTINGS_FILE string = "settings.yaml"

var ACTIVE_SETTINGS Settings

func init() {
     if _, err := os.Stat(SETTINGS_FILE); os.IsNotExist(err) {
     	ACTIVE_SETTINGS = DEFAULT_SETTINGS
     	err := ACTIVE_SETTINGS.Save()
	ACTIVE_SETTINGS.Window()
	if err != nil {
	   panic(err)
	}
     } else {
       dat, _ := ioutil.ReadFile(SETTINGS_FILE)
       	    err := yaml.Unmarshal(dat, &ACTIVE_SETTINGS)
	    if err != nil {
	       panic(err)
	    }
     }
}


func (s *Settings) Window() time.Duration {
     switch s.Granularity {
     case "hour":
     	  return time.Hour
     case "minute":
     	  return time.Minute
     default:
     	  return time.Hour * 24
     }
}

func (s *Settings) Save() error {
     bytes, err := yaml.Marshal(s)
     if err == nil {
     	err = ioutil.WriteFile(SETTINGS_FILE, bytes, 0644)
     }
     return err
}