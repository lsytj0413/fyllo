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
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lsytj0413/fyllo/pkg/random"
	"github.com/lsytj0413/fyllo/pkg/segment"
	"github.com/lsytj0413/fyllo/pkg/snowflake"
	"github.com/lsytj0413/fyllo/pkg/version"
)

type versionService struct {
}

func (s *versionService) Install(engine *gin.Engine) error {
	engine.Use(errorMiddleware()).Use(logMiddleware()).Use(jsonRespMiddleware())
	engine.GET("/version", wrapperHandler(s.Version))
	return nil
}

func (s *versionService) Version(c *gin.Context) (interface{}, error) {
	return map[string]string{
		"Version":     version.Version,
		"Commit":      version.Commit,
		"Name":        "fyllo",
		"Description": "A distributed, unique ID generation service.",
	}, nil
}

type snowflakeService struct {
	provider snowflake.Provider
}

func (s *snowflakeService) Install(engine *gin.Engine) error {
	engine.GET("/api/snowflake", wrapperHandler(s.Next))
	return nil
}

func (s *snowflakeService) Next(c *gin.Context) (interface{}, error) {
	tag := c.Query("tag")
	tagID, err := strconv.Atoi(tag)
	if err != nil {
		return nil, err
	}

	n, err := s.provider.Next(&snowflake.Arguments{
		Tag: uint64(tagID),
	})
	if err != nil {
		return nil, err
	}
	return n, nil
}

type randomService struct {
	provider random.Provider
}

func (s *randomService) Install(engine *gin.Engine) error {
	engine.GET("/api/random", wrapperHandler(s.Next))
	return nil
}

func (s *randomService) Next(c *gin.Context) (interface{}, error) {
	n, err := s.provider.Next(&random.Arguments{})
	if err != nil {
		return nil, err
	}
	return n, nil
}

type segmentService struct {
	provider segment.Provider
}

func (s *segmentService) Install(engine *gin.Engine) error {
	engine.GET("/api/segment", wrapperHandler(s.Next))
	return nil
}

func (s *segmentService) Next(c *gin.Context) (interface{}, error) {
	tag := c.Query("tag")
	if tag == "" {
		return nil, fmt.Errorf("tag can't be empty")
	}

	n, err := s.provider.Next(&segment.Arguments{
		Tag: tag,
	})
	if err != nil {
		return nil, err
	}
	return n, nil
}
