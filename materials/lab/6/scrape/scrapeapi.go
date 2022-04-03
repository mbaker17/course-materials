package scrape

// scrapeapi.go HAS TEN TODOS - TODO_5-TODO_14 and an OPTIONAL "ADVANCED" ASK

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "os"
    "path/filepath"
    "strconv"
    "github.com/gorilla/mux"
    "regexp"
)

const LOG_LEVEL = 2


//==========================================================================\\

// Helper function walk function, modfied from Chap 7 BHG to enable passing in of
// additional parameter http responsewriter; also appends items to global Files and 
// if responsewriter is passed, outputs to http 

func walkFn(w http.ResponseWriter) filepath.WalkFunc {
    count := 0
    return func(path string, f os.FileInfo, err error) error {
        w.Header().Set("Content-Type", "application/json")

        for _, r := range regexes {
            if r.MatchString(path) {
                var tfile FileInfo
                dir, filename := filepath.Split(path)
                tfile.Filename = string(filename)
                tfile.Location = string(dir)

                //TODO_5: As it currently stands the same file can be added to the array more than once 
                //TODO_5: Prevent this from happening by checking if the file AND location already exist as a single record
                exists := false

                for i := 0; i < len(Files); i++ {
                    if tfile.Filename == Files[i].Filename && tfile.Location == Files[i].Location {
                        exists = true
                    }
                }

                if(exists == false) {
                    Files = append(Files, tfile)
                    count++
                }

                if w != nil && len(Files)>0 {

                    //TODO_6: The current key value is the LEN of Files (this terrible); 
                    //TODO_6: Create some variable to track how many files have been added
                    w.Write([]byte(`"`+(strconv.FormatInt(int64(count), 10))+`":  `))
                    json.NewEncoder(w).Encode(tfile)
                    w.Write([]byte(`,`))

                } 
                
                if LOG_LEVEL == 2{
                    log.Printf("[+] HIT: %s\n", path)
                }
            }

        }
        return nil
    }

}

//TODO_7: One of the options for the API is a query command
//TODO_7: Create a walkFn2 function based on the walkFn function, 
//TODO_7: Instead of using the regexes array, define a single regex 
//TODO_7: Hint look at the logic in scrape.go to see how to do that; 
//TODO_7: You won't have to itterate through the regexes for loop in this func!

func walkFn2(w http.ResponseWriter, query string) filepath.WalkFunc {
    count := 0
    return func(path string, f os.FileInfo, err error) error {
        w.Header().Set("Content-Type", "application/json")

        r := regexp.MustCompile(`(?i)`+query)


        if r.MatchString(path) {
            var tfile FileInfo
            dir, filename := filepath.Split(path)
            tfile.Filename = string(filename)
            tfile.Location = string(dir)

            //TODO_5: As it currently stands the same file can be added to the array more than once 
            //TODO_5: Prevent this from happening by checking if the file AND location already exist as a single record
            exists := false

            for i := 0; i < len(Files); i++ {
                if tfile.Filename == Files[i].Filename && tfile.Location == Files[i].Location {
                    exists = true
                }
            }

            if(exists == false) {
                Files = append(Files, tfile)
                count++
            }


            

            if w != nil && len(Files)>0 {

                //TODO_6: The current key value is the LEN of Files (this terrible); 
                //TODO_6: Create some variable to track how many files have been added
                w.Write([]byte(`"`+(strconv.FormatInt(int64(count), 10))+`":  `))
                json.NewEncoder(w).Encode(tfile)
                w.Write([]byte(`,`))

            } 
            
            if LOG_LEVEL == 2{
                log.Printf("[+] HIT: %s\n", path)
            }

        }
        return nil
    }
}

//==========================================================================\\

