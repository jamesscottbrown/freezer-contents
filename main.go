package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	port := flag.String("port", ":8080", "port to serve on")
	flag.Parse()
	fmt.Println("Serving on port", *port)

	mux := http.NewServeMux()

	mux.HandleFunc("/state", handleStateRequest)
	mux.HandleFunc("/remove", handleRemoveRequest)
	mux.HandleFunc("/move", handleMoveRequest)
	mux.HandleFunc("/add", handleAddRequest)

	mux.HandleFunc("/list", handleListRequest)
	mux.HandleFunc("/", handleRootRequest)
	http.ListenAndServe(*port, mux)

}

func handleRootRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	handleListRequest(w, r)
}

func handleListRequest(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(`WOO!`))
}

func readState(f string) ([]byte, error) {
	out, err := ioutil.ReadFile("contents.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return out, nil

}

func handleStateRequest(w http.ResponseWriter, r *http.Request) {
	// read the contents.json file
	f := "contents.json"
	out, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(out))
	return
}

type AddBody struct {
	Name       string
	Date       string
	Freezer    string
	Containers []string
}

func handleAddRequest(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var t AddBody
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error parsing request body:", err)
		return
	}

	contents, err := readContents("contents.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// check if freezer exists
	var freezerExists bool
	for _, freezer := range contents.Freezers {
		if freezer.Name == t.Freezer {
			freezerExists = true
		}
	}
	if !freezerExists {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error: freezer does not exist")
		return
	}

	// add item to freezer
	for i, freezer := range contents.Freezers {
		if freezer.Name == t.Freezer {
			contents.Freezers[i].Contents = append(contents.Freezers[i].Contents, Item{t.Name, t.Date, t.Containers})
		}
	}

	// save JSON to file
	err = writeContents("contents.json", contents)

	// return the json
	json := json.NewEncoder(w)
	err = json.Encode(contents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to convert updated contents to JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	return

}

// remove item from contents.json
type RemoveBody struct {
	Container string
}

func handleRemoveRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t RemoveBody
	err := decoder.Decode(&t)
	if err != nil {
		// panic(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error parsing request body:", err)
		return
	}

	contents, err := readContents("contents.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// find the container to remove
	for i, freezer := range contents.Freezers {
		for j, item := range freezer.Contents {

			for k, container := range item.Containers {
				// remove container
				if container == t.Container {
					contents.Freezers[i].Contents[j].Containers = append(contents.Freezers[i].Contents[j].Containers[:k], contents.Freezers[i].Contents[j].Containers[k+1:]...)
				}
			}

			// if no containers, remove item too
			if len(contents.Freezers[i].Contents[j].Containers) == 0 {
				contents.Freezers[i].Contents = append(contents.Freezers[i].Contents[:j], contents.Freezers[i].Contents[j+1:]...)
			}
		}
	}

	// save JSON to file
	err = writeContents("contents.json", contents)

	// return the json
	json := json.NewEncoder(w)
	err = json.Encode(contents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to convert updated contents to JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	return
}

// remove item from contents.json
type MoveBody struct {
	Container  string
	NewFreezer string
}

func handleMoveRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t MoveBody
	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error parsing request body:", err)
		return
	}

	contents, err := readContents("contents.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// find the container to move
	var moveItem Item
	for i, freezer := range contents.Freezers {
		for j, item := range freezer.Contents {

			for k, container := range item.Containers {
				// remove container from old freezer
				if container == t.Container {
					moveItem = contents.Freezers[i].Contents[j]
					contents.Freezers[i].Contents[j].Containers = append(contents.Freezers[i].Contents[j].Containers[:k], contents.Freezers[i].Contents[j].Containers[k+1:]...)
				}
			}

			// if no containers, remove parent item too
			if len(contents.Freezers[i].Contents[j].Containers) == 0 {
				contents.Freezers[i].Contents = append(contents.Freezers[i].Contents[:j], contents.Freezers[i].Contents[j+1:]...)
			}
		}
	}

	// add item to new freezer
	for i, freezer := range contents.Freezers {
		if freezer.Name == t.NewFreezer {

			// Look for item in new freezer
			itemExists := false
			for _, item := range freezer.Contents {
				if item.Name == moveItem.Name && item.Date == moveItem.Date {
					itemExists = true
					item.Containers = append(item.Containers, t.Container)
				}
			}

			// If item doesn't exist in new freezer, add it
			if !itemExists {
				contents.Freezers[i].Contents = append(contents.Freezers[i].Contents, Item{moveItem.Name, moveItem.Date, []string{t.Container}})
			}
		}
	}

	// save JSON to file
	err = writeContents("contents.json", contents)

	// return the json
	json := json.NewEncoder(w)
	err = json.Encode(contents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Failed to convert updated contents to JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	return
}
