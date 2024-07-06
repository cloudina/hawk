type GCS_Manager struct {

}

func (self *GCS_Manager ) readFile(bucket string, item string) ([] byte, error) 
func (self *GCS_Manager ) copyFile(bucket string, item string, other string) (error)
func (self *GCS_Manager ) deleteFile(bucket string, item string) (error) 
