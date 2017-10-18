package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

const defaultURL = "https://docs.aws.amazon.com/ja_jp/general/latest/gr/glos-chap.html"

type userConfig struct {
	URL   string
	Rules []userRule
}

type userRule struct {
	Expected string
	Patterns []string
}

func main() {
	url := defaultURL
	userConf, err := loadConfig()
	if err == nil {
		fmt.Printf("%+v\n\n", userConf)
		url = userConf.URL
	}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatalf("+%v\n", err)
	}

	wordsSet := make(map[string]struct{})
	doc.Find("dt").Each(func(_ int, s *goquery.Selection) {
		addWords(s.Text(), wordsSet)
	})
	fmt.Printf("%+v\n", wordsSet)
	saveYamlFile(wordsSet, userConf)
}

func addWords(word string, wordsSet map[string]struct{}) {
	text := strings.TrimSpace(word)
	if text != "" {
		for _, v := range strings.Split(text, "(") {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			for _, z := range strings.Split(v, ")") {
				z = strings.TrimSpace(z)
				if z == "" {
					continue
				}
				if utf8.RuneCountInString(z) <= 1 {
					continue
				}
				wordsSet[z] = struct{}{}

				for _, x := range strings.Split(z, " ") {
					x = strings.TrimSpace(x)
					if x == "" {
						continue
					}
					if utf8.RuneCountInString(x) <= 1 {
						continue
					}
					wordsSet[x] = struct{}{}
				}
			}
		}
	}
}

func saveYamlFile(wordsSet map[string]struct{}, config *userConfig) error {
	path, err := getYamlPath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	words := make([]string, 0, len(wordsSet)+1)
	for k := range wordsSet {
		words = append(words, k)
	}
	if len(words) > 0 {
		sort.Slice(words, func(i, j int) bool {
			return utf8.RuneCountInString(words[i]) > utf8.RuneCountInString(words[j])
		})
	}

	// write file
	text := "version: 1\n"
	_, err = f.WriteString(text)
	if err != nil {
		return err
	}

	text = "rules:\n"
	_, err = f.WriteString(text)
	if err != nil {
		return err
	}

	for _, v := range words {
		text = strings.TrimSpace(v)

		if text == "" {
			continue
		}
		text = fmt.Sprintf("  - expected: '%s'\n", text)

		_, err = f.WriteString(text)
		if err != nil {
			return err
		}

		sl := getPatterns(config, v)
		if len(sl) > 0 {
			text = "    patterns:\n"
			_, err = f.WriteString(text)
			if err != nil {
				return err
			}

			for _, z := range sl {
				text = strings.TrimSpace(escapeText(z))
				text = fmt.Sprintf("      - '%s'\n", text)

				_, err = f.WriteString(text)
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}

func escapeText(text string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`*`, `\*`,
		`+`, `\+`,
		`.`, `\.`,
		`?`, `\?`,
		`{`, `\{`,
		`}`, `\}`,
		`(`, `\(`,
		`)`, `\)`,
		`[`, `\[`,
		`]`, `\]`,
		`|`, `\|`,
		`^`, `\^`,
		`-`, `\-`,
		`$`, `\$`)

	return replacer.Replace(text)
}

func getPatterns(config *userConfig, exp string) []string {
	var sl []string

	expText := strings.TrimSpace(exp)
	if strings.Index(expText, " ") > 0 {
		sl = append(sl, strings.Replace(expText, " ", "", -1))
	}

	if config != nil {
		for _, v := range config.Rules {
			if expText == v.Expected {
				//sl = append(sl, v.Patterns...)
				for _, z := range v.Patterns {
					z = strings.TrimSpace(z)
					if z == "" {
						continue
					}
					sl = append(sl, z)
				}
			}
		}
	}

	if len(sl) > 0 {
		sort.Slice(sl, func(i, j int) bool {
			return utf8.RuneCountInString(sl[i]) > utf8.RuneCountInString(sl[j])
		})
	}

	return sl
}

func getExecDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(execPath), nil
}

func getYamlPath() (string, error) {
	dir, err := getExecDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "aws_words.yml"), err
}

func getConfigPath() (string, error) {
	dir, err := getExecDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), err
}

func loadConfig() (*userConfig, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &userConfig{}
	err = json.NewDecoder(f).Decode(conf)

	return conf, err
}

func saveConfig(config *userConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(config)
	if err != nil {
		return err
	}

	return err
}
