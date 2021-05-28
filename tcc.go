package dtm

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/yedf/dtm/common"
)

type Tcc struct {
	TccData
	Server string
}

type TccData struct {
	Gid           string    `json:"gid"`
	TransType     string    `json:"trans_type"`
	Steps         []TccStep `json:"steps"`
	QueryPrepared string    `json:"query_prepared"`
}
type TccStep struct {
	Try     string `json:"try"`
	Confirm string `json:"confirm"`
	Cancel  string `json:"cancel"`
	Data    string `json:"data"`
}

func TccNew(server string, gid string) *Tcc {
	return &Tcc{
		TccData: TccData{
			Gid:       gid,
			TransType: "tcc",
		},
		Server: server,
	}
}
func (s *Tcc) Add(try string, confirm string, cancel string, data interface{}) error {
	logrus.Printf("tcc %s Add %s %s %s %v", s.Gid, try, confirm, cancel, data)
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	step := TccStep{
		Try:     try,
		Confirm: confirm,
		Cancel:  cancel,
		Data:    string(d),
	}
	s.Steps = append(s.Steps, step)
	return nil
}

func (s *Tcc) Prepare(queryPrepared string) error {
	s.QueryPrepared = queryPrepared
	logrus.Printf("preparing %s body: %v", s.Gid, &s.TccData)
	resp, err := common.RestyClient.R().SetBody(&s.TccData).Post(fmt.Sprintf("%s/prepare", s.Server))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("prepare failed: %v", resp.Body())
	}
	return nil
}

func (s *Tcc) Commit() error {
	logrus.Printf("committing %s body: %v", s.Gid, &s.TccData)
	resp, err := common.RestyClient.R().SetBody(&s.TccData).Post(fmt.Sprintf("%s/commit", s.Server))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("commit failed: %v", resp.Body())
	}
	return nil
}