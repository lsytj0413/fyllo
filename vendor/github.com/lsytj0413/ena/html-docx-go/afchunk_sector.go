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
	"strings"
	"text/template"

	"github.com/lsytj0413/ena/common"
)

// afchunkSector will write word/afchunk.mht file
type afchunkSector struct {
	name    string
	content string
	tmpl    *template.Template
}

func (p *afchunkSector) Name() string {
	return p.name
}

func (p *afchunkSector) Write(c Converter, w *zip.Writer, html string) error {
	f, err := w.CreateHeader(&zip.FileHeader{
		Name:   p.name,
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	htmlSource := strings.Replace(html, "=", "=3D", -1)

	buf := bytes.NewBuffer([]byte{})
	err = p.tmpl.Execute(buf, htmlSource)
	if err != nil {
		return err
	}

	return common.Writen(f, buf.Bytes())
}

const (
	afchunkName    = "word/afchunk.mht"
	afchunkContent = `MIME-Version: 1.0
	Content-Type: multipart/related;
		type="text/html";
		boundary="----=mhtDocumentPart"
	
	
	------=mhtDocumentPart
	Content-Type: text/html; charset="utf-8";
	Content-Transfer-Encoding: quoted-printable
	Content-Location: file:///C:/fake/document.html
	
	{{.}}
	
	------=mhtDocumentPart--
	`
)

func newAfchunkSector() (sector, error) {
	var err error

	p := &afchunkSector{
		name:    afchunkName,
		content: afchunkContent,
	}

	p.tmpl, err = template.New("afchunk").Parse(p.content)
	if err != nil {
		return nil, err
	}

	return p, nil
}
