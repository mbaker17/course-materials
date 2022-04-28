package hscan

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

//==========================================================================\\

var shalookup map[string]string
var md5lookup map[string]string

func GuessSingle(sourceHash string, filename string) string {

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	shalookup = make(map[string]string)
	md5lookup = make(map[string]string)

	for scanner.Scan() {
		password := scanner.Text()

		// TODO - From the length of the hash you should know which one of these to check ...
		// add a check and logicial structure
		if len(sourceHash) == 32 {
			hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
			if hash == sourceHash {
				fmt.Printf("[+] Password found (MD5): %s\n", password)
				return password
			}
		} else if len(sourceHash) == 64 {
			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
			if hash == sourceHash {
				fmt.Printf("[+] Password found (SHA-256): %s\n", password)
				return password
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
	return ""
}

var mutex = &sync.Mutex{}

func Addmd5 (password string) {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	mutex.Lock()
	md5lookup[hash] = password
	mutex.Unlock()
}

var mutextwo = &sync.Mutex{}

func Addsha (password string) {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	mutextwo.Lock()
	shalookup[hash] = password
	mutextwo.Unlock()
}

func addlist (password string) {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	mutex.Lock()
	shalookup[hash] = password
	mutex.Unlock()
	hash = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	mutextwo.Lock()
	md5lookup[hash] = password
	mutextwo.Unlock()
}

func GenHashMaps(filename string) {

	//TODO
	//itterate through a file (look in the guessSingle function above)
	//rather than check for equality add each hash:passwd entry to a map SHA and MD5 where the key = hash and the value = password
	//TODO at the very least use go subroutines to generate the sha and md5 hashes at the same time
	//OPTIONAL -- Can you use workers to make this even faster

	//TODO create a test in hscan_test.go so that you can time the performance of your implementation
	//Test and record the time it takes to scan to generate these Maps
	// 1. With and without using go subroutines
	// 2. Compute the time per password (hint the number of passwords for each file is listed on the site...)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		password := scanner.Text()

			// go Addmd5(password)
			// go Addsha(password)

		//    go addlist(password)

			// Add wait group?
		      go addlist(password)


			// hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
			// md5lookup[hash] = password

			// hash = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
			// shalookup[hash] = password	
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

}

func GetSHA(hash string) (string, error) {
	password, ok := shalookup[hash]
	if ok {
		fmt.Printf("[+] Password found: %s\n", password)
		return password, nil

	} else {

		return "", errors.New("password does not exist")

	}
}

//TODO
func GetMD5(hash string) (string, error) {
	password, ok := md5lookup[hash]
	if ok {
		fmt.Printf("[+] Password found: %s\n", password)
		return password, nil

	} else {

		return "", errors.New("password does not exist")

	}
}
