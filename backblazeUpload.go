package main
import (
  "os"
  "path/filepath"
  "gopkg.in/kothar/go-backblaze.v0"

)


func uploadFileToBackblaze(filename string){
    b2, _ := backblaze.NewB2(backblaze.Credentials{
      AccountID: "23547fcec776",
      ApplicationKey: "0016ab4da23ef8548aa6d19c77e0eada59ae55764e",
  })

  bucket, _ := b2.Bucket("Hack4Impact")

  path:= filename
  reader, _ := os.Open(path)
  name := filepath.Base(path)
  metadata := make(map[string]string)

      bucket.UploadFile(name, metadata, reader)
}
