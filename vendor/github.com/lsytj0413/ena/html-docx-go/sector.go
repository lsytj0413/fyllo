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

	"github.com/lsytj0413/ena/common"
)

// sector is process while DoDocx is called
type sector interface {
	Write(c Converter, w *zip.Writer, html string) error
	Name() string
}

// defSector will write content to filename
type defSector struct {
	name    string
	content string
}

func (p *defSector) Write(c Converter, w *zip.Writer, html string) error {
	f, err := w.CreateHeader(&zip.FileHeader{
		Name:   p.name,
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	return common.Writen(f, []byte(p.content))
}

func (p *defSector) Name() string {
	return p.name
}
