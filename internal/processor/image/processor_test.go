package image_test

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/isutare412/imageer/internal/processor/domain"
	"github.com/isutare412/imageer/internal/processor/image"
	"github.com/isutare412/imageer/pkg/images"
)

//go:embed testdata/*
var testFS embed.FS

var nonAllowedChars = regexp.MustCompile(`[^a-zA-Z\-]`)

func TestProcessor_Process(t *testing.T) {
	type testSet struct {
		name         string // description of this test case
		fileName     string
		prepareInput func(t *testing.T, tt testSet) domain.RawImage
		preset       domain.Preset
	}

	tests := []testSet{
		{
			name:     "jpeg-astronaut-cover",
			fileName: "testdata/jpeg-astronaut-2000x1360.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitCover),
				Anchor:  lo.ToPtr(images.AnchorSmart),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](400),
			},
		},
		{
			name:     "jpeg-astronaut-only-width",
			fileName: "testdata/jpeg-astronaut-2000x1360.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Width:   lo.ToPtr[int64](400),
			},
		},
		{
			name:     "jpeg-astronaut-no-dimension",
			fileName: "testdata/jpeg-astronaut-2000x1360.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(30),
			},
		},
		{
			name:     "jpeg-astronaut-fill",
			fileName: "testdata/jpeg-astronaut-2000x1360.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitFill),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](400),
			},
		},
		{
			name:     "jpeg-astronaut-contain",
			fileName: "testdata/jpeg-astronaut-2000x1360.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitContain),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](400),
			},
		},
		{
			name:     "jpeg-mountain-cover",
			fileName: "testdata/jpeg-mountain-2000x1332.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitCover),
				Anchor:  lo.ToPtr(images.AnchorSmart),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](400),
			},
		},
		{
			name:     "jpeg-mountain-fill",
			fileName: "testdata/jpeg-mountain-2000x1332.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitFill),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](400),
			},
		},
		{
			name:     "jpeg-mountain-contain",
			fileName: "testdata/jpeg-mountain-2000x1332.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitContain),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](400),
			},
		},
		{
			name:     "jpeg-sprout-larger-dimension",
			fileName: "testdata/jpeg-sprout-620x427.jpg",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatJPEG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitCover),
				Anchor:  lo.ToPtr(images.AnchorSmart),
				Width:   lo.ToPtr[int64](800),
				Height:  lo.ToPtr[int64](800),
			},
		},
		{
			name:     "png-mistletoe-cover",
			fileName: "testdata/png-mistletoe-1920x1920.png",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatPNG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitCover),
				Anchor:  lo.ToPtr(images.AnchorSmart),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](300),
			},
		},
		{
			name:     "png-mistletoe-contain",
			fileName: "testdata/png-mistletoe-1920x1920.png",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatPNG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitContain),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](300),
			},
		},
		{
			name:     "png-mistletoe-fill",
			fileName: "testdata/png-mistletoe-1920x1920.png",
			prepareInput: func(t *testing.T, tt testSet) domain.RawImage {
				buf, err := testFS.ReadFile(tt.fileName)
				require.NoError(t, err)
				return domain.RawImage{
					Data:   buf,
					Format: images.FormatPNG,
				}
			},
			preset: domain.Preset{
				Format:  images.FormatWebp,
				Quality: images.Quality(90),
				Fit:     lo.ToPtr(images.FitFill),
				Width:   lo.ToPtr[int64](400),
				Height:  lo.ToPtr[int64](300),
			},
		},
	}

	cleanUpTestOutputs(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := image.NewProcessor()
			input := tt.prepareInput(t, tt)
			out, err := c.Process(t.Context(), input, tt.preset)
			require.NoError(t, err)

			outFileName := testImageOutputName(tt.name, tt.fileName, out.Format)
			err = os.WriteFile(outFileName, out.Data, 0o644)
			require.NoError(t, err)
		})
	}
}

func cleanUpTestOutputs(t *testing.T) {
	files, err := filepath.Glob("testdata/*.out.*")
	require.NoError(t, err)

	for _, f := range files {
		err := os.Remove(f)
		require.NoError(t, err)
	}
}

func testImageOutputName(testName, fileName string, format images.Format) string {
	safeName := nonAllowedChars.ReplaceAllString(testName, "-")

	dir := filepath.Dir(fileName)
	base := filepath.Base(fileName)
	ext := filepath.Ext(base)

	return fmt.Sprintf("%s/%s-%s.out.%s", dir, base[0:len(base)-len(ext)],
		safeName, format.ToExtension())
}
