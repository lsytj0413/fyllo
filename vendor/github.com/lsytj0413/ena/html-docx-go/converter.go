// Copyright (c) 2018 soren yang
//
// Licensed under the MIT License
// you may not use this file except in complicance with the License.
// You may obtain a copy of the License at
//
//     https://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package html2docx

import (
	"archive/zip"
	"fmt"
	"io"
)

type defConverter struct {
	option  Option
	sectors []sector
}

func (c *defConverter) ToDocx(html string, iow io.Writer) error {
	w := zip.NewWriter(iow)
	defer w.Close()

	for _, p := range c.sectors {
		err := p.Write(c, w, html)
		if err != nil {
			return fmt.Errorf("html2docx ToDocx failed: %s, %s", p.Name(), err.Error())
		}
	}

	return nil
}

func (c *defConverter) Options() Option {
	return c.option
}

// NewConverter will return a Converter instance
func NewConverter() (Converter, error) {
	option, err := NewOption()
	if err != nil {
		return nil, err
	}
	return NewConverterFromOption(option)
}

const (
	contentTypeName    = "[Content_Types].xml"
	contentTypeContent = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	  <Default Extension="rels" ContentType=
		"application/vnd.openxmlformats-package.relationships+xml" />
	  <Override PartName="/word/document.xml" ContentType=
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
	  <Override PartName="/word/afchunk.mht" ContentType="message/rfc822"/>
	</Types>
	`

	relsName    = "_rels/.rels"
	relsContent = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
	  <Relationship
		  Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
		  Target="/word/document.xml" Id="R09c83fafc067488e" />
	</Relationships>
	`

	relsDocumentName    = "word/_rels/document.xml.rels"
	relsDocumentContent = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
	  <Relationship Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/aFChunk"
		Target="/word/afchunk.mht" Id="htmlChunk" />
	</Relationships>`
)

// NewConverterFromOption will return a Converter instance with option
func NewConverterFromOption(option Option) (Converter, error) {
	c := &defConverter{option: option, sectors: make([]sector, 0, 5)}
	c.sectors = append(c.sectors, &defSector{
		name:    contentTypeName,
		content: contentTypeContent,
	})
	c.sectors = append(c.sectors, &defSector{
		name:    relsName,
		content: relsContent,
	})
	p, err := newDocumentXMLSector()
	if err != nil {
		return nil, err
	}
	c.sectors = append(c.sectors, p)

	p, err = newAfchunkSector()
	if err != nil {
		return nil, err
	}
	c.sectors = append(c.sectors, p)

	c.sectors = append(c.sectors, &defSector{
		name:    relsDocumentName,
		content: relsDocumentContent,
	})
	return c, nil
}
