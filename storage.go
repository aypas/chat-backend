package main 
import (
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
)

const dur = time.Hour*24

type store struct {
	Timestamp 	time.Time
	Logs 		map[string][]message
	Name 		string
}

//run when ok is false in hub's broadcast select...means user is not there to receive messages
func newUserStore(name string) store {
	return store{
		Timestamp: time.Now(),
		Logs: make(map[string][]message), //testing with int
		Name: name,
	}
}

func fileExists(name string) bool {
	_, err := os.Stat("store/"+name+".json")
	if err != nil {
		return false
	}
	return true
}

func (s *store) saveFile() error {
	//if file exists
	b, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		fmt.Println("big problem...")
		return err
	}
	f := ioutil.WriteFile("store/"+s.Name+".json", b, 0644)
	if f != nil {
		fmt.Println("couldn't save to file")
		return f
	}
	return nil
} 

func readFile(name string) (store, error) {
	fmt.Println(name)
	b, err := ioutil.ReadFile("store/"+name+".json")
	if err != nil {
		fmt.Println("couldn't read file", err)
		return  store{}, err
	}
	var s *store = &store{}
	json.Unmarshal(b, s)
	return *s, nil
}

//call as goroutine
func storeMsg(msg *message) {
	//msg to from msg payload
	var str store
	if fileExists(msg.To) {
		str, _ = readFile(msg.To)
	} else {
		str = newUserStore(msg.To)
	}
	str.Logs[msg.From] = append(str.Logs[msg.From], *msg)
	str.saveFile() 
}

//on register, look for logs, send if found
func sendLogs(name string, c *client) {
	if !fileExists(name) {
		return
	}
	s, e := readFile(name)
	if e != nil {
		return
	}
	if len(s.Logs) == 0 {
		return
	}
	c.writeTo <- &message{To: name, From: name, PayloadType: "storedData", Msg: s.Logs}
}