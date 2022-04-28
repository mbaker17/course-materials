// Optional Todo

package hscan

import (
	"testing"
)

func TestGuessSingle(t *testing.T) {
	got := GuessSingle("77f62e3524cd583d698d51fa24fdff4f", "/home/matt/Desktop/cstwo/course-materials/materials/lab/7/main/wordlist.txt") // Currently function returns only number of open ports
	want := "Nickelback4life"
	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}

}

func TestGenHashMaps(t *testing.T) {
	GenHashMaps("/home/matt/Desktop/cstwo/course-materials/materials/lab/7/main/wordlist.txt")
}