/*
Description: This creates a server that can view and edited data in the form of key/value pair.
			 To add, delete, modify and view, make sure your path contains /test/add,  /test/delete,
			 /test/modify, and /test/view respectively. Here are an example for your to try:
			 -------------------------------------------------------------------------------------
			 add data: 				http://localhost:8081/test/add?key=test1&data=192.168.1.84
		  	 modify specified data: http://localhost:8081/test/modify?key=test1&data=192.168.1.202
			 delete specified data: http://localhost:8081/test/delete?key=test1
			 view specified data: 	http://localhost:8081/test/view?key=test1
			 view all data: 	    http://localhost:8081/test/view
			 -------------------------------------------------------------------------------------
Name: Johnson Zhuang
Date: 6/25/2019
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

//OpenFile opens the file named data.json
func OpenFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := make(map[string]string)

	JSONDecoder := json.NewDecoder(file)
	JSONDecoder.Decode(&d)

	return d, nil
}

//view will check if the key's validity and shows the data if valid
func viewPair(w http.ResponseWriter, r *http.Request, pair map[string]string, checks ...bool) {

	//outputs everything when no query
	if !checks[0] {
		for k, v := range pair {
			fmt.Fprintf(w, "key=%s data=%s\n", k, v)
		}
		return
	}

	//when there are no "key" in the query or that the key does not contain anything
	if !checks[1] {
		errHandler(w, r, 404)
		log.Println("Invalid query")
		return
	}

	q := r.URL.Query()
	key := q["key"][0]

	if !checks[2] {
		//when the key give is invalid
		errHandler(w, r, 404)
		fmt.Fprintf(w, "\nThe key \"%v\" cannot be viewed becasuse it does not exist", key)
		log.Println("Error info: cannot find matching key in the file")
		return
	}

	fmt.Fprint(w, pair[key])
}

//edit will check if the key's validity and edit the data if valid
func modifyPair(w http.ResponseWriter, r *http.Request, pair map[string]string, checks ...bool) {

	//when there are no "key" in the query or no query
	if !checks[0] || !checks[1] {
		errHandler(w, r, 404)
		log.Println("Query invalid, either no query or invalid length")
		return
	}

	q := r.URL.Query()
	key := q["key"][0]

	if !checks[2] {
		//when the key given is invalid
		errHandler(w, r, 400)
		fmt.Fprintf(w, "\nThe key \"%v\" cannot be modified because it does not exist", key)
		log.Println("Cannot find matching key")
		return
	}

	//when there are no "data" in the query or no query, or data does not contain anything
	if !checks[3] {
		errHandler(w, r, 404)
		fmt.Fprint(w, "\nThe given data is invalid")
		log.Println("Invalid data, does not exist or empty string")
		return
	}

	data := q["data"][0]

	fmt.Fprintf(w, "Modified key \"%v\": \nOriginal data: %v\nNew data: %v ", key, pair[key], data)
	pair[key] = data
	WriteFile(pair, "data.json")
}

//add will add the key/value pair
func addPair(w http.ResponseWriter, r *http.Request, pair map[string]string, checks ...bool) {

	//when there are no "key" in the query or no query
	if !checks[0] || !checks[1] {
		errHandler(w, r, 404)
		log.Println("Query invalid, either no query or invalid length")
		return
	}

	q := r.URL.Query()
	key := q["key"][0]

	if checks[2] {
		//when the key given is invalid
		errHandler(w, r, 400)
		fmt.Fprintf(w, "\nThe key \"%v\" cannot be added because it already exist", key)
		log.Println("Request to add existing key")
		return
	}

	//when there are no "data" in the query or no query
	if !checks[3] {
		errHandler(w, r, 404)
		fmt.Fprint(w, "\nThe give data is invalid")
		log.Println("Invalid data, does not exist or empty string")
		return
	}

	data := q["data"][0]
	pair[key] = data
	WriteFile(pair, "data.json")
	fmt.Fprintf(w, "New data added: %v", "\n"+key+"="+data)
}

//delete will delete the key/value pair
func deletePair(w http.ResponseWriter, r *http.Request, pairs map[string]string, checks ...bool) {

	//when there are no "key" in the query or no query
	if !checks[0] || !checks[1] {
		errHandler(w, r, 404)
		log.Println("Query invalid, either no query or invalid length")
		return
	}

	q := r.URL.Query()
	key := q["key"][0]
	if !checks[2] {
		//when the key given is invalid
		errHandler(w, r, 400)
		fmt.Fprintf(w, "\nThe key \"%v\" cannot be deleted because it does not exist", key)
		log.Println("Cannot find matching key")
		return
	}

	delete(pairs, key)
	WriteFile(pairs, "data.json")
	fmt.Fprintf(w, "The key \"%v\" and its data has been deleted", key)
}

//WriteFile takes in the new map and writes the content into the json file
func WriteFile(input map[string]string, filename string) {

	JSONWrite, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer JSONWrite.Close()

	jsEnconder := json.NewEncoder(JSONWrite)
	jsEnconder.Encode(input)
}

//errHandler handles errors
func errHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "Cannot find the page, error 404")
	} else if status == http.StatusBadRequest {
		fmt.Fprint(w, "Bad request, error 400")
	}
}

//testHandler checks to see if the path is right or not
func testHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("The path does not exists")
	errHandler(w, r, http.StatusNotFound)

}

//makeHandler opens the file and saves the info in a map called pair. The it checks the length
//of the query, if "key" exists, if "key" has any values, if that valye exists in pair, and then
//"data" is valid or not. It then sends all the information to the handler and they will apply correspoding actions
func makeHandler(f func(http.ResponseWriter, *http.Request, map[string]string, ...bool)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pair, err := OpenFile("data.json")
		if err != nil {
			fmt.Fprint(w, "Failed to open the file")
			log.Fatal("failed to opent the requested file")
		}

		q := r.URL.Query()
		length, keyword, pairValue, dataValid := false, false, false, false

		if len(q) > 0 {
			length = true
		}

		key, ok := q["key"]
		if ok {
			keyword = true
			if _, ok = pair[key[0]]; ok {
				pairValue = true
			}
		}

		data, ok := q["data"]
		if ok && data[0] != "" {
			dataValid = true
		}

		f(w, r, pair, length, keyword, pairValue, dataValid)
	}
}

//server starts the servers
func server() {
	log.Println("Server started...")

	http.HandleFunc("/", testHandler)
	http.HandleFunc("/test/view", makeHandler(viewPair))
	http.HandleFunc("/test/modify", makeHandler(modifyPair))
	http.HandleFunc("/test/add", makeHandler(addPair))
	http.HandleFunc("/test/delete", makeHandler(deletePair))

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	server()
}