func APISTATUS(w http.ResponseWriter, r *http.Request) {
    if LOG_LEVEL == 2{
		log.Printf("Entering %s end point", r.URL.Path)
	}
	
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{ "status" : "API is up and running ",`))
    var regexstrings []string
    
    for _, regex := range regexes{
        regexstrings = append(regexstrings, regex.String())
    }

    w.Write([]byte(` "regexs" :`))
    json.NewEncoder(w).Encode(regexstrings)
    w.Write([]byte(`}`))
	log.Println(regexes)

}


func MainPage(w http.ResponseWriter, r *http.Request) {
	if LOG_LEVEL == 2{
		log.Printf("Entering %s end point", r.URL.Path)
	}
    
    w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(http.StatusOK)
    //TODO_8 - Write out something better than this that describes what this api does

	fmt.Fprintf(w, "<html><body><H1>Welcome to my awesome File page</H1><h3>Here is a list of endpoints and a brief description</h3><ul><li>/api-status: Shows status of the API and Regexs</li><li>/indexer: Search specified location for regexes</li><li>/search: Search for if a something exists</li><li>/addsearch/{regex}: Add a regular expression to the array of search terms</li><li>/clear: Clear the list of search terms</li><li>/reset: Reset search list to original </li></ul></body>")
}

var FOUND bool = false

func FindFile(w http.ResponseWriter, r *http.Request) {
    if LOG_LEVEL == 2{
    log.Printf("Entering %s end point", r.URL.Path)
    }
    q, ok := r.URL.Query()["q"]

    w.WriteHeader(http.StatusOK)
    if ok && len(q[0]) > 0 {
        if LOG_LEVEL == 2{
        log.Printf("Entering search with query=%s",q[0])
        }

        // ADVANCED: Create a function in scrape.go that returns a list of file locations; call and use the result here
        // e.g., func finder(query string) []string { ... }

        for _, File := range Files {
            if File.Filename == q[0] {
                json.NewEncoder(w).Encode(File.Location)
                FOUND = true
            }
        }
        if FOUND == false{
            json.NewEncoder(w).Encode("No matches exist")
        }
        //TODO_9: Handle when no matches exist; print a useful json response to the user; hint you might need a "FOUND variable" to check here ...

    } else {
        // didn't pass in a search term, show all that you've found
        w.Write([]byte(`"files":`))    
        json.NewEncoder(w).Encode(Files)
    }
}

func IndexFiles(w http.ResponseWriter, r *http.Request) {
    if LOG_LEVEL == 2{
		log.Printf("Entering %s end point", r.URL.Path)
	}
    w.Header().Set("Content-Type", "application/json")

    location, locOK := r.URL.Query()["location"]
    regex, regexOK := r.URL.Query()["regex"]

    
    //TODO_10: Currently there is a huge risk with this code ... namely, we can search from the root /
    //TODO_10: Assume the location passed starts at /home/ (or in Windows pick some "safe?" location)
    //TODO_10: something like ...  rootDir string := "???"
    //TODO_10: create another variable and append location[0] to rootDir (where appropriate) to patch this hole
    rootDir := "/home/matt"

    if locOK && len(location[0]) > 0 {
        w.WriteHeader(http.StatusOK)

    } else {
        w.WriteHeader(http.StatusFailedDependency)
        w.Write([]byte(`{ "parameters" : {"required": "location",`))    
        w.Write([]byte(`"optional": "regex"},`))    
        w.Write([]byte(`"examples" : { "required": "/indexer?location=/xyz",`))
        w.Write([]byte(`"optional": "/indexer?location=/xyz&regex=(i?).md"}}`))
        return 
    }

    //wrapper to make "nice json"
    w.Write([]byte(`{ `))
    
    // TODO_11: Currently the code DOES NOT do anything with an optionally passed regex parameter
    // Define the logic required here to call the new function walkFn2(w,regex[0])
    // Hint, you need to grab the regex parameter (see how it's done for location above...) 
    
    // if regexOK
    //   call filepath.Walk(location[0], walkFn2(w, `(i?)`+regex[0]))
    // else run code to locate files matching stored regular expression
    if regexOK {
        filepath.Walk(rootDir+location[0], walkFn2(w, `(i?)`+regex[0]))
    } else {
        filepath.Walk(rootDir+location[0], walkFn(w))
    }

    // if err := filepath.Walk(rootDir+location[0], walkFn(w)); err != nil {
	// 	log.Panicln(err)
	// }

    //wrapper to make "nice json"
    w.Write([]byte(` "status": "completed"} `))

}


//TODO_12 create endpoint that calls resetRegEx AND *** clears the current Files found; ***
//TODO_12 Make sure to connect the name of your function back to the reset endpoint main.go!
func RESETARRAY(w http.ResponseWriter, r *http.Request) {
    if LOG_LEVEL == 2{
        log.Printf("Entering %s end point", r.URL.Path)
    }
    w.Header().Set("Content-Type", "application/json")
    resetRegEx()
    Files = nil
}


//TODO_13 create endpoint that calls clearRegEx ; 
//TODO_12 Make sure to connect the name of your function back to the clear endpoint main.go!
func CLEARREGEX(w http.ResponseWriter, r *http.Request) {
    if LOG_LEVEL == 2{
        log.Printf("Entering %s end point", r.URL.Path)
    }
    w.Header().Set("Content-Type", "application/json")
    clearRegEx()
}

//TODO_14 create endpoint that calls addRegEx ; 
//TODO_12 Make sure to connect the name of your function back to the addsearch endpoint in main.go!
// consider using the mux feature
// params := mux.Vars(r)
// params["regex"] should contain your string that you pass to addRegEx
// If you try to pass in (?i) on the command line you'll likely encounter issues
// Suggestion : prepend (?i) to the search query in this endpoint
func AddRegEx(w http.ResponseWriter, r *http.Request) {
    if LOG_LEVEL == 2{
        log.Printf("Entering %s end point", r.URL.Path)
    }
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    addRegEx(`(?i)`+params["regex"])
}