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
	"github.com/lsytj0413/fyllo/pkg/random"
	randomBuilder "github.com/lsytj0413/fyllo/pkg/random/builder"
	"github.com/lsytj0413/fyllo/pkg/segment"
	segmentBuilder "github.com/lsytj0413/fyllo/pkg/segment/builder"
	"github.com/lsytj0413/fyllo/pkg/snowflake"
	snowflakeBuilder "github.com/lsytj0413/fyllo/pkg/snowflake/builder"
)

func buildRandomProvider(name string, args string) (random.Provider, error) {
	builder, err := randomBuilder.NewBuilder(&randomBuilder.Options{
		ProviderName: name,
		ProviderArgs: args,
	})
	if err != nil {
		return nil, err
	}

	return builder.Build()
}

func buildSnowflakeProvider(name string, args string) (snowflake.Provider, error) {
	builder, err := snowflakeBuilder.NewBuilder(&snowflakeBuilder.Options{
		ProviderName: name,
		ProviderArgs: args,
	})
	if err != nil {
		return nil, err
	}

	return builder.Build()
}

func buildSegmentProvider(name string, args string) (segment.Provider, error) {
	builder, err := segmentBuilder.NewBuilder(&segmentBuilder.Options{
		ProviderName: name,
		ProviderArgs: args,
	})
	if err != nil {
		return nil, err
	}

	return builder.Build()
}
