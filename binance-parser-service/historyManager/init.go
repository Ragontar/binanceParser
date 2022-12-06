package historyManager

var HistoryStorageDB *DBHistoryStorage

func init() {
	var err error
	if HistoryStorageDB, err = NewDBHistoryStorage(); err != nil {
		panic(err)
	}
}
