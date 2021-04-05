package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"istio.io/istio/pilot/test/util"
)

func TestGolden(t *testing.T) {
	cases := []string{"input"}
	for _, tt := range cases {
		t.Run(tt, func(t *testing.T) {
			inp, err := os.ReadFile(fmt.Sprintf("testdata/%s.yaml", tt))
			if err != nil {
				t.Fatal(err)
			}
			out := capture(func() {
				fmt.Println(runMigration(inp))
			})
			goldenFile := fmt.Sprintf("testdata/%s.yaml.golden", tt)
			if util.Refresh() {
				if err := ioutil.WriteFile(goldenFile, []byte(out), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			expected, err := ioutil.ReadFile(goldenFile)
			if err != nil {
				t.Fatal(err)
			}
			if out != string(expected) {
				t.Fatalf("expected %v, got %v", string(expected), out)
			}
		})
	}
}

func capture(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	stderr := os.Stderr
	os.Stderr = w
	defer func() {
		os.Stderr = stderr
	}()

	lw := log.Writer()
	log.SetOutput(w)
	defer func() {
		log.SetOutput(lw)
	}()

	f()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
