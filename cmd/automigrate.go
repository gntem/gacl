package cmd

func Automigrate(db *DB) {
	db.DropTableIfExists(&User{}, &Group{}, &Permission{})
	db.AutoMigrate(&User{}, &Group{}, &Permission{})
}
