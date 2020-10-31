package main

import   (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	var files []string

	const (
		host     = "localhost"
		port     = 5432
		user     = "moraghan"
		password = "Coleus92"
		dbname   = "moraghan"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err1 := sql.Open("postgres", psqlInfo)
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()

	err1 = db.Ping()
	if err1 != nil {
		panic(err1)
	}

	fmt.Print(psqlInfo)
	//root := "/Volumes/1Tb_exfat/Movies/"
	root := "/Volumes/video/movies/"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for a, file := range files {
		fileStat, _ := os.Stat(file)
		extension := filepath.Ext(file)

		if extension != ".xml" {

			var containsYear, _ = regexp.MatchString(`\b(19|20)\d{2}\b`, file)

			fmt.Println(a, file, extension)
			fmt.Println("File Name:", fileStat.Name())        // Base name of the file
			fmt.Println("Size:", fileStat.Size())             // Length in bytes for regular files
			fmt.Println("Last Modified:", fileStat.ModTime()) // Last modification time
			fmt.Println("Contains Year: ", containsYear)

			if fileStat.Name() != ".DS_Store" && !fileStat.IsDir() {
				sqlStatement := `	INSERT INTO video (id, file_name, file_size, contains_year_ind)
									VALUES ($1, $2, $3, $4)`
				_, err = db.Exec(sqlStatement, a, fileStat.Name(), fileStat.Size(), containsYear)

				if err != nil {
					panic(err)
				}
			}
		}

	}

}
