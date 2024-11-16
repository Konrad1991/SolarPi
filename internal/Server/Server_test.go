package Server

import (
  "testing"
  "github.com/gin-gonic/gin"
  "reflect"
)

func TestCreateRouter(t *testing.T) {
  got := createRouter("192.168.0.1")
  want := gin.Default()
  
  if reflect.TypeOf(got) != reflect.TypeOf(want) {
    t.Errorf("type of Router object is not gin default router")
  }


}
