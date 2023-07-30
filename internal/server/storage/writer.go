package storage

//func WriteStorageToFile(s *Storage, interval time.Duration, path string, restore bool) {
//	log := logger.Default()
//	var file *os.File
//	var err error
//	if restore {
//		file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//		if err != nil {
//			log.Fatalf("cannot open file: %e", err)
//			return
//		}
//		// тут должна быть функция парсинга данных из хранилища
//	} else {
//		file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
//		if err != nil {
//			log.Fatalf("cannot open file: %e", err)
//			return
//		}
//	}
//}
