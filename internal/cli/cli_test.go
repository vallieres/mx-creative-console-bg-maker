package cli_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/cli"
	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

// newMockedApp builds an App backed by a Service with mock dependencies and
// returns the mock file system so tests can seed files.
func newMockedApp() (*cli.App, *processor.TestMockFileSystem) {
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	return cli.NewAppWithProcessor(service), fs
}

func TestNewApp(t *testing.T) {
	t.Parallel()

	app := cli.NewApp()

	assert.NotNil(t, app)
}

func TestNewAppWithProcessor(t *testing.T) {
	t.Parallel()

	customProcessor := processor.NewService()
	app := cli.NewAppWithProcessor(customProcessor)

	assert.NotNil(t, app)
}

func TestApp_Run_Split_Success(t *testing.T) {
	t.Parallel()

	app, fs := newMockedApp()

	fs.AddFile("/test/image.jpg", []byte("fake image"))

	err := app.Run([]string{"ccbm", "split", "-f", "/test/image.jpg"})
	assert.NoError(t, err)
}

func TestApp_Run_Split_SpecialCharactersInPath(t *testing.T) {
	t.Parallel()

	app, fs := newMockedApp()

	path := "/test/dir with spaces/imágé (1).jpg"
	fs.AddFile(path, []byte("fake image"))

	err := app.Run([]string{"ccbm", "split", "--file", path})
	assert.NoError(t, err)
}

func TestApp_Run_Split_MissingFileFlag(t *testing.T) {
	t.Parallel()

	app, _ := newMockedApp()

	err := app.Run([]string{"ccbm", "split"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")
}

func TestApp_Run_NoArguments(t *testing.T) {
	t.Parallel()

	// Without a subcommand, the root command shows help and succeeds.
	app, _ := newMockedApp()

	for _, args := range [][]string{{}, {"ccbm"}} {
		err := app.Run(args)
		assert.NoError(t, err)
	}
}

func TestApp_Run_UnknownCommand(t *testing.T) {
	t.Parallel()

	app, _ := newMockedApp()

	err := app.Run([]string{"ccbm", "bogus"})

	require.Error(t, err)
}

func TestApp_Run_Split_ProcessorError(t *testing.T) {
	t.Parallel()

	fs := processor.NewTestMockFileSystem()
	fs.OpenFunc = func(_ string) (io.ReadCloser, error) {
		return nil, errors.New("file not found")
	}

	decoder := processor.NewTestMockImageDecoder(nil, "", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	err := app.Run([]string{"ccbm", "split", "-f", "/nonexistent/image.jpg"})

	require.Error(t, err)
}
