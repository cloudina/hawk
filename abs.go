type GCS_Manager struct {

}

func (self *GCS_Manager ) readFile(bucket string, item string) ([] byte, error) {
func (self *S3_Manager ) copyFile(bucket string, item string, other string) (error){
func (self *S3_Manager ) deleteFile(bucket string, item string) (error) {
