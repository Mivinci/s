package model

import (
	"time"

	"github.com/mivinci/shortid"
)

type User struct {
	ID    int       `json:"id" storm:"increment"`
	Name  string    `json:"name" storm:"index"`
	Mail  string    `json:"mail" storm:"unique"`
	Perm  int8      `json:"perm" storm:"index"`
	Ctime time.Time `json:"ctime"`
}

func (u User) Uid() string {
	return shortid.String(u.ID)
}

// func (u User) MarshalJSON() ([]byte, error) {
// 	type Alias User
// 	return json.Marshal(struct {
// 		Ctime int64 `json:"ctime"`
// 		Alias
// 	}{
// 		Ctime: time.Now().Unix(),
// 		Alias: (Alias)(u),
// 	})
// }

// func (u *User) UnmarshalJSON(b []byte) error {
// 	type Alias User
// 	obj := struct {
// 		Ctime int64 `json:"ctime"`
// 		Alias
// 	}{
// 		Alias: (Alias)(*u),
// 	}
// 	if err := json.Unmarshal(b, &obj); err != nil {
// 		return err
// 	}
// 	u.Ctime = time.Unix(obj.Ctime, 0)
// 	return nil
// }
