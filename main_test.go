package main

import (
	"os"
	"testing"
)

func TestVscmrVerify(t *testing.T) {
	realArgs := os.Args
	defer func() {
		os.Args = realArgs
	}()

	os.Args = []string{"", "vscmr", "verify", "ue/test-data/o00407-0100700090001.zip"}

	main()
}

func TestBuVerify(t *testing.T) {
	realArgs := os.Args
	defer func() {
		os.Args = realArgs
	}()

	os.Args = []string{"", "bu", "verify", "ue/test-data/o00407-0100700090001.zip"}
	main()

	os.Args = []string{"", "bu", "verify", "ue/test-data/urna.bu"}
	main()
}

func TestBuCount(t *testing.T) {
	realArgs := os.Args
	defer func() {
		os.Args = realArgs
	}()

	os.Args = []string{"", "bu", "count", "-cargo", "Presidente", "ue/test-data/o00407-0100700090001.zip"}
	main()

	os.Args = []string{"", "bu", "count", "-cargo", "Presidente", "ue/test-data/urna.bu"}
	main()
}

func TestBuCsv(t *testing.T) {
	realArgs := os.Args
	defer func() {
		os.Args = realArgs
	}()

	os.Args = []string{"", "bu", "csv", "-candidatos", "Nulo,Branco", "-cargo", "Presidente", "ue/test-data/o00407-0100700090001.zip"}
	main()

	os.Args = []string{"", "bu", "csv", "-candidatos", "Nulo,Branco", "-cargo", "Presidente", "ue/test-data/urna.bu"}
	main()
}

func TestRdvCsv(t *testing.T) {
	realArgs := os.Args
	defer func() {
		os.Args = realArgs
	}()

	os.Args = []string{"", "rdv", "csv", "ue/test-data/o00407-0100700090001.zip"}
	main()

	os.Args = []string{"", "rdv", "csv", "ue/test-data/urna.rdv"}
	main()
}
