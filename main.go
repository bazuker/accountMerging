package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Account struct {
	Application string   `json:"application"`
	Emails      []string `json:"emails"`
	Name        string   `json:"name"`
}

type Person struct {
	Applications []string `json:"applications"`
	Emails       []string `json:"emails"`
	Name         string   `json:"name"`

	index          int
	representative int
}

func mergeArrays[T comparable](slices ...[]T) []T {
	m := map[T]bool{}
	for _, slice := range slices {
		for _, key := range slice {
			m[key] = true
		}
	}
	result := make([]T, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	return result
}

func merge(accounts []*Account) []*Person {
	// Convert Account structs into Person structs.
	// Potentially this could be done during unmarshalling.
	people := make([]*Person, len(accounts))
	for i, a := range accounts {
		people[i] = &Person{
			index:          i,
			representative: i,
			Name:           a.Name,
			Emails:         a.Emails,
			Applications:   []string{a.Application},
		}
	}

	// Merge accounts with DSU.
	emailsMap := make(map[string]*Person)
	for i := 0; i < len(people); i++ {
		for _, e := range people[i].Emails {
			if p, ok := emailsMap[e]; ok {
				// Pick a bigger group and merge into it.
				// The bigger (major) group will be a representative of a smaller (minor) group.
				major := p
				minor := people[i]
				if len(people[i].Emails) > len(p.Emails) {
					major = people[i]
					minor = p
				}
				// Find the group representative.
				for major.representative != major.index {
					major = people[major.representative]
				}
				minor.representative = major.index
				// Merge accounts.
				major.Emails = mergeArrays[string](major.Emails, minor.Emails)
				major.Applications = mergeArrays[string](major.Applications, minor.Applications)
			} else {
				emailsMap[e] = people[i]
			}
		}
	}

	// Only keep the groups that represent themselves.
	result := make([]*Person, 0)
	for i := 0; i < len(people); i++ {
		if people[i].representative == people[i].index {
			result = append(result, people[i])
		}
	}
	return result
}

func ProcessAccounts(accountsFilename string) (string, error) {
	data, err := os.ReadFile(accountsFilename)
	if err != nil {
		return "", fmt.Errorf("failed to read the accounts: %s", err)
	}

	var accounts []*Account
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal accounts: %s", err)
	}

	result := merge(accounts)

	data, err = json.MarshalIndent(result, "", "\t")
	if err != nil {
		return "", fmt.Errorf("failed to marshal accounts: %s", err)
	}

	return string(data), nil
}

func main() {
	output, err := ProcessAccounts("accounts.json")
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(output)
}
