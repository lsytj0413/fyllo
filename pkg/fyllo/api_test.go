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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"github.com/lsytj0413/fyllo/pkg/server"
)

type apiTestSuite struct {
	suite.Suite
}

func getRouter(t *testing.T) *gin.Engine {
	r := gin.Default()
	installers := []server.Installer{
		&versionService{},
		&segmentService{},
		&randomService{},
		&snowflakeService{},
	}

	for _, installer := range installers {
		err := installer.Install(r)
		if err != nil {
			t.Fatalf("install failed, %v", err)
		}
	}

	return r
}

func (s *apiTestSuite) TestVersionOk() {
	w := httptest.NewRecorder()
	r := getRouter(s.T())
	req, _ := http.NewRequest("GET", "/version", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		s.Failf("HttpStatus Failed", "expect[%d], got[%d]", http.StatusOK, w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
	s.NoError(err)

	resp := make(map[string]string)
	err = json.Unmarshal(body, &resp)
	s.NoError(err)

	keys := []string{"Name", "Commit", "Version", "Description"}
	for _, key := range keys {
		if _, ok := resp[key]; !ok {
			s.Failf("HttpBofy Failed", "expect key[%s] exists, but not found", key)
		}
	}
}

func TestApiTestSuite(t *testing.T) {
	s := &apiTestSuite{}
	suite.Run(t, s)
}
