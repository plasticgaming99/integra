package util

import (
	"errors"

	"github.com/cavaliergopher/grab/v3"
)

type RemoteDB struct {
	repos []string // including url
	path  string   // dbdir
}

func (rdb *RemoteDB) AddRepo(url string) {
	rdb.repos = append(rdb.repos, url)
}

func (rdb *RemoteDB) SyncRemoteDB() error {
	for _, s := range rdb.repos {
		_, err := grab.Get(rdb.path, s)
		if err != nil {
			return errors.New("error synchronizing database")
		}
	}
	return nil
}
