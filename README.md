### Connect to database

```go
// mysql
var db, err = gorm.Open(mysql.New(mysql.Config{
DSN: "root:rootpass@tcp(127.0.0.1:3306)/gift-finder?charset=utf8&parseTime=True&loc=Local",
}), &gorm.Config{})

// postgresql
var db, err = gorm.Open(postgres.New(postgres.Config{
DSN: "host=127.0.0.1 user=root password=rootpass dbname=gift-finder port=5432 sslmode=disable",
}), &gorm.Config{})
```