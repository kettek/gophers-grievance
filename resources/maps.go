package resources

import (
	"bufio"
	"bytes"
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"
	"strings"
)

var maps = make([]Map, 0)

type Map struct {
	Name          string
	Background    color.RGBA
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
			Name:       strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
			Background: color.RGBA{128, 128, 128, 255},
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
			} else if t[0] == '@' {
				parts := strings.Split(strings.TrimSpace(t[1:]), ",")
				for i, p := range parts {
					p = strings.TrimSpace(p)
					if i == 0 {
						if c, err := strconv.ParseInt(p, 10, 32); err == nil {
							m.Background.R = uint8(c)
						} else {
							fmt.Println(err)
						}
					} else if i == 1 {
						if c, err := strconv.ParseInt(p, 10, 32); err == nil {
							m.Background.G = uint8(c)
						} else {
							fmt.Println(err)
						}
					} else if i == 2 {
						if c, err := strconv.ParseInt(p, 10, 32); err == nil {
							m.Background.B = uint8(c)
						} else {
							fmt.Println(err)
						}
					}
				}
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
