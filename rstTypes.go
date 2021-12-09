package main

type RstDirectiveMap struct {
	Files []RstFile
}

func NewRstDirectiveMap() *RstDirectiveMap {
	return &RstDirectiveMap{}
}

func (rdm *RstDirectiveMap) AddConstants(refmap RefMap, fileName, input string) *RstDirectiveMap {
	var rstFile RstFile
	rstFile.Name = fileName
	rstFile.Links = make([]HTTPLink, 0)
	return rdm
}
