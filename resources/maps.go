package resources

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

var maps = make([]Map, 0)

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
			t := scanner.Text()
			if t[0] == '!' {
				m.Name = strings.TrimSpace(t[1:])
			} else {
				if len(t) > m.Columns {
					m.Columns = len(t)
				}
				m.Rows++
				m.Cells = append(m.Cells, t)
			}
		}
		// Now read our runes.
		maps = append(maps, m)
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

// GetNextMap bogusly gets the next map from the provided one.
func GetNextMap(m Map) Map {
	nextIsIt := false
	for _, v := range maps {
		if nextIsIt {
			return v
		}
		if v.Name == m.Name {
			nextIsIt = true
		}
	}
	return Map{}
}
