package env

var (
	DB_USER     string
	DB_PASSWORD string
	DB_ADDR     string
	DB_DATABASE string
)

// TODO read from environment

func init() {
	// var ok bool
	// if DB_USER, ok = os.LookupEnv("DB_USER"); !ok {
	// 	panic(fmt.Sprintf("Cannot lookup %s", "DB_USER"))
	// }
	// if DB_PASSWORD, ok = os.LookupEnv("DB_PASSWORD"); !ok {
	// 	panic(fmt.Sprintf("Cannot lookup %s", "DB_PASSWORD"))
	// }
	// if DB_ADDR, ok = os.LookupEnv("DB_ADDR"); !ok {
	// 	panic(fmt.Sprintf("Cannot lookup %s", "DB_ADDR"))
	// }
	// if DB_DATABASE, ok = os.LookupEnv("DB_DATABASE"); !ok {
	// 	panic(fmt.Sprintf("Cannot lookup %s", "DB_DATABASE"))
	// }
}
