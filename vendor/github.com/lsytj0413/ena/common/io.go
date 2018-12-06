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

package common

import (
	"io"
)

func read(r io.Reader, buf []byte) error {
	var readed int
	lens := cap(buf)
	for {
		n, err := r.Read(buf[readed:])
		if err != nil {
			return err
		}

		readed += n
		if readed == lens {
			break
		}
	}

	return nil
}

// Readn will read nbytes from Reader
func Readn(r io.Reader, n uint32) ([]byte, error) {
	buf := make([]byte, n)
	err := read(r, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Writen will write buf to Writer
func Writen(w io.Writer, buf []byte) error {
	var writed int
	len := len(buf)
	for {
		n, err := w.Write(buf[writed:])
		if err != nil {
			return err
		}

		writed += n
		if writed == len {
			break
		}
	}

	return nil
}
