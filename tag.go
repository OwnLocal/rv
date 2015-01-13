package rv

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var sourceMap = map[string]Source{
	"path":  PATH,
	"query": QUERY,
	"json":  JSON,
	"form":  FORM,
}

func ParseTag(tags string) (Args, error) {
	tagmap, err := splitTags(tags)
	if err != nil {
		return Args{}, err
	}
	tag, ok := tagmap["rv"]
	if !ok {
		return Args{}, nil
	}
	return parseRvTag(tag)
}

var (
	wordRe     = regexp.MustCompile(`[^.\s=]+`)
	notSpace   = regexp.MustCompile(`[^\s]*`)
	maybeSpace = regexp.MustCompile(`\s*`)
	keyRe      = regexp.MustCompile(`[^:]+`)
	valRe      = regexp.MustCompile(`"(\\.|[^"])*"`)
)

func parseRvTag(tag string) (Args, error) {
	args := Args{}
	l := lexer{input: tag}

	args.Source = sourceMap[l.match(wordRe)]
	l.matchString(".")
	args.Field = l.match(wordRe)

	for !l.eof && l.err == nil {
		l.match(maybeSpace)
		opt := l.match(wordRe)
		l.matchString("=")
		val := l.match(notSpace)
		if l.err == nil {
			args.fillOpt(opt, val)
		}
	}

	return args, l.err
}

func (a *Args) fillOpt(opt, val string) error {
	var err error

	switch opt {

	case "range":
		parts := strings.Split(val, ",")
		if len(parts) < 2 {
			return fmt.Errorf("range expects 2 comma-separated ints, got '%s'", val)
		}
		a.Range.Start, err = strconv.Atoi(parts[0])
		if err == nil {
			a.Range.End, err = strconv.Atoi(parts[1])
		}
		return err

	case "options":
		if len(val) < 1 {
			return fmt.Errorf("options can't be empty")
		}
		a.Options = strings.Split(val, ",")

	case "default":
		a.Default = val

	}

	return nil
}

func splitTags(tags string) (map[string]string, error) {
	tagmap := map[string]string{}

	l := lexer{input: tags}
	for !l.eof && l.err == nil {
		l.match(maybeSpace)
		key := l.match(keyRe)
		l.matchString(":")
		val, _ := strconv.Unquote(l.match(valRe))

		if key != "" {
			tagmap[key] = val
		}
	}
	return tagmap, l.err
}

type lexer struct {
	input string
	err   error
	eof   bool
}

func (l *lexer) match(r *regexp.Regexp) string {
	if l.err != nil {
		return ""
	}

	return l.handleMatch(r.FindStringIndex(l.input), r.String())
}

func (l *lexer) matchString(s string) string {
	if l.err != nil {
		return ""
	}

	i := strings.Index(l.input, s)
	return l.handleMatch([]int{i, i + len(s)}, s)
}

func (l *lexer) handleMatch(match []int, matcher string) string {
	if match == nil || match[0] != 0 {
		l.err = fmt.Errorf("Failed match '%s' in '%s'", matcher, l.input)
		return ""
	} else {
		token := l.input[0:match[1]]
		l.input = l.input[match[1]:]
		if len(l.input) == 0 {
			l.eof = true
		}
		return token
	}
}
