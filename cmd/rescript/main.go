package main

import (
	//	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/akeil/rmtool"
	"github.com/akeil/rmtool/pkg/api"
	"golang.org/x/sync/errgroup"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"

	"github.com/akeil/rescript"
)

const (
	checkmark = "\u2713"
	crossmark = "\u2717"
	ellipsis  = "\u2026"
)

var langs = map[string]rescript.LanguageCode{
	"en": rescript.LangEN,
	"de": rescript.LangDE,
}

func main() {
	app := kingpin.New("hwr", "reMarkable Handwriting Recogntion")
	app.HelpFlag.Short('h')

	var (
		name   = app.Arg("name", "Name of the notebook to convert").Required().String()
		dst    = app.Arg("dir", "Directory for output document").Default(".").String()
		format = app.Flag("format", "Output format").Short('f').Default("txt").Enum("txt", "md")
		lang   = app.Flag("lang", "Language of the notebook").Short('l').Default("en").String()
	)

	kingpin.MustParse(app.Parse(os.Args[1:]))

	rmtool.SetLogLevel("error")

	err := run(*name, *dst, *lang, *format)
	if err != nil {
		fmt.Printf("%v Error: %v\n", crossmark, err)
		os.Exit(1)
	}

	fmt.Printf("%v Done.\n", checkmark)
}

func run(name, dst, lang, format string) error {
	lc, ok := langs[lang]
	if !ok {
		return fmt.Errorf("invalid language %q", lang)
	}

	s, err := loadSettings()
	if err != nil {
		return err
	}

	rec := rescript.NewRecognizer(s.AppKey, s.HmacKey, s.hwrCache())

	c, err := initClient(s)
	if err != nil {
		return err
	}

	r := api.NewRepository(c, s.CacheDir)

	items, err := r.List()
	if err != nil {
		return err
	}
	root := rmtool.BuildTree(items)
	root = root.Filtered(rmtool.IsDocument, rmtool.MatchName(name))

	cmp := selectComposer(format)

	// do recognition for each matching document
	var group errgroup.Group
	root.Walk(func(n *rmtool.Node) error {
		if n.Type() == rmtool.CollectionType {
			return nil
		}

		group.Go(func() error {
			fmt.Printf("%v download notebook %q\n", ellipsis, n.Name())
			doc, err := rmtool.ReadDocument(r, n)
			if err != nil {
				return err
			}

			fmt.Printf("%v recognize handwriting for %q\n", ellipsis, n.Name())
			results, err := rec.Recognize(doc, lc)
			if err != nil {
				return err
			}

			path := filepath.Join(dst, doc.Name()+"."+format)
			f, err := os.Create(path)
			if err != nil {
				return nil
			}
			defer f.Close()

			err = cmp(f, doc, results)
			if err != nil {
				return err
			}
			fmt.Printf("%v write %q to %q\n", checkmark, n.Name(), path)
			return nil
		})
		return nil
	})
	return group.Wait()
}

func initClient(s settings) (*api.Client, error) {
	token, err := loadToken(s.tokenPath())
	if err != nil {
		return nil, err
	}
	client := api.DefaultClient(token)

	err = register(s, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func register(s settings, c *api.Client) error {
	if c.IsRegistered() {
		return nil
	}

	code, err := readInput("Enter one time code from https://my.remarkable.com/")
	if err != nil {
		return err
	}

	token, err := c.Register(code)
	if err != nil {
		return err
	}

	err = saveToken(s.tokenPath(), token)
	if err != nil {
		return err
	}

	return nil
}

func readInput(msg string) (string, error) {
	var reply string

	fmt.Printf("%v: \n", msg)
	_, err := fmt.Scanf("%s", &reply)

	return reply, err
}

func selectComposer(t string) rescript.ComposeFunc {
	switch t {
	case "txt":
		return rescript.NewPlaintextComposer()
	case "md":
		return rescript.NewMarkdownComposer()
	default:
		return rescript.NewPlaintextComposer()
	}
}

func loadToken(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(d), err
}

func saveToken(path, token string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(token))
	return err
}

type settings struct {
	DataDir  string
	CacheDir string
	AppKey   string
	HmacKey  string
}

func (s settings) tokenPath() string {
	return filepath.Join(s.DataDir, "device-token")
}

func (s settings) hwrCache() string {
	return filepath.Join(s.CacheDir, "hwr")
}

func loadSettings() (settings, error) {
	s := settings{}
	config, err := os.UserConfigDir()
	if err != nil {
		return s, err
	}

	path := filepath.Join(config, "rmhwr-conf.yaml")
	f, err := os.Open(path)
	if err != nil {
		return s, err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&s)
	if err != nil {
		return s, err
	}

	return s, nil
}
