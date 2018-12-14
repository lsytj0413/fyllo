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

	"github.com/lsytj0413/fyllo/pkg/random"
	"github.com/lsytj0413/fyllo/pkg/segment"
	"github.com/lsytj0413/fyllo/pkg/snowflake"
)

type versionService struct {
}

func (s *versionService) Install(engine *gin.Engine) error {
	return nil
}

type snowflakeService struct {
	provider snowflake.Provider
}

func (s *snowflakeService) Install(engine *gin.Engine) error {
	return nil
}

type randomService struct {
	provider random.Provider
}

func (s *randomService) Install(engine *gin.Engine) error {
	return nil
}

type segmentService struct {
	provider segment.Provider
}

func (s *segmentService) Install(engine *gin.Engine) error {
	return nil
}
