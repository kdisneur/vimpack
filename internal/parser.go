package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

type Parser struct {
	err       error
	namespace *Namespace
	plugins   []*Plugin
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(filepath string) ([]*Plugin, error) {
	filereader, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("can't open %s: %s", filepath, err)
	}

	return p.Parse(filepath, filereader)
}

func (p *Parser) Parse(filename string, config io.Reader) ([]*Plugin, error) {
	var s scanner.Scanner
	s.Init(config)
	s.Filename = filename
	s.Mode = scanner.ScanComments | scanner.SkipComments | scanner.ScanIdents | scanner.ScanStrings
	// Do not recognize `\n` as Whitespace
	s.Whitespace ^= 1 << '\n'

	p.namespace = nil
	p.err = nil
	p.plugins = []*Plugin{}

	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		switch s.TokenText() {
		case "\n":
			continue
		case "namespace":
			p.parseNamespace(&s, token)
		case "onstart", "ondemand":
			p.parsePlugin(&s, token)
		default:
			p.err = fmt.Errorf("%s: unexpected token: %s", s.Position, s.TokenText())
		}

		if p.err != nil {
			return nil, p.err
		}
	}

	return p.plugins, p.err
}

func (p *Parser) parsePlugin(s *scanner.Scanner, token rune) {
	var loading Loading
	if s.TokenText() == "onstart" {
		loading = LoadingStart
	} else if s.TokenText() == "ondemand" {
		loading = LoadingOnDemand
	} else {
		p.err = fmt.Errorf("%s: expected onstart or ondemand got: %s", s.Position, s.TokenText())
		return
	}

	if p.namespace == nil {
		p.err = fmt.Errorf("%s: expected onstart to be in a namespace", s.Position)
		return
	}

	token = s.Scan()
	if token != scanner.String {
		p.err = fmt.Errorf("%s: expected string got: %s", s.Position, s.TokenText())
		return
	}

	name := strings.Trim(s.TokenText(), "\"")
	if name == "" {
		p.err = fmt.Errorf("%s: repository can't be empty", s.Position)
		return
	}

	token = s.Scan()
	if token != '\n' {
		p.err = fmt.Errorf("%s: expected end of line got: %s", s.Position, s.TokenText())
		return
	}

	names := strings.Split(name, "/")
	if len(names) != 2 {
		p.err = fmt.Errorf("%s: expected repositoy name to be '<repository>/<name>', got: %s", s.Position, name)
		return
	}

	plugin := NewPlugin(names[1], *p.namespace, NewGitHub(name))
	plugin.Loading = loading
	p.plugins = append(p.plugins, plugin)
}

func (p *Parser) parseNamespace(s *scanner.Scanner, token rune) {
	token = s.Scan()
	if token != scanner.String {
		p.err = fmt.Errorf("%s: expected string got: %s", s.Position, s.TokenText())
		return
	}

	namespace := strings.Trim(s.TokenText(), "\"")
	if namespace == "" {
		p.err = fmt.Errorf("%s: namespace can't be empty", s.Position)
		return
	}

	token = s.Scan()
	if token != '\n' {
		p.err = fmt.Errorf("%s: expected end of line got: %s", s.Position, s.TokenText())
		return
	}

	ns := Namespace(namespace)
	p.namespace = &ns
}
