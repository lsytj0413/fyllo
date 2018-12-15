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
	"github.com/gin-gonic/gin"
	"github.com/lsytj0413/ena/cerror"
	"github.com/lsytj0413/ena/logger"

	ierror "github.com/lsytj0413/fyllo/pkg/error"
)

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if nerr := recover(); nerr != nil {
				if err, ok := nerr.(error); ok {
					var cerr *cerror.Error
					var ok bool
					if cerr, ok = err.(*cerror.Error); !ok {
						cerr = ierror.NewError(ierror.EcodeUnknown, err.Error())
					}
					cerr.WriteTo(c.Writer)
				}
			}
		}()

		c.Next()
	}
}

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Infof("Accept Request, url=%s", c.Request.URL.String())

		c.Next()
	}
}

func jsonRespMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")

		c.Next()
	}
}
