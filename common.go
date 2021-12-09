package main

type RstFile struct {
	Name  string
	Links []HTTPLink
}

type HTTPLink string
