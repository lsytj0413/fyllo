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

package fyllo

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/lsytj0413/ena/logger"
)

func wrapperHandler(f func(*gin.Context) (interface{}, error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := f(c)
		if err != nil {
			logger.Errorf("handler failed, %v", err)
			panic(err)
		}

		data, err := json.Marshal(resp)
		if err != nil {
			logger.Errorf("handler marshal[%v] failed, %v", resp, err)
			panic(err)
		}

		_, err = c.Writer.Write(data)
		if err != nil {
			logger.Errorf("handler write response[%v] failed, %v", string(data), err)
			panic(err)
		}
	}
}
