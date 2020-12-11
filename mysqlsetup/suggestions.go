package mysqlsetup

import (
	"database/sql"
	"fmt"
	"os"

	"encoding/json"

	fire "github.com/Nv7-Github/firebase"
	database "github.com/Nv7-Github/firebase/db"
	_ "github.com/go-sql-driver/mysql" // mysql
)

const serviceAccount = `{ "type": "service_account", "project_id": "elementalserver-8c6d0", "private_key_id": "78ba8b0bb00e5233e4ac4cb5e640ec6d6c56eb6b", "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDEw7lgdH3/H5F2\nMUkY6ImrSkS9sqvKmgAlW1CLxn8p/V4yeSbFvFiUv9mBTlKYHAnE/wmez8F1OyOT\ncnvG0AFicmkV/2QJSjkW8yvkFB7S2Br7rWkPsmacA1RlWHlctybSUo9/HkooDl5j\nkqtbJsT2Vt/nWlsDMmZTKTCJoRQuP8jATbeA16XcQ/Gl8mFap1lC89OE7wxrx5Eq\nTKqt4AkhKhBQhoNr9VTFRCBkhaBMq830/MKaZpbLDVa6qTk4c9QOs9hohtpUAup6\n2nvxX6k9l/K9nPJI9Gp7lOhQ951/yws8c+DGbuSdjHnj1sKQlWmOY6hMkS2i4YHX\nWWcEheYnAgMBAAECggEAGb+DApw74KbA4jaQ2jGT0lZlqG05DcoZOso4QBI5kcUW\nDoTMDhQXg1+XltQo+r6wiJbXK3EEX9LdVO4mRF3z0G4oUjiZXp3X2qj3lWEMp4qf\n/U8z8FnoE4JcCOcK+pb8/YjQPlI4YgV/VIhc5BCutY2ovx2Ty1dNDJTXRStO+L4k\n1rnkaL/zwbxUUBLEb4rYuAjw7cOyxk081yXdogs34IQd/9aLVZnIri5A6rSDYFFw\n1A2aJzNBolur+/5m57dT6OqVjgD99zjWEJ9lq61oNMtYVetGY9bFis59JnUdOUHg\n11B2RMMOAu4lqDaCFEvmqiw09OEVgynlA72XUhaAAQKBgQD9+FMOjM9dai9RoDBZ\nrgEQQfqwPXS1z24Fqlx1zBgI6wUCvVU1qfnMBCQsrAwzfev+nt44yOBctuPr7ZUC\nvBvSGO6Dd50Wwnd24wHrhd+45FX6pHbdSItbkpMyuA6OHHigXJNKaHib6G7EqHs5\nIK88GrawYtYD2abWfQJwtJhTgQKBgQDGVlhqyOrGfBub1NSFsey5oltCH7KuMV1F\nozNoOlpU4EtOaYE1g71+9uNSf3n1vO8PLCLj8ZoF2fk6VV2ttME2fXUjfZz9HUgw\ncfOZpD5j3XM/s8mJFncbR+5eC/KAoXA0oUVvhkffdEaGtrInrobOOcyETGqdPcjA\n6oAvKFntpwKBgQCWJDA18djFiPjgcKsk2VGXounpNuvAcBjDEKwIl9e9rfMQY430\nY8BhdDFOl4e/CTpzFMibGWZKaXTlDVeCfmKUGlknL5eW1PB7QEjqTAKu845A1unO\neAyq3kRXP6ibKwnFA/Wvj4N96DNT36a5ZzExfzlxnXyYWhvfwZenuZw0AQKBgH0c\nqKer2BWe4leZmPpBM4AiP4jlr/QcJacw/NOpw6O43Sg4e45DbTzzBpDa4xc1uGOM\nxvGdTTiVuJaolPBnjl4OI99gdLBiUVBmAXGQ3t5mKjYr9lyotDecV2wyAyZLMBmz\nBbcFML9vfLGr+5P2jwj2AuINxk8sU0AGbRfST3APAoGBAMOJsmiBoXBiOxi52lyM\nZN8jxyTHd9LwnTgPHQd2JedBi7EIJ3j3T+QP3Z3SENMMImQr6MOda8otrTyqpMTp\nDS+pTomSwTCCEir7bSVpi7QMejchURVYM/PmMwhso1vocZBM3YHvxLtGAnFOu7BM\nQ39vHDC9jyj00STzo/+fD6X3\n-----END PRIVATE KEY-----\n", "client_email": "firebase-adminsdk-7nmm6@elementalserver-8c6d0.iam.gserviceaccount.com", "client_id": "113854670633531537114", "auth_uri": "https://accounts.google.com/o/oauth2/auth", "token_uri": "https://oauth2.googleapis.com/token", "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs", "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-7nmm6%40elementalserver-8c6d0.iam.gserviceaccount.com" }`

// Suggestion has the data for a suggestion
type Suggestion struct {
	Creator string   `json:"creator"`
	Name    string   `json:"name"`
	Votes   int      `json:"votes"`
	Color   Color    `json:"color"`
	Voted   []string `json:"voted"`
}

// Color has the data for a suggestion's color
type Color struct {
	Base       string  `json:"base"`
	Lightness  float32 `json:"lightness"`
	Saturation float32 `json:"saturation"`
}

const (
	dbUser     = "u29_c99qmCcqZ3"
	dbPassword = "j8@tJ1vv5d@^xMixUqUl+NmA"
	dbName     = "s29_nv7haven"
)

// Mysqlsetup adds the elements to the mysql db
func Mysqlsetup() {
	// mysql
	db, err := sql.Open("mysql", "jdbc:mysql://"+dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("opened db")

	firebaseapp, err := fire.CreateAppWithServiceAccount("https://elementalserver-8c6d0.firebaseio.com", "AIzaSyCsqvV3clnwDTTgPHDVO2Yatv5JImSUJvU", []byte(serviceAccount))
	if err != nil {
		panic(err)
	}

	firedb := database.CreateDatabase(firebaseapp)

	var suggs []Suggestion
	data, err := firedb.Get("suggestions")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &suggs)

	fmt.Println("Got suggs")

	insElem, err := db.Prepare("INSERT INTO suggestions VALUES( ?, ?, ?, ?, ? )")
	if err != nil {
		panic(err)
	}
	defer insElem.Close()
	fmt.Println("Prepared command")
	for _, val := range suggs {
		a, _ := json.Marshal(val.Voted)
		fmt.Println("ready to exec")
		_, err = insElem.Exec(val.Name, val.Color, val.Creator, a, val.Votes)
		if err != nil {
			panic(err)
		}
		fmt.Println("execed!")
	}
}
