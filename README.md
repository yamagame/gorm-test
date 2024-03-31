# gorm 速度計測テスト

MacBook Pro M1 Max で計測  
join1 は 13.621s、join2 は 5.037s、DBへの問い合わせを減らした方が早い。

```sh
$ docker-compose up -d
$ go build -o build/ main.go
$ build/main migrate
# 10,000件作成
$ time build/main create 10000

real    0m29.018s
user    0m1.998s
sys     0m4.814s
$ time build/main join1

real    0m13.621s
user    0m0.893s
sys     0m1.447s
$ time build/main join2

real    0m5.037s
user    0m0.337s
sys     0m0.500s
```

以下、join1 と join2 のコード

```go
func join1(db *gorm.DB, ulids []string) []*Result {
	users := []*Result{}
	for _, ulid := range ulids {
		var user model.User
		db.Select("ULID", "Name").Where("ul_id = ?", ulid).Take(&user)
		var addr model.Addr
		db.Select("ULID", "City").Where("ul_id = ?", ulid).Take(&addr)
		var school model.School
		db.Select("ULID", "Name").Where("ul_id = ?", ulid).Take(&school)
		users = append(users, &Result{
			Name:       user.Name,
			City:       addr.City,
			SchoolName: school.Name,
		})
	}
	return users
}
```

```go
func join2(db *gorm.DB, ulids []string) []*Result {
	users := []*Result{}
	for _, ulid := range ulids {
		var user Result
		db.Table("users").
			Select("users.ul_id as ulid", "users.name", "addrs.city", "schools.name as school_name").
			Where("users.ul_id = ?", ulid).
			Joins("left join addrs on addrs.ul_id = users.ul_id").
			Joins("left join schools on schools.ul_id = users.ul_id").
			Take(&user)
		users = append(users, &user)
	}
	return users
}
```
