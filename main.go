package main

func main() {
	a := App{}
	db_username := "root"
	db_password := "Marin123"
	db_name := "budget"
	a.Initialize(db_username, db_password, db_name)
	a.Run(":8080")
}
