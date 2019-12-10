package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

// siteDir is the target directory into which the HTML gets generated. Its
// default is set here but can be changed by an argument passed into the
// program.
const syslRoot = "../../"
const siteDir = syslRoot + "docs/website/content/docs/byexample/"
const assetDir = syslRoot + "docs/website/static/assets/byexample/"
const pygmentizeBin = syslRoot + "docs/website/byexample/vendor/pygments/pygmentize"
const templates = syslRoot + "docs/website/byexample/templates"
const cacheDir = "./.tmp/gobyexample-cache"
const orderingfile = "ordering.yaml"

var imageFiles = []string{".png", ".svg"}

func main() {
	ensureDir(siteDir)
	examples := parseExamples()
	renderExamples(examples)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ensureDir(dir string) {
	err := os.MkdirAll(dir, 0755)
	check(err)
}

func copyFile(src, dst string) {
	dat, err := ioutil.ReadFile(src)
	check(err)
	err = ioutil.WriteFile(dst, dat, 0644)
	check(err)
}

func pipe(bin string, arg []string, src string) []byte {
	fmt.Println(bin, arg)

	cmd := exec.Command(bin, arg...)
	in, err := cmd.StdinPipe()
	check(err)
	out, err := cmd.StdoutPipe()
	check(err)
	err = cmd.Start()
	check(err)
	_, err = in.Write([]byte(src))
	check(err)
	err = in.Close()
	check(err)
	bytes, err := ioutil.ReadAll(out)
	check(err)
	err = cmd.Wait()
	check(err)
	return bytes
}

func sha1Sum(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	b := h.Sum(nil)
	return fmt.Sprintf("%x", b)
}

func mustReadFile(path string) string {
	bytes, err := ioutil.ReadFile(path)
	check(err)
	return string(bytes)
}

func cachedPygmentize(lex string, src string) string {
	ensureDir(cacheDir)
	arg := []string{"-l", lex, "-f", "html"}
	cachePath := cacheDir + "/pygmentize-" + strings.Join(arg, "-") + "-" + sha1Sum(src)
	// fmt.Println(cachePath)
	fmt.Println(arg)

	cacheBytes, cacheErr := ioutil.ReadFile(cachePath)
	if cacheErr == nil {
		return string(cacheBytes)
	}
	// code := ""
	// renderBytes := pipe(pygmentizeBin, arg, src)
	// for _, i := range
	code := chromaFormat(src)

	return code
	// Newer versions of Pygments add silly empty spans.
	// renderCleanString := strings.Replace(string(renderBytes), "<span></span>", "", -1)
	// writeErr := ioutil.WriteFile(cachePath, []byte(renderCleanString), 0600)
	// check(writeErr)
	// return renderCleanString
}
func chromaFormat(code string) string {
	lexer := lexers.Get("go")
	if lexer == nil {
		panic("")
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("swapoff")
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.WithClasses(true))
	// if formatter == nil {
	// 	formatter = formatters.Fallback
	// }

	iterator, err := lexer.Tokenise(nil, string(code))
	check(err)
	buf := new(bytes.Buffer)
	err = formatter.Format(buf, style, iterator)
	check(err)
	return buf.String()

}
func markdown(src string) string {
	return string(blackfriday.Run([]byte(src)))
}

func readLines(path string) []string {
	src := mustReadFile(path)
	return strings.Split(src, "\n")
}

func mustGlob(glob string) []string {
	paths, err := filepath.Glob(glob)
	check(err)
	return paths
}

func whichLexer(path string) string {
	if strings.HasSuffix(path, ".sysl") {
		return "sysl"
	} else if strings.HasSuffix(path, ".sh") {
		return "console"
	} else if strings.HasSuffix(path, ".png") {
		return "png"
	}
	panic("No lexer for " + path)
}

func debug(msg string) {
	if os.Getenv("DEBUG") == "1" {
		fmt.Fprintln(os.Stderr, msg)
	}
}

var docsPat = regexp.MustCompile("^\\s*(\\/\\/|#)\\s")
var dashPat = regexp.MustCompile("\\-+")

// Seg is a segment of an example
type Seg struct {
	Docs, DocsRendered              string
	Code, CodeRendered, CodeForJs   string
	CodeEmpty, CodeLeading, CodeRun bool
}

// Example is info extracted from an example file
type Example struct {
	ID, Name                    string
	Topic                       string
	Weight                      int
	Images                      []string
	GoCode, GoCodeHash, URLHash string
	Segs                        [][]*Seg
	PrevExample                 *Example
	NextExample                 *Example
}

func parseSegs(sourcePath string) ([]*Seg, string) {
	var lines []string
	// Convert tabs to spaces for uniform rendering.
	for _, line := range readLines(sourcePath) {
		lines = append(lines, strings.Replace(line, "\t", "    ", -1))
	}
	filecontent := strings.Join(lines, "\n")
	segs := []*Seg{}
	lastSeen := ""
	for _, line := range lines {
		if line == "" {
			lastSeen = ""
			continue
		}
		matchDocs := docsPat.MatchString(line)
		matchCode := !matchDocs
		newDocs := (lastSeen == "") || ((lastSeen != "docs") && (segs[len(segs)-1].Docs != ""))
		newCode := (lastSeen == "") || ((lastSeen != "code") && (segs[len(segs)-1].Code != ""))
		if newDocs || newCode {
			debug("NEWSEG")
		}
		if matchDocs {
			trimmed := docsPat.ReplaceAllString(line, "")
			if newDocs {
				newSeg := Seg{Docs: trimmed, Code: ""}
				segs = append(segs, &newSeg)
			} else {
				segs[len(segs)-1].Docs = segs[len(segs)-1].Docs + trimmed
			}
			debug("DOCS: " + line)
			lastSeen = "docs"
		} else if matchCode {
			if newCode {
				newSeg := Seg{Docs: "", Code: line}
				segs = append(segs, &newSeg)
			} else {
				segs[len(segs)-1].Code = segs[len(segs)-1].Code + "\n" + line
			}
			debug("CODE: " + line)
			lastSeen = "code"
		}
	}
	for i, seg := range segs {
		seg.CodeEmpty = (seg.Code == "")
		seg.CodeLeading = (i < (len(segs) - 1))
		seg.CodeRun = i == 1
	}
	return segs, filecontent
}

func parseAndRenderSegs(sourcePath string) ([]*Seg, string) {
	segs, filecontent := parseSegs(sourcePath)
	lexer := whichLexer(sourcePath)
	for _, seg := range segs {
		if seg.Docs != "" {
			seg.DocsRendered = markdown(seg.Docs)
		}
		if seg.Code != "" {
			seg.CodeRendered = cachedPygmentize(lexer, seg.Code)

			// adding the content to the js code for copying to the clipboard
			if strings.HasSuffix(sourcePath, ".sysl") {
				seg.CodeForJs = strings.Trim(seg.Code, "\n") + "\n"
			}
		}
	}
	if lexer != "sysl" {
		filecontent = ""
	}
	return segs, filecontent
}

// unmarshalYaml unmarshals a yaml file of form
// key1:
// 		- value 1
// 		- value 2
func unmarshalYaml(filename string) map[string][]string {
	source, err := ioutil.ReadFile(filename)
	check(err)
	this := make(map[string][]string)
	err = yaml.Unmarshal(source, this)
	check(err)
	return this
}

func parseExamples() []*Example {
	var examples []*Example
	ordering := unmarshalYaml(orderingfile)
	weight := 0
	for topic, tutorial := range ordering {
		for _, exampleName := range tutorial {
			weight++
			example := Example{Name: exampleName}
			exampleID := strings.ToLower(exampleName)
			exampleID = strings.Replace(exampleID, " ", "-", -1)
			exampleID = strings.Replace(exampleID, "/", "-", -1)
			exampleID = strings.Replace(exampleID, "'", "", -1)
			exampleID = dashPat.ReplaceAllString(exampleID, "-")
			example.ID = exampleID
			example.Weight = weight + 1
			example.Topic = topic
			example.Segs = make([][]*Seg, 0)
			sourcePaths := mustGlob(exampleID + "/*")

			for _, sourcePath := range sourcePaths {
				if ok, ext := isImageFile(sourcePath); ok {
					destination := assetDir + "images/" + exampleID + strconv.Itoa(weight) + ext
					copyFile(sourcePath, destination)

					// This is the path that gets rendered in the markdown file
					imagesRelativeToSite := "/assets/byexample/images/"

					example.Images = append(example.Images, imagesRelativeToSite+exampleID+strconv.Itoa(weight)+ext)
				} else {
					sourceSegs, filecontents := parseAndRenderSegs(sourcePath)
					if filecontents != "" {
						example.GoCode = filecontents
					}
					example.Segs = append(example.Segs, sourceSegs)
				}
			}
			newCodeHash := sha1Sum(example.GoCode)
			if example.GoCodeHash != newCodeHash {

			}
			examples = append(examples, &example)
		}
	}

	for i, example := range examples {
		if i > 0 {
			example.PrevExample = examples[i-1]
		}
		if i < (len(examples) - 1) {
			example.NextExample = examples[i+1]
		}
	}

	return examples
}

func renderExamples(examples []*Example) {
	exampleTmpl := template.New("example")
	_, err := exampleTmpl.Parse(mustReadFile(templates + "/example.tmpl"))
	check(err)
	for _, example := range examples {
		exampleF, err := os.Create(siteDir + example.ID + ".md")
		check(err)
		exampleTmpl.Execute(exampleF, example)
	}
}
func isImageFile(filename string) (bool, string) {
	for _, extension := range imageFiles {
		if strings.HasSuffix(filename, extension) {
			return true, extension
		}
	}
	return false, ""
}
