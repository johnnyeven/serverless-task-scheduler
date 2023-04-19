package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
)

func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(
		gen.Config{
			OutPath: "db/api",
			Mode:    gen.WithoutContext,
			//if you want the nullable field generation property to be pointer type, set FieldNullable true
			FieldNullable: true,
			//if you want to generate index tags from database, set FieldWithIndexTag true
			/* FieldWithIndexTag: true,*/
			//if you want to generate type tags from database, set FieldWithTypeTag true
			/* FieldWithTypeTag: true,*/
			//if you need unit tests for query code, set WithUnitTest true
			/* WithUnitTest: true, */
		},
	)

	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessray or it will panic
	db, _ := gorm.Open(mysql.Open(os.Getenv("DATABASE_DSN")))
	g.UseDB(db)

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
	g.ApplyBasic(
		g.GenerateModelAs("t_task", "Task"),
	)

	// execute the action of code generation
	g.Execute()
}
