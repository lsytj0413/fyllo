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
	"bytes"
	"text/template"

	"github.com/lsytj0413/ena/common"
)

// documentXMLSector will write word/document.xml file
type documentXMLSector struct {
	name    string
	content string
	tmpl    *template.Template
}

func (p *documentXMLSector) Name() string {
	return p.name
}

func (p *documentXMLSector) Write(c Converter, w *zip.Writer, html string) error {
	f, err := w.CreateHeader(&zip.FileHeader{
		Name:   p.name,
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})
	err = p.tmpl.Execute(buf, c.Options())
	if err != nil {
		return err
	}

	return common.Writen(f, buf.Bytes())
}

const (
	documentXMLName    = "word/document.xml"
	documentXMLContent = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<w:document
	  xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
	  xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math"
	  xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
	  xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
	  xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
	  xmlns:ns6="http://schemas.openxmlformats.org/schemaLibrary/2006/main"
	  xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart"
	  xmlns:ns8="http://schemas.openxmlformats.org/drawingml/2006/chartDrawing"
	  xmlns:dgm="http://schemas.openxmlformats.org/drawingml/2006/diagram"
	  xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture"
	  xmlns:ns11="http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing"
	  xmlns:dsp="http://schemas.microsoft.com/office/drawing/2008/diagram"
	  xmlns:ns13="urn:schemas-microsoft-com:office:excel"
	  xmlns:o="urn:schemas-microsoft-com:office:office"
	  xmlns:v="urn:schemas-microsoft-com:vml"
	  xmlns:w10="urn:schemas-microsoft-com:office:word"
	  xmlns:ns17="urn:schemas-microsoft-com:office:powerpoint"
	  xmlns:odx="http://opendope.org/xpaths"
	  xmlns:odc="http://opendope.org/conditions"
	  xmlns:odq="http://opendope.org/questions"
	  xmlns:odi="http://opendope.org/components"
	  xmlns:odgm="http://opendope.org/SmartArt/DataHierarchy"
	  xmlns:ns24="http://schemas.openxmlformats.org/officeDocument/2006/bibliography"
	  xmlns:ns25="http://schemas.openxmlformats.org/drawingml/2006/compatibility"
	  xmlns:ns26="http://schemas.openxmlformats.org/drawingml/2006/lockedCanvas">
	  <w:body>
		<w:altChunk r:id="htmlChunk" />
		<w:sectPr>
		  <w:pgSz w:w="{{.Width}}" w:h="{{.Height}}" w:orient="{{.Orient}}" />
		  <w:pgMar w:top="{{.Margin.Top}}"
				   w:right="{{.Margin.Right}}"
				   w:bottom="{{.Margin.Bottom}}"
				   w:left="{{.Margin.Left}}"
				   w:header="{{.Margin.Header}}"
				   w:footer="{{.Margin.Footer}}"
				   w:gutter="{{.Margin.Gutter}}"/>
		</w:sectPr>
	  </w:body>
	</w:document>
	`
)

func newDocumentXMLSector() (sector, error) {
	var err error

	p := &documentXMLSector{
		name:    documentXMLName,
		content: documentXMLContent,
	}

	p.tmpl, err = template.New("documentXML").Parse(p.content)
	if err != nil {
		return nil, err
	}

	return p, nil
}
