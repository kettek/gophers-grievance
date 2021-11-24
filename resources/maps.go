package resources

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

var maps = map[string]Map{}

type Map struct {
	Name          string
	Columns, Rows int
	Cells         []string
}

func loadMaps() error {
	results, err := f.ReadDir("maps")
	if err != nil {
		return err
	}
	for _, file := range results {
		if file.IsDir() {
			continue
		}
		m := Map{
			Name: strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
		}
		data, err := f.ReadFile("maps/" + file.Name())

		if err != nil {
			fmt.Errorf("%w", err)
			continue
		}
		// This is inefficent, but get our dimensions first.
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			if len(scanner.Text()) > m.Columns {
				m.Columns = len(scanner.Text())
			}
			m.Rows++
			m.Cells = append(m.Cells, scanner.Text())
		}
		// Now read our runes.
		maps[m.Name] = m
	}

	return nil
}

// GetAnyMap returns any map. If no map exists, it returns an empty map.
func GetAnyMap() Map {
	for _, v := range maps {
		return v
	}
	return Map{}
}
