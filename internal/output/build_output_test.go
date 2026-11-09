package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBuildOutput(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		output      BuildOutputOption
		invalid     bool
	}{
		{
			description: "hyphen",
			input:       "-",
			output: BuildOutputOption{
				Type: BuildOutputStdout,
			},
		},
		{
			description: "just-a-path",
			input:       "/tmp",
			output: BuildOutputOption{
				Type: BuildOutputLocalDir,
				Path: "/tmp",
			},
		},
		{
			description: "normal-path",
			input:       "type=local,dest=/tmp",
			output: BuildOutputOption{
				Type: BuildOutputLocalDir,
				Path: "/tmp",
			},
		},
		{
			description: "image",
			input:       "type=image",
			output: BuildOutputOption{
				Type: BuildOutputImage,
			},
		},
		{
			description: "image-ignores-dest",
			input:       "type=image,dest=/tmp",
			output: BuildOutputOption{
				Type: BuildOutputImage,
			},
		},
		{
			description: "image-with-name",
			input:       "name=example.com/image:tag,type=image",
			output: BuildOutputOption{
				Type: BuildOutputImage,
				Name: "example.com/image:tag",
			},
		},
		{
			description: "image-push-false",
			input:       "type=image,push=false",
			output: BuildOutputOption{
				Type: BuildOutputImage,
			},
		},
		{
			description: "image-push-true",
			input:       "type=image,push=true,name=example.com/image",
			output: BuildOutputOption{
				Type: BuildOutputRegistry,
				Name: "example.com/image",
			},
		},
		{
			description: "registry",
			input:       "type=registry",
			output: BuildOutputOption{
				Type: BuildOutputRegistry,
			},
		},
		{
			description: "registry-with-name",
			input:       "type=registry,name=example.com/image",
			output: BuildOutputOption{
				Type: BuildOutputRegistry,
				Name: "example.com/image",
			},
		},
		{
			description: "missing-type",
			input:       "name=example.com/image",
			invalid:     true,
		},
		{
			description: "invalid-type",
			input:       "type=wrong,dest=/tmp",
			invalid:     true,
		},
		{
			description: "duplicate-name",
			input:       "type=image,name=a,name=b",
			invalid:     true,
		},
		{
			description: "invalid-push",
			input:       "type=image,push=maybe",
			invalid:     true,
		},
		{
			description: "push-with-registry",
			input:       "type=registry,push=true",
			invalid:     true,
		},
		{
			description: "name-with-local",
			input:       "type=local,dest=/tmp,name=example.com/image",
			invalid:     true,
		},
		{
			description: "push-with-tar",
			input:       "type=tar,dest=/tmp/out.tar,push=true",
			invalid:     true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			result, err := GetBuildOutput(testCase.input)
			if testCase.invalid {
				assert.Errorf(t, err, "expected %q to be rejected", testCase.input)
				return
			}
			require.NoErrorf(t, err, "expected to be able to parse %q", testCase.input)
			assert.Equal(t, testCase.output, result)
		})
	}
}
