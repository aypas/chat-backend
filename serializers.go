package main

type peoplePayload struct {
	People      []string  `json:"people"`
	PayloadType   string  `json:"payloadType"`
}

type message struct {
	To          string 	 `json:"to"`
	Msg         string 	 `json:"msg"`
	PayloadType string   `json:"payloadType"`
}

// func newMessage(to, msg, payloadType string) []byte {
// 	var s message = message{To: to, Msg: msg, PayloadType: payloadType}
// 	b, e := json.Marshal(s)
// 	if e != nil {
// 		fmt.Println("something went wrong ", e)
// 		return nil
// 	}
// 	return b
// }

// func newPeoplePayload(people, payload string) []byte {
// 	var s peoplePayload = peoplePayload{People: people, PayloadType: payload}
// 	b, e := json.Marshal(s)
// 	if e != nil {
// 		fmt.Println("something went wrong ", e)
// 		return nil
// 	}
// 	return b 
// }

// func main() {
// 	p := peoplePayload{People: []string{"hehe", "popo"}}
// 	e, _ := json.Marshal(p)
// 	fmt.Println(string(e))
// 	var pp *peoplePayload = &peoplePayload{}
// 	json.Unmarshal(e, pp)
// 	fmt.Printf("%T",*pp)
// }

