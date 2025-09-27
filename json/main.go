package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "os"
    "time"
)

type Person struct {
    Name      string    `json:"name"`
    Height    string    `json:"height"`
    Mass      string    `json:"mass"`
    HairColor string    `json:"hair_color"`
    SkinColor string    `json:"skin_color"`
    EyeColor  string    `json:"eye_color"`
    BirthYear string    `json:"birth_year"`
    Gender    string    `json:"gender"`
    Homeworld string    `json:"homeworld"`
    Films     []string  `json:"films"`
    Species   []string  `json:"species"`
    Vehicles  []string  `json:"vehicles"`
    Starships []string  `json:"starships"`
    Created   time.Time `json:"created"`
    Edited    time.Time `json:"edited"`
    URL       string    `json:"url"`
}

func main() {
    //var p Person
    //p = unmarshal()
    //fmt.Printf("%# v\n", pretty.Formatter(p))

    //p = unmarshalFromAPI()
    //fmt.Printf("%# v\n", pretty.Formatter(p))

    //o := unstructured()
    //fmt.Printf("%# v\n", pretty.Formatter(o))

    //ch := make(chan Person)
    //go decode(ch)
    //for {
    //	person, ok := <-ch
    //	if !ok {
    //		break
    //	}
    //	fmt.Printf("%# v\n", pretty.Formatter(person))
    //}
}

func unmarshal() (person Person) {
    file, err := os.Open("skywalker.json")
    if err != nil {
        log.Println("Error opening json file:", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        log.Println("Error reading json data:", err)
    }

    err = json.Unmarshal(data, &person)
    if err != nil {
        log.Println("Error unmarshalling json data:", err)
    }
    return
}

func unmarshalFromAPI() (person Person) {
    resp, err := http.Get("https://swapi.dev/api/people/1")
    if err != nil {
        log.Println("Cannot get from URL", err)
    }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Error reading json data:", err)
    }

    err = json.Unmarshal(data, &person)
    if err != nil {
        log.Println("Error unmarshalling json data:", err)
    }
    return
}

func unstructured() (output map[string]any) {
    file, err := os.Open("unstructured.json")
    if err != nil {
        log.Println("Error opening json file:", err)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        log.Println("Error reading json data:", err)
    }

    err = json.Unmarshal(data, &output)
    if err != nil {
        log.Println("Error unmarshalling json data:", err)
    }
    return
}

func decode(p chan Person) {
    file, err := os.Open("people_stream.json")
    if err != nil {
        log.Println("Error opening json file:", err)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    for {
        var person Person
        err := decoder.Decode(&person)
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Println("Error decoding json data:", err)
            break
        }
        p <- person
    }
    close(p)
}
