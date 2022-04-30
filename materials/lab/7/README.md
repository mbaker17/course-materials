# Lab 7
Due April 9th at 11:59PM

## Development Work [18 points]
- [3pt] Complete the TODOs in [main.go](course-materials/materials/lab/7/main/main.go)
- [12pt] Complete the TODOs in [hscan.go](course-materials/materials/lab/7/hscan/hscan.go)
- [3pt] Create at least one new test in [hscan_test.go](course-materials/materials/lab/7/hscan/hscan_test.go)

## Capture  details [2pts]
- Capture Timing Details (per hscan.go) for various implementation of creating the hash maps

## Submit 
Two Passwords found.
1. letmein
2. p@ssword

1. Link to your Github Repository [16pts]
2. Report the numbers [4pts]
Time without go routines
797.270039ms

Time with go routines
783.070521ms

Since I had to use mutex locks with every go routine it didnt seem to speed up my searches any significant amount.

2. List of Collaborator
No one
