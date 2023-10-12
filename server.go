package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
        "encoding/json"

	"github.com/deiu/rdf2go"
)

const defaultPort = 8069

type assocMap map[string]interface{}

func readStringBody(inputStream io.Reader) string {
	buffer := new(strings.Builder)
	if _, err := io.Copy(buffer, inputStream); err != nil {
		panic(err)
	}

	return buffer.String()
}

type turtleObject struct {
	IRI string `json:"iri"`
}

type turtlePredicate struct {
	From	*turtleObject	`json:"from"`
	Into	*turtleObject	`json:"into"`
	Name	string		`json:"name"`
}

type turtleDescr struct {
	Objects		[]*turtleObject		`json:"objects"`
	Predicates	[]*turtlePredicate	`json:"predicates"`
}

func parseTurtle(input string) (*turtleDescr, error) {
	graph := rdf2go.NewGraph("")
	reader := strings.NewReader(input)

	if err := graph.Parse(reader, "text/turtle"); err != nil {
		return nil, err
	}

	graphLen := graph.Len()
	connections := make([]*turtlePredicate, 0, graphLen)
	objects := make(map[string]*turtleObject)
	

	for tripple := range graph.IterTriples() {
		objectStrVal := tripple.Object.RawValue()

		// try object
		if _, ok := objects[objectStrVal]; !ok {
			objects[objectStrVal] = &turtleObject{
				IRI: objectStrVal,
			}
		}

		subjectStrVal := tripple.Subject.RawValue()

		// try subject
		if _, ok := objects[subjectStrVal]; !ok {
			objects[subjectStrVal] = &turtleObject {
				IRI: subjectStrVal,
			}
		}

		connections = append(connections, &turtlePredicate{
			From: objects[objectStrVal],
			Into: objects[subjectStrVal],
			Name: tripple.Predicate.RawValue(),
		})
	}

	turtleObjs := make([]*turtleObject, 0, len(objects))
	for _, val := range objects {
		turtleObjs = append(turtleObjs, val)
	}

	return &turtleDescr{
		Objects: turtleObjs,
		Predicates: connections,
	}, nil
}

func handleIntoTripples(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")

	turtleStrInput := readStringBody(r.Body)
	
        enc := json.NewEncoder(w)
	turtle, err := parseTurtle(turtleStrInput)

	if err != nil {
            fmt.Printf("Error: %s", err.Error())
            w.WriteHeader(http.StatusInternalServerError)
            enc.Encode(map[string]any{
                "error": err.Error(),
            })
            return
	}

        w.WriteHeader(http.StatusOK)
        enc.Encode(turtle)

	return
}

func main() {
        http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Hello!")
        })

        http.HandleFunc("/intoTripples", handleIntoTripples)

	var port uint = defaultPort
	if env := os.Getenv("SERVER_PORT"); env != "" {
		if parsedPort, err := strconv.Atoi(env); err == nil {
			port = uint(parsedPort)
		}
		
	}

	fmt.Printf("Listening on port %d\n", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

