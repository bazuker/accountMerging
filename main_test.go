package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

// TestMerge demonstrates that the merge algorithm works with different sets of input data.
func TestMerge(t *testing.T) {
	// 2 in 1
	accounts := []*Account{
		{
			Application: "x",
			Emails:      []string{"a", "b", "c"},
			Name:        "Person 1",
		},
		{
			Application: "y",
			Emails:      []string{"c", "d"},
			Name:        "Person 1",
		},
	}

	people := merge(accounts)

	assert.Len(t, people, 1)

	assert.Len(t, people[0].Applications, 2)
	assert.Contains(t, people[0].Applications, "x")
	assert.Contains(t, people[0].Applications, "y")

	assert.Len(t, people[0].Emails, 4)
	assert.Contains(t, people[0].Emails, "a")
	assert.Contains(t, people[0].Emails, "b")
	assert.Contains(t, people[0].Emails, "c")
	assert.Contains(t, people[0].Emails, "d")

	// 3 into 1
	accounts = []*Account{
		{
			Application: "1",
			Emails:      []string{"c", "b"},
			Name:        "Person 1",
		},
		{
			Application: "2",
			Emails:      []string{"a"},
			Name:        "Person 2",
		},
		{
			Application: "3",
			Emails:      []string{"a", "b"},
			Name:        "Person 3",
		},
	}

	people = merge(accounts)

	assert.Len(t, people, 1)

	assert.Len(t, people[0].Applications, 3)
	assert.Contains(t, people[0].Applications, "1")
	assert.Contains(t, people[0].Applications, "2")
	assert.Contains(t, people[0].Applications, "3")

	assert.Len(t, people[0].Emails, 3)
	assert.Contains(t, people[0].Emails, "a")
	assert.Contains(t, people[0].Emails, "b")
	assert.Contains(t, people[0].Emails, "c")

	// 4 in 2
	accounts = []*Account{
		{
			Application: "1",
			Emails:      []string{"a@gmail.com", "b@gmail.com"},
			Name:        "A",
		},
		{
			Application: "1",
			Emails:      []string{"c@gmail.com", "d@gmail.com"},
			Name:        "B",
		},
		{
			Application: "2",
			Emails:      []string{"a@yahoo.com"},
			Name:        "C",
		},
		{
			Application: "3",
			Emails:      []string{"a@gmail.com", "a@yahoo.com"},
			Name:        "D",
		},
	}

	people = merge(accounts)

	assert.Len(t, people, 2)

	assert.Len(t, people[0].Applications, 3)
	assert.Contains(t, people[0].Applications, "1")
	assert.Contains(t, people[0].Applications, "2")
	assert.Contains(t, people[0].Applications, "3")

	assert.Len(t, people[0].Emails, 3)
	assert.Contains(t, people[0].Emails, "a@gmail.com")
	assert.Contains(t, people[0].Emails, "a@yahoo.com")
	assert.Contains(t, people[0].Emails, "b@gmail.com")

	assert.Len(t, people[1].Applications, 1)
	assert.Contains(t, people[1].Applications, "1")

	assert.Len(t, people[1].Emails, 2)
	assert.Contains(t, people[1].Emails, "c@gmail.com")
	assert.Contains(t, people[1].Emails, "d@gmail.com")
}

// TestProcessAccounts demonstrates that the accounts will be successfully processed/merged
// and the output will be a valid JSON.
// No exact JSON output match is assertion is done because contents of the input file may vary.
func TestProcessAccounts(t *testing.T) {
	output, err := ProcessAccounts("accounts.json")
	assert.NoError(t, err)

	var persons []Person
	err = json.Unmarshal([]byte(output), &persons)
	assert.NoError(t, err)
}

// TestProcessAccountsErrorReadingFile demonstrates that an error will be returned if the input file does not exist.
func TestProcessAccountsErrorReadingFile(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	_, err := ProcessAccounts("/tmp/does-not-exist-123")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to read the accounts")
}

// TODO more tests. TDD?
