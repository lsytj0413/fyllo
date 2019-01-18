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

package mysql

import (
	"database/sql"
	"fmt"

	mysqlDriver "github.com/go-sql-driver/mysql"

	ierror "github.com/lsytj0413/fyllo/pkg/error"
	"github.com/lsytj0413/fyllo/pkg/segment"
	"github.com/lsytj0413/fyllo/pkg/segment/internal"
)

var _ = mysqlDriver.DeregisterLocalFile

const (
	// ProviderName for the mysql segment provider
	ProviderName = "mysql"

	// +-------------+--------------+------+-----+-------------------+-----------------------------+
	// | Field       | Type         | Null | Key | Default           | Extra                       |
	// +-------------+--------------+------+-----+-------------------+-----------------------------+
	// | tag         | varchar(128) | NO   | PRI | NULL              |                             |
	// | max_id      | bigint(20)   | NO   |     | 0                 |                             |
	// | step        | bigint(20)   | NO   |     | 100               |                             |
	// | description | varchar(256) | NO   |     |                   |                             |
	// | create_time | timestamp    | NO   |     | CURRENT_TIMESTAMP |                             |
	// | update_time | timestamp    | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
	// +-------------+--------------+------+-----+-------------------+-----------------------------+

	// CREATE TABLE `fyllo_segment` (
	// 	`tag` varchar(128) NOT NULL,
	// 	`max_id` bigint(20) NOT NULL DEFAULT '0',
	// 	`step` bigint(20) NOT NULL DEFAULT '100',
	// 	`description` varchar(256) NOT NULL DEFAULT '',
	// 	`create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	// 	`update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	// 	PRIMARY KEY (`tag`)
	//   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	tableName = "fyllo_segment"
)

// Options is mysql segment provider option
type Options struct {
	Args string
}

type mysqlStorage struct {
	db *sql.DB
}

func (m *mysqlStorage) List() ([]string, error) {
	rows, err := m.db.Query("select tag from fyllo_segment")
	if err != nil {
		return nil, ierror.NewError(ierror.EcodeSegmentQueryFailed, fmt.Sprintf("Query failed, %v", err))
	}
	if rows == nil {
		return nil, ierror.NewError(ierror.EcodeSegmentQueryFailed, "Query failed, return nil Rows")
	}
	defer rows.Close()

	r := make([]string, 0)
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, ierror.NewError(ierror.EcodeSegmentQueryFailed, fmt.Sprintf("Query failed, %v", err))
		}

		r = append(r, name)
	}
	return r, nil
}

func (m *mysqlStorage) Obtain(tag string) (*internal.TagItem, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, ierror.NewError(ierror.EcodeSegmentQueryFailed, fmt.Sprintf("Query failed, %v", err))
	}
	defer tx.Rollback()
	item, err := func() (*internal.TagItem, error) {
		var (
			maxID uint64
			step  uint64
			desc  string
		)
		var rs sql.Result
		rs, err = tx.Exec("update fyllo_segment set max_id=max_id+step where tag = ?", tag)
		if err != nil {
			return nil, err
		}
		var rowAffected int64
		rowAffected, err = rs.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rowAffected != 1 {
			return nil, fmt.Errorf("no row")
		}

		err = tx.QueryRow("select max_id, step, description from fyllo_segment where tag = ?", tag).Scan(&maxID, &step, &desc)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, err
			}
			return nil, err
		}
		return &internal.TagItem{
			Tag:         tag,
			Max:         maxID,
			Min:         maxID - step + 1,
			Description: desc,
		}, nil
	}()
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return item, nil
}

// NewProvider return mysql segment provider implement
func NewProvider(options *Options) (segment.Provider, error) {
	if 0 == len(options.Args) {
		return nil, ierror.NewError(ierror.EcodeInitFailed, "segment provider[mysql] arguments should not be empty")
	}

	db, err := sql.Open("mysql", options.Args)
	if err != nil {
		return nil, ierror.NewError(ierror.EcodeInitFailed, fmt.Sprintf("open %s failed, %v", options.Args, err))
	}

	return internal.NewProvider(ProviderName, &mysqlStorage{
		db: db,
	})
}
