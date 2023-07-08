package pdf

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	context2 "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestNewFromHTML(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.TODO(), context2.FieldStage, "dev")
	res, err := NewFromHTML(ctx, `<html>
<body><h1>Hello <mark>world!</mark></h1></body>
</html>`)
	assert.NoError(t, err)
	fmt.Println(base64.StdEncoding.EncodeToString(res))
}
