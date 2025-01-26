package openai

import (
	"fmt"
	"os"
	"testing"
)

func TestModels(t *testing.T) {
	client := NewClient(&ClientConfig{
		BaseUrl: "https://api.deepseek.com/v1",
		ApiKey:  os.Getenv("OPENAI_KEY"),
	})
	models, e := client.Models()
	if e != nil {
		t.Error(e.Error())
	}
	fmt.Println(models)
}
